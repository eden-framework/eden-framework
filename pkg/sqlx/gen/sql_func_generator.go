package gen

import (
	"go/build"
	"go/parser"
	"go/types"
	"path"

	"golang.org/x/tools/go/loader"

	"github.com/profzone/eden-framework/pkg/codegen"
	"github.com/profzone/eden-framework/pkg/codegen/loaderx"
)

type Config struct {
	StructName          string
	TableName           string
	Database            string
	WithComments        bool
	WithTableInterfaces bool

	FieldPrimaryKey string
	FieldSoftDelete string
	FieldCreatedAt  string
	FieldUpdatedAt  string

	ConstSoftDeleteTrue  string
	ConstSoftDeleteFalse string
}

func (g *Config) Defaults() {
	if g.FieldSoftDelete == "" {
		g.FieldSoftDelete = "Enabled"
	}

	if g.FieldCreatedAt == "" {
		g.FieldCreatedAt = "CreateTime"
	}

	if g.FieldUpdatedAt == "" {
		g.FieldUpdatedAt = "UpdateTime"
	}

	if g.ConstSoftDeleteTrue == "" {
		g.ConstSoftDeleteTrue = "github.com/profzone/eden-framework/pkg/courier/enumeration.BOOL__TRUE"
	}

	if g.ConstSoftDeleteFalse == "" {
		g.ConstSoftDeleteFalse = "github.com/profzone/eden-framework/pkg/courier/enumeration.BOOL__FALSE"
	}

	if g.TableName == "" {
		g.TableName = toDefaultTableName(g.StructName)
	}
}

type SqlFuncGenerator struct {
	Config
	pkgImportPath string
	program       *loader.Program
	model         *Model
}

func (g *SqlFuncGenerator) Load(cwd string) {
	ldr := loader.Config{
		AllowErrors: true,
		ParserMode:  parser.ParseComments,
	}

	pkgImportPath := codegen.GetPackageImportPath(cwd)
	ldr.Import(pkgImportPath)

	p, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	g.pkgImportPath = pkgImportPath
	g.program = p

	g.Defaults()
}

func (g *SqlFuncGenerator) Pick() {
	for pkg, pkgInfo := range g.program.AllPackages {
		if pkg.Path() != g.pkgImportPath {
			continue
		}
		for ident, obj := range pkgInfo.Defs {
			if typeName, ok := obj.(*types.TypeName); ok {
				if typeName.Name() == g.StructName {
					if _, ok := typeName.Type().Underlying().(*types.Struct); ok {
						comments := loaderx.CommentsOf(g.program.Fset, ident, pkgInfo.Files...)
						g.model = NewModel(g.program, typeName, comments, &g.Config)
					}
				}

			}
		}
	}
}

func (g *SqlFuncGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}

	if g.model != nil {
		pkg, _ := build.Import(g.model.TypeName.Pkg().Path(), "", build.FindOnly)
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkg.Dir, codegen.ToLowerSnakeCase(g.StructName)+".go")), g.model.Render())
	}

	return outputs
}
