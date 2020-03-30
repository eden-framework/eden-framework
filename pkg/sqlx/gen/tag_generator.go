package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"

	"github.com/profzone/eden-framework/pkg/codegen"
	"github.com/profzone/eden-framework/pkg/codegen/loaderx"
)

type TagGenerator struct {
	StructNames   []string
	pkgImportPath string
	WithDefaults  bool
	program       *loader.Program
	outputs       codegen.Outputs
}

func (g *TagGenerator) Load(cwd string) {
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
	g.outputs = codegen.Outputs{}
}

func (g *TagGenerator) Pick() {
	for pkg, pkgInfo := range g.program.AllPackages {
		if pkg.Path() != g.pkgImportPath {
			continue
		}
		for ident, obj := range pkgInfo.Defs {
			if typeName, ok := obj.(*types.TypeName); ok {
				for _, structName := range g.StructNames {
					if typeName.Name() == structName {
						if typeStruct, ok := typeName.Type().Underlying().(*types.Struct); ok {
							modifyTag(ident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType), typeStruct, g.WithDefaults)
							file := loaderx.FileOf(ident, pkgInfo.Files...)
							g.outputs.Add(g.program.Fset.Position(file.Pos()).Filename, loaderx.StringifyAst(g.program.Fset, file))
						}
					}
				}
			}
		}
	}
}

func toTags(tags map[string]string) (tag string) {
	names := make([]string, 0)
	for name := range tags {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		tag += fmt.Sprintf("%s:%s ", name, strconv.Quote(tags[name]))
	}
	return strings.TrimSpace(tag)
}

func getTags(tag string) (tags map[string]string) {
	tags = make(map[string]string)
	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		tags[name] = value

	}
	return
}

func modifyTag(structType *ast.StructType, typeStruct *types.Struct, withDefaults bool) {
	for i := 0; i < typeStruct.NumFields(); i++ {
		f := typeStruct.Field(i)
		if f.Anonymous() {
			continue
		}
		tags := getTags(typeStruct.Tag(i))
		astField := structType.Fields.List[i]

		if tags["db"] == "" {
			tags["db"] = fmt.Sprintf("F_%s", codegen.ToLowerSnakeCase(f.Name()))
		}
		if tags["json"] == "" {
			tags["json"] = codegen.ToLowerCamelCase(f.Name())
			switch f.Type().(type) {
			case *types.Basic:
				if f.Type().(*types.Basic).Kind() == types.Uint64 {
					tags["json"] = tags["json"] + ",string"
				}
			}
		}
		if tags["sql"] == "" {
			tpe := f.Type()
			switch codegen.DeVendor(tpe.String()) {
			case "github.com/profzone/eden-framework/pkg/timelib.MySQLDatetime":
				tags["sql"] = "datetime NOT NULL"
			case "github.com/profzone/eden-framework/pkg/timelib.MySQLTimestamp":
				tags["sql"] = toSqlFromKind(types.Typ[types.Int64].Kind(), withDefaults)
			default:
				tpe, err := IndirectType(tpe)
				if err != nil {
					logrus.Warnf("%s, make sure type of Field `%s` have sql.Valuer and sql.Scanner interface", err, f.Name())
				}
				switch tpe.(type) {
				case *types.Basic:
					tags["sql"] = toSqlFromKind(tpe.(*types.Basic).Kind(), withDefaults)
				default:
					tags["sql"] = WithDefaults("varchar(255) NOT NULL", withDefaults, "")
				}
			}
		}
		astField.Tag = &ast.BasicLit{Kind: token.STRING, Value: "`" + toTags(tags) + "`"}
	}
}

func IndirectType(tpe types.Type) (types.Type, error) {
	switch tpe.(type) {
	case *types.Basic:
		return tpe.(*types.Basic), nil
	case *types.Struct, *types.Slice, *types.Array, *types.Map:
		return nil, fmt.Errorf("unsupport type %s", tpe)
	case *types.Pointer:
		return IndirectType(tpe.(*types.Pointer).Elem())
	default:
		return IndirectType(tpe.Underlying())
	}
}

func WithDefaults(dataType string, withDefaults bool, defaultValue string) string {
	if withDefaults {
		return dataType + fmt.Sprintf(" DEFAULT '%s'", defaultValue)
	}
	return dataType
}

func toSqlFromKind(kind types.BasicKind, withDefaults bool) string {
	switch kind {
	case types.Bool:
		return WithDefaults("tinyint(1) NOT NULL", withDefaults, "0")
	case types.Int8:
		return WithDefaults("tinyint NOT NULL", withDefaults, "0")
	case types.Int16:
		return WithDefaults("smallint NOT NULL", withDefaults, "0")
	case types.Int, types.Int32:
		return WithDefaults("int NOT NULL", withDefaults, "0")
	case types.Int64:
		return WithDefaults("bigint NOT NULL", withDefaults, "0")
	case types.Uint8:
		return WithDefaults("tinyint unsigned NOT NULL", withDefaults, "0")
	case types.Uint16:
		return WithDefaults("smallint unsigned NOT NULL", withDefaults, "0")
	case types.Uint, types.Uint32:
		return WithDefaults("int unsigned NOT NULL", withDefaults, "0")
	case types.Uint64:
		return WithDefaults("bigint unsigned NOT NULL", withDefaults, "0")
	case types.Float32:
		return WithDefaults("float NOT NULL", withDefaults, "0")
	case types.Float64:
		return WithDefaults("double NOT NULL", withDefaults, "0")
	default:
		// string
		return WithDefaults("varchar(255) NOT NULL", withDefaults, "")
	}
}

func (g *TagGenerator) Output(cwd string) codegen.Outputs {
	return g.outputs
}
