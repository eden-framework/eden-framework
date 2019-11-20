package scanner

import (
	"go/ast"
	"profzone/eden-framework/internal/generator/api"
	"strings"
)

type ModelScanner struct {
	inputModelReferer  map[string][]*api.OperatorMethod
	outputModelReferer map[string][]*api.OperatorMethod
}

func NewModelScanner() *ModelScanner {
	return &ModelScanner{
		inputModelReferer:  make(map[string][]*api.OperatorMethod),
		outputModelReferer: make(map[string][]*api.OperatorMethod),
	}
}

func (m *ModelScanner) RegisterInputModelWithReferer(name string, method *api.OperatorMethod) {
	if _, ok := m.inputModelReferer[name]; !ok {
		m.inputModelReferer[name] = make([]*api.OperatorMethod, 0)
	}
	m.inputModelReferer[name] = append(m.inputModelReferer[name], method)
}

func (m *ModelScanner) RegisterOutputModelWithReferer(name string, method *api.OperatorMethod) {
	if _, ok := m.outputModelReferer[name]; !ok {
		m.outputModelReferer[name] = make([]*api.OperatorMethod, 0)
	}
	m.outputModelReferer[name] = append(m.outputModelReferer[name], method)
}

func (m *ModelScanner) Visit(node ast.Node) ast.Visitor {
	file, ok := node.(*ast.File)
	if !ok {
		return nil
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
				var isInput bool
				var operatorList []*api.OperatorMethod
				if ops, ok := m.inputModelReferer[typeSpec.Name.Name]; ok {
					operatorList = ops
					isInput = true
				} else if ops, ok := m.outputModelReferer[typeSpec.Name.Name]; ok {
					operatorList = ops
					isInput = false
				} else {
					continue
				}

				model := api.NewOperatorModel(typeSpec.Name.Name)
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range structType.Fields.List {
					var key, keyType string
					if field.Names != nil {
						key = field.Names[0].Name
					}
					switch field.Type.(type) {
					case *ast.Ident:
						keyType = field.Type.(*ast.Ident).Name
					case *ast.SelectorExpr:
						selectorExpr := field.Type.(*ast.SelectorExpr)
						switch selectorExpr.X.(type) {
						case *ast.Ident:
							keyType = selectorExpr.X.(*ast.Ident).Name
							if path, ok := importPath[keyType]; ok {
								model.AddImport(path, keyType)
							}
						}
						if keyType == "" {
							keyType = selectorExpr.Sel.Name
						} else {
							keyType = strings.Join([]string{keyType, selectorExpr.Sel.Name}, ".")
						}
					}

					model.AddField(key, keyType)
				}

				for _, op := range operatorList {
					if isInput {
						op.AddInput(model)
					} else {
						op.AddOutput(model)
					}
				}
			}
		}
	}

	return m
}
