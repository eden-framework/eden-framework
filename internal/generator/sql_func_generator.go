package generator

import (
	"github.com/eden-framework/codegen"
	"github.com/eden-framework/packagex"
	"github.com/eden-framework/sqlx/generator"
	"github.com/sirupsen/logrus"
	"go/types"
	"os"
	"path"
	"path/filepath"
)

func NewSqlFuncGenerator() *SqlFuncGenerator {
	return &SqlFuncGenerator{}
}

type SqlFuncGenerator struct {
	generator.Config
	pkg   *packagex.Package
	model *generator.Model
}

func (g *SqlFuncGenerator) Load(cwd string) {
	var err error
	if len(cwd) == 0 {
		cwd, err = os.Getwd()
		if err != nil {
			logrus.Panicf("get current working directory err: %v, cwd: %s", err, cwd)
		}
	}
	_, err = os.Stat(cwd)
	if err != nil {
		if !os.IsExist(err) {
			logrus.Panicf("entry path does not exist: %s", cwd)
		}
	}
	pkg, err := packagex.Load(cwd)
	if err != nil {
		logrus.Panic(err)
	}

	g.pkg = pkg

}

func (g *SqlFuncGenerator) Pick() {
	for ident, obj := range g.pkg.TypesInfo.Defs {
		if typeName, ok := obj.(*types.TypeName); ok {
			if typeName.Name() == g.StructName {
				if _, ok := typeName.Type().Underlying().(*types.Struct); ok {
					g.model = generator.NewModel(g.pkg, typeName, g.pkg.CommentsOf(ident), &g.Config)
				}
			}
		}
	}
}

func (g *SqlFuncGenerator) Output(outputPath string) Outputs {
	if g.model == nil {
		return nil
	}

	dir, _ := filepath.Rel(outputPath, filepath.Dir(g.pkg.GoFiles[0]))
	filename := codegen.GeneratedFileSuffix(codegen.LowerSnakeCase(g.StructName) + ".go")

	file := codegen.NewFile(g.pkg.Name, filename)
	g.model.WriteTo(file)

	return Outputs{
		path.Join(dir, filename): string(file.Bytes()),
	}
}

type Config struct {
	StructName string
	TableName  string
	Database   string

	WithComments        bool
	WithTableName       bool
	WithTableInterfaces bool
	WithMethods         bool

	FieldPrimaryKey   string
	FieldKeyDeletedAt string
	FieldKeyCreatedAt string
	FieldKeyUpdatedAt string
}

func (g *Config) SetDefaults() {
	if g.FieldKeyDeletedAt == "" {
		g.FieldKeyDeletedAt = "DeletedAt"
	}

	if g.FieldKeyCreatedAt == "" {
		g.FieldKeyCreatedAt = "CreatedAt"
	}

	if g.FieldKeyUpdatedAt == "" {
		g.FieldKeyUpdatedAt = "UpdatedAt"
	}

	if g.TableName == "" {
		g.TableName = toDefaultTableName(g.StructName)
	}
}
