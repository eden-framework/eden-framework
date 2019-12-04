package scanner

import (
	"github.com/profzone/eden-framework/internal/generator/api"
	"github.com/profzone/eden-framework/internal/generator/importer"
	"github.com/profzone/eden-framework/pkg/enumeration"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/constant"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strconv"
	"strings"
)

type EnumScanner struct {
	Enums map[string]api.Enum
}

func NewEnumScanner() *EnumScanner {
	return &EnumScanner{
		Enums: make(map[string]api.Enum),
	}
}

func (s *EnumScanner) HasOffset(typeFullName string) bool {
	pkgPath, typeName := importer.GetPackagePathAndDecl(typeFullName)
	if typeName == "" {
		logrus.Panic("typeFullName must have path and typeName")
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax,
	}, pkgPath)
	if err != nil {
		logrus.Panic(err)
	}
	pkg := pkgs[0]
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range gen.Specs {
				if value, ok := spec.(*ast.ValueSpec); ok {
					if value.Names == nil {
						continue
					}
					if value.Names[0].Name == str.ToUpperSnakeCase(typeName)+"_OFFSET" {
						return true
					}
				}
			}
		}
	}

	return false
}

func (s *EnumScanner) Enum(typeFullName string) api.Enum {
	if enumOptions, ok := s.Enums[typeFullName]; ok {
		return enumOptions.Sort()
	}

	pkgPath, typeName := importer.GetPackagePathAndDecl(typeFullName)
	if typeName == "" {
		logrus.Panic("typeFullName must have path and typeName")
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedDeps,
	}, pkgPath)
	if err != nil {
		logrus.Panic(err)
	}
	errCount := packages.PrintErrors(pkgs)
	if errCount > 0 {
		logrus.Panicf("packages.Load errors count: %d", errCount)
	}
	pkg := pkgs[0]

	for ident, def := range pkg.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type().String() != typeFullName {
			continue
		}
		name := typeConst.Name()
		if strings.HasSuffix(name, "_OFFSET") {
			continue
		}
		if strings.HasPrefix(name, str.ToUpperSnakeCase(typeName)+"_") {
			val := typeConst.Val()
			label := strings.TrimSpace(ident.Obj.Decl.(*ast.ValueSpec).Comment.Text())
			values := strings.SplitN(name, "__", 2)
			if len(values) == 2 {
				value := values[1]
				typeFullName := strings.Join([]string{pkg.PkgPath, typeName}, ".")
				s.addEnum(typeFullName, value, getConstVal(val), label)
			}
		}
	}

	return s.Enums[typeFullName].Sort()
}

func (s *EnumScanner) addEnum(typeFullName string, value interface{}, val interface{}, label string) {
	s.Enums[typeFullName] = append(s.Enums[typeFullName], enumeration.EnumOption{
		Val:   val,
		Value: value,
		Label: label,
	})
}

func getConstVal(constVal constant.Value) interface{} {
	switch constVal.Kind() {
	case constant.String:
		stringVal, _ := strconv.Unquote(constVal.String())
		return stringVal
	case constant.Int:
		intVal, _ := strconv.ParseInt(constVal.String(), 10, 64)
		return intVal
	case constant.Float:
		floatVal, _ := strconv.ParseFloat(constVal.String(), 10)
		return floatVal
	}
	return nil
}
