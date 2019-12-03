package scanner

import (
	"fmt"
	"github.com/profzone/eden-framework/internal/generator/api"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"strings"
)

type ModelScanner struct {
	Api *api.Api
	*EnumScanner
	models map[string]*api.OperatorModel
}

func NewModelScanner(enum *EnumScanner) *ModelScanner {
	return &ModelScanner{
		models:      make(map[string]*api.OperatorModel),
		EnumScanner: enum,
	}
}

func (m *ModelScanner) NewModel(name string, pkgID string) *api.OperatorModel {
	id := strings.Join([]string{pkgID, name}, ".")
	if id == "time.Time" {
		fmt.Println()
	}
	if _, ok := m.models[id]; !ok {
		model := api.NewOperatorModel(name, pkgID)
		m.models[id] = &model
	}

	return m.models[id]
}

func (m *ModelScanner) RegisterModel(model *api.OperatorModel) {
	if _, ok := m.models[model.ID]; !ok {
		m.models[model.ID] = model
	}
}

func (m *ModelScanner) GetModel(name string, pkgID string) *api.OperatorModel {
	id := strings.Join([]string{pkgID, name}, ".")
	if _, ok := m.models[id]; !ok {
		return nil
	}

	return m.models[id]
}

func (m *ModelScanner) GetModelByID(id string) *api.OperatorModel {
	if _, ok := m.models[id]; !ok {
		return nil
	}

	return m.models[id]
}

func (m *ModelScanner) NewInspector(pkg *packages.Package) func(node ast.Node) bool {
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
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					model := m.ResolveStdType(typeSpec.Name.Name, pkg.ID)
					if model != nil {
						m.RegisterModel(model)
						continue
					}
					commentScanner := NewCommentScanner(pkg.Fset, file)
					doc := commentScanner.CommentsOf(typeSpec)
					doc, hasEnum := ParseEnum(doc)
					if hasEnum {
						enumFullPath := strings.Join([]string{pkg.PkgPath, typeSpec.Name.Name}, ".")
						enum := m.Enum(enumFullPath)
						m.Api.AddEnum(enumFullPath, enum)
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					model = m.NewModel(typeSpec.Name.Name, pkg.ID)
					for _, field := range structType.Fields.List {
						var key, keyType, alias, path string
						if field.Names != nil {
							key = field.Names[0].Name
						}
						keyType, pkgName := ParseType(field.Type)
						if pkgName != "" {
							if _, ok := importPath[pkgName]; ok {
								path = importPath[pkgName]
								alias = pkgName
							}
						}
						model.AddField(key, keyType, alias, path)
					}
				}
			}
		}

		return true
	}
}

func (m *ModelScanner) ResolveStdType(typeName, pkgPath string) *api.OperatorModel {
	if typeName == "Time" && pkgPath == "time" {
		model := api.NewTimeModel()
		return &model
	}

	return nil
}
