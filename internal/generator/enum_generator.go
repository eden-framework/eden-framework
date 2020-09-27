package generator

import (
	"github.com/eden-framework/eden-framework/internal/generator/files"
	"github.com/eden-framework/eden-framework/internal/generator/scanner"
	"github.com/eden-framework/eden-framework/pkg/packagex"
	str "github.com/eden-framework/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/build"
	"golang.org/x/tools/go/packages"
	"os"
	"path"
	"path/filepath"
)

type EnumGenerator struct {
	pkg      *packagex.Package
	scanner  *scanner.EnumScanner
	TypeName string
}

func NewEnumGenerator(typeName string) *EnumGenerator {
	return &EnumGenerator{
		TypeName: typeName,
	}
}

func (e *EnumGenerator) Load(cwd string) {
	_, err := os.Stat(cwd)
	if err != nil {
		if !os.IsExist(err) {
			logrus.Panicf("entry path does not exist: %s", cwd)
		}
	}
	pkg, err := packagex.Load(cwd)
	if err != nil {
		logrus.Panic(err)
	}

	e.pkg = pkg
	e.scanner = scanner.NewEnumScanner(pkg)
}

func (e *EnumGenerator) Pick() {
	packages.Visit(e.pkg.AllPackages, nil, func(i *packages.Package) {
		for _, f := range i.Syntax {
			commentScanner := scanner.NewCommentScanner(i.Fset, f)
			ast.Inspect(f, func(node ast.Node) bool {
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}
				doc := commentScanner.CommentsOf(typeSpec)
				doc, hasEnum := scanner.ParseEnum(doc)

				if hasEnum {
					if e.TypeName != "" {
						if e.TypeName == typeSpec.Name.Name {
							e.scanner.Enum(e.pkg.TypeName(typeSpec.Name.Name))
						}
					} else {
						e.scanner.Enum(e.pkg.TypeName(typeSpec.Name.Name))
					}
				}

				return true
			})
		}
	})
}

func (e *EnumGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}
	for typeName, enum := range e.scanner.EnumSet {
		p, _ := build.Import(typeName.Pkg().Path(), "", build.FindOnly)
		dir, _ := filepath.Rel(outputPath, p.Dir)

		enum := files.NewEnum(typeName.Pkg().Path(), typeName.Pkg().Name(), typeName.Name(), enum, e.scanner.HasOffset(typeName))
		outputs.Add(GeneratedSuffix(path.Join(dir, str.ToLowerSnakeCase(typeName.Name())+".go")), enum.String())
	}
	return outputs
}
