package scanner

import (
	"fmt"
	"github.com/henrylee2cn/erpc/v6"
	"github.com/sirupsen/logrus"
	"go/ast"
	"profzone/eden-framework/internal/generator/api"
	"reflect"
)

type OperatorScanner struct {
	Api           api.Api
	orphanMethods map[string][]*api.OperatorMethod
	modelScanner  *ModelScanner
}

func NewOperatorScanner(modelScanner *ModelScanner) *OperatorScanner {
	return &OperatorScanner{
		Api:           api.NewApi(),
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

func (v *OperatorScanner) Visit(node ast.Node) (w ast.Visitor) {
	switch node.(type) {
	case *ast.TypeSpec:
		typeSpec := node.(*ast.TypeSpec)
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return nil
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
	case *ast.FuncDecl:
		funcDecl := node.(*ast.FuncDecl)
		if funcDecl.Recv == nil {
			return nil
		}
		starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil
		}
		groupNameIdent, ok := starExpr.X.(*ast.Ident)
		if !ok {
			return nil
		}
		groupName := groupNameIdent.Name
		funcName := funcDecl.Name.Name
		method := api.NewOperatorMethod(nil, funcName)
		if group := v.Api.GetGroup(groupName); group != nil {
			method.Group = group
			group.AddMethod(method)
		} else {
			v.AddOrphanMethod(method, groupName)
		}

		if funcDecl.Type.Params != nil {
			for _, input := range funcDecl.Type.Params.List {
				inputExpr, ok := input.Type.(*ast.StarExpr)
				if !ok {
					break
				}
				inputSelectorExpr, ok := inputExpr.X.(*ast.SelectorExpr)
				if !ok {
					break
				}
				inputIdent := inputSelectorExpr.Sel.Name
				v.modelScanner.RegisterInputModelWithReferer(inputIdent, method)
				method.AddInputDef(inputIdent)
			}
		}

		if funcDecl.Type.Results != nil {
			for _, output := range funcDecl.Type.Results.List {
				outputExpr, ok := output.Type.(*ast.StarExpr)
				if !ok {
					break
				}
				outputSelectorExpr, ok := outputExpr.X.(*ast.SelectorExpr)
				if !ok {
					break
				}
				outputIdent := outputSelectorExpr.Sel.Name
				if outputIdent == reflect.TypeOf(erpc.Status{}).Name() {
					// Status
				} else {
					v.modelScanner.RegisterOutputModelWithReferer(outputIdent, method)
					method.AddOutputDef(outputIdent)
				}
			}
		}
	}
	return v
}
