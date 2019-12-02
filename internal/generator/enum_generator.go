package generator

import (
	"github.com/profzone/eden-framework/internal/generator/files"
	"github.com/profzone/eden-framework/internal/generator/importer"
	"github.com/profzone/eden-framework/internal/generator/scanner"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/build"
	"golang.org/x/tools/go/packages"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type EnumGenerator struct {
	pkgs        []*packages.Package
	EnumScanner *scanner.EnumScanner
	TypeName    string
}

func NewEnumGenerator(scanner *scanner.EnumScanner, typeName string) *EnumGenerator {
	return &EnumGenerator{
		EnumScanner: scanner,
		TypeName:    typeName,
	}
}

func (e *EnumGenerator) Load(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			logrus.Panicf("entry path does not exist: %s", path)
		}
	}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedDeps | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes,
		Dir:  path,
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		logrus.Panic(err)
	}

	errs := packages.PrintErrors(pkgs)
	if errs > 0 {
		logrus.Panicf("packages.PrintErrors(a.pkgs) = %d", errs)
	}

	e.pkgs = pkgs
}

func (e *EnumGenerator) Pick() {
	packages.Visit(e.pkgs, nil, func(i *packages.Package) {
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
							e.EnumScanner.Enum(strings.Join([]string{i.PkgPath, typeSpec.Name.Name}, "."))
						}
					} else {
						e.EnumScanner.Enum(strings.Join([]string{i.PkgPath, typeSpec.Name.Name}, "."))
					}
				}

				return true
			})
		}
	})
}

func (e *EnumGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}
	for typeFullName, enum := range e.EnumScanner.Enums {
		pkgPath, typeName := importer.GetPackagePathAndDecl(typeFullName)

		p, _ := build.Import(pkgPath, "", build.FindOnly)
		dir, _ := filepath.Rel(outputPath, p.Dir)

		enum := files.NewEnum(pkgPath, scanner.RetrievePackageName(pkgPath), typeName, enum, e.EnumScanner.HasOffset(typeFullName))
		outputs.Add(GeneratedSuffix(path.Join(dir, str.ToLowerSnakeCase(typeName)+".go")), enum.String())
	}
	return outputs
}
