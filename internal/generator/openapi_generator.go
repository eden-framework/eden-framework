package generator

import (
	"encoding/json"
	"fmt"
	"github.com/go-courier/oas"
	"github.com/profzone/eden-framework/internal/generator/scanner"
	"github.com/profzone/eden-framework/internal/project"
	"github.com/profzone/eden-framework/pkg/packagex"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/types"
	"os"
	"path"
	"regexp"
	"strings"
)

type OpenApiGenerator struct {
	api           *oas.OpenAPI
	pkg           *packagex.Package
	routerScanner *scanner.RouterScanner
}

func NewOpenApiGenerator() *OpenApiGenerator {
	return &OpenApiGenerator{
		api: oas.NewOpenAPI(),
	}
}

func (a *OpenApiGenerator) Load(cwd string) {
	entryPath := path.Join(cwd, "cmd")
	_, err := os.Stat(entryPath)
	if err != nil {
		if !os.IsExist(err) {
			logrus.Panicf("entry path does not exist: %s", entryPath)
		}
	}
	pkg, err := packagex.Load(entryPath)
	if err != nil {
		logrus.Panic(err)
	}

	a.pkg = pkg

	proj := project.Project{}
	err = proj.UnmarshalFromFile(cwd, "")
	if err != nil {
		logrus.Panic(err)
	}

	a.routerScanner = scanner.NewRouterScanner(pkg)
}

func (a *OpenApiGenerator) Pick() {
	defer func() {
		a.routerScanner.OperatorScanner().BindSchemas(a.api)
	}()

	var routerVar = findRootRouter(a.pkg)
	if routerVar == nil {
		return
	}

	router := a.routerScanner.Router(routerVar)
	routes := router.Routes()
	operationIDs := map[string]*scanner.Route{}
	for _, r := range routes {
		method := r.Method()
		operation := a.OperationByOperatorTypes(method, r.Operators...)
		if _, exists := operationIDs[operation.OperationId]; exists {
			panic(fmt.Errorf("operationID %s should be unique", operation.OperationId))
		}
		operationIDs[operation.OperationId] = r
		a.api.AddOperation(oas.HttpMethod(strings.ToLower(method)), a.patchPath(r.Path(), operation), operation)
	}
}

func (a *OpenApiGenerator) OperationByOperatorTypes(method string, operatorTypes ...*scanner.OperatorWithTypeName) *oas.Operation {
	operation := &oas.Operation{}

	length := len(operatorTypes)

	for idx := range operatorTypes {
		operatorTypes[idx].BindOperation(method, operation, idx == length-1)
	}

	return operation
}

var reHttpRouterPath = regexp.MustCompile("/:([^/]+)")

func (a *OpenApiGenerator) patchPath(openapiPath string, operation *oas.Operation) string {
	return reHttpRouterPath.ReplaceAllStringFunc(openapiPath, func(str string) string {
		name := reHttpRouterPath.FindAllStringSubmatch(str, -1)[0][1]

		var isParameterDefined = false

		for _, parameter := range operation.Parameters {
			if parameter.In == "path" && parameter.Name == name {
				isParameterDefined = true
			}
		}

		if isParameterDefined {
			return "/{" + name + "}"
		}

		return "/0"
	})
}

func (a *OpenApiGenerator) Output(outputPath string) Outputs {
	data, err := json.MarshalIndent(a.api, "", "    ")
	if err != nil {
		logrus.Panic(err)
	}
	return Outputs{
		path.Join(outputPath, "openapi.json"): string(data),
	}
}

func runnerFunc(node ast.Node) (runner *ast.FuncDecl) {
	switch n := node.(type) {
	case *ast.CallExpr:
		if len(n.Args) > 0 {
			if selectorExpr, ok := n.Fun.(*ast.SelectorExpr); ok {
				if selectorExpr.Sel.Name == "NewApplication" {
					switch node := n.Args[0].(type) {
					case *ast.SelectorExpr:
						runner = node.Sel.Obj.Decl.(*ast.FuncDecl)
					case *ast.Ident:
						runner = node.Obj.Decl.(*ast.FuncDecl)
					case *ast.FuncLit:
						funcDec := &ast.FuncDecl{
							Doc:  nil,
							Recv: nil,
							Name: nil,
							Type: node.Type,
							Body: node.Body,
						}
						runner = funcDec
					}
					return
				}
			}
		}
	}
	return nil
}

func rootRouter(node ast.Node, p *packagex.Package) *types.Var {
	switch n := node.(type) {
	case *ast.CallExpr:
		if len(n.Args) > 0 {
			if selectorExpr, ok := n.Fun.(*ast.SelectorExpr); ok {
				if selectorExpr.Sel.Name == "Serve" {
					switch node := n.Args[0].(type) {
					case *ast.SelectorExpr:
						return p.TypesInfo.ObjectOf(node.Sel).(*types.Var)
					case *ast.Ident:
						return p.TypesInfo.ObjectOf(node).(*types.Var)
					}
				}
			}
		}
	}
	return nil
}

func findRootRouter(p *packagex.Package) (router *types.Var) {
	for ident, def := range p.TypesInfo.Defs {
		if typFunc, ok := def.(*types.Func); ok {
			// 搜寻main函数
			if typFunc.Name() != "main" {
				continue
			}

			// 搜寻runner方法
			var runner *ast.FuncDecl
			ast.Inspect(ident.Obj.Decl.(*ast.FuncDecl), func(node ast.Node) bool {
				runnerDecl := runnerFunc(node)
				if runnerDecl != nil {
					runner = runnerDecl
					return false
				}
				return true
			})

			// 搜寻router入口
			if runner != nil {
				ast.Inspect(runner, func(node ast.Node) bool {
					if routerVar := rootRouter(node, p); routerVar != nil {
						router = routerVar
						return false
					}
					return true
				})
			}
			return
		}
	}
	return
}
