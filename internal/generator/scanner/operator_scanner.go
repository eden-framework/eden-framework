package scanner

import (
	"fmt"
	"github.com/henrylee2cn/erpc/v6"
	"github.com/profzone/eden-framework/internal/generator/api"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"reflect"
	"strings"
)

type OperatorScanner struct {
	Api           *api.Api
	orphanMethods map[string][]*api.OperatorMethod
	modelScanner  *ModelScanner
}

func NewOperatorScanner(modelScanner *ModelScanner) *OperatorScanner {
	return &OperatorScanner{
		orphanMethods: make(map[string][]*api.OperatorMethod),
		modelScanner:  modelScanner,
	}
}

func (v *OperatorScanner) AddOrphanMethod(method *api.OperatorMethod, groupName string) {
	if _, ok := v.orphanMethods[groupName]; !ok {
		v.orphanMethods[groupName] = make([]*api.OperatorMethod, 0)
	}
	v.orphanMethods[groupName] = append(v.orphanMethods[groupName], method)
}

func (v *OperatorScanner) ResolveOrphanMethod(groupName string) error {
	group := v.Api.GetGroup(groupName)
	if group == nil {
		return fmt.Errorf("the group %s does not exist", groupName)
	}
	if methods, ok := v.orphanMethods[groupName]; ok {
		group.AddMethods(methods...)
	}

	return nil
}

func (v *OperatorScanner) NewInspector(pkg *packages.Package) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		file, ok := node.(*ast.File)
		if !ok {
			return false
		}

		importPath := make(map[string]string)
		for _, ipt := range file.Imports {
			var packageName string
			if ipt.Name == nil {
				packageName = RetrievePackageName(ipt.Path.Value)
			} else {
				packageName = ipt.Name.Name
			}
			importPath[packageName] = strings.Trim(ipt.Path.Value, "\"")
		}

		for _, decl := range file.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.(*ast.GenDecl).Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						structType, ok := typeSpec.Type.(*ast.StructType)
						if !ok {
							return false
						}

						for _, filed := range structType.Fields.List {
							selectorExpr, ok := filed.Type.(*ast.SelectorExpr)
							if !ok {
								continue
							}

							indentPackage, ok := selectorExpr.X.(*ast.Ident)
							if !ok {
								continue
							}

							if indentPackage.Name == "erpc" && selectorExpr.Sel.Name == "CallCtx" {
								group := v.Api.AddGroup(typeSpec.Name.Name)
								if err := v.ResolveOrphanMethod(group.Name); err != nil {
									logrus.Panic(err)
								}
							} else if indentPackage.Name == "erpc" && selectorExpr.Sel.Name == "PushCtx" {
								group := v.Api.AddGroup(typeSpec.Name.Name)
								group.IsPush = true
								if err := v.ResolveOrphanMethod(group.Name); err != nil {
									logrus.Panic(err)
								}
							}
						}
					}
				}
			case *ast.FuncDecl:
				funcDecl := decl.(*ast.FuncDecl)
				if funcDecl.Recv == nil {
					return false
				}
				starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
				if !ok {
					return false
				}
				groupNameIdent, ok := starExpr.X.(*ast.Ident)
				if !ok {
					return false
				}
				groupName := groupNameIdent.Name
				funcName := funcDecl.Name.Name
				method := api.NewOperatorMethod(nil, funcName, str.ToLowerSlashCase(funcName))
				if group := v.Api.GetGroup(groupName); group != nil {
					method.Group = group
					group.AddMethod(method)
				} else {
					v.AddOrphanMethod(method, groupName)
				}

				if funcDecl.Type.Params != nil {
					for _, input := range funcDecl.Type.Params.List {
						var model *api.OperatorModel
						inputExpr, ok := input.Type.(*ast.StarExpr)
						if !ok {
							break
						}
						switch inputExpr.X.(type) {
						case *ast.SelectorExpr:
							inputSelectorExpr, ok := inputExpr.X.(*ast.SelectorExpr)
							if !ok {
								break
							}
							inputAlias, ok := inputSelectorExpr.X.(*ast.Ident)
							if !ok {
								break
							}
							pkgID, ok := importPath[inputAlias.Name]
							if !ok {
								break
							}
							inputIdent := inputSelectorExpr.Sel.Name
							model = v.modelScanner.GetModel(inputIdent, pkgID)
						case *ast.Ident:
							inputIdent := inputExpr.X.(*ast.Ident).Name
							model = v.modelScanner.GetModel(inputIdent, pkg.ID)
						}
						if model != nil {
							method.AddInput(model)
						}
					}
				}

				if funcDecl.Type.Results != nil {
					for _, output := range funcDecl.Type.Results.List {
						var model *api.OperatorModel
						var outputIdent string
						outputExpr, ok := output.Type.(*ast.StarExpr)
						if !ok {
							break
						}

						switch outputExpr.X.(type) {
						case *ast.SelectorExpr:
							outputSelectorExpr, ok := outputExpr.X.(*ast.SelectorExpr)
							if !ok {
								break
							}
							outputAlias, ok := outputSelectorExpr.X.(*ast.Ident)
							if !ok {
								break
							}
							pkgID, ok := importPath[outputAlias.Name]
							if !ok {
								break
							}
							outputIdent = outputSelectorExpr.Sel.Name
							model = v.modelScanner.GetModel(outputIdent, pkgID)
						case *ast.Ident:
							outputIdent = outputExpr.X.(*ast.Ident).Name
							model = v.modelScanner.GetModel(outputIdent, pkg.ID)
						}

						if outputIdent != reflect.TypeOf(erpc.Status{}).Name() {
							if model != nil {
								method.AddOutput(model)
							}
						} else {
							// Status
						}
					}
				}
			}
		}

		return true
	}
}
