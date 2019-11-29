package scanner

import (
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/constant"
	"go/types"
	"golang.org/x/tools/go/packages"
	"profzone/eden-framework/internal/generator/importer"
	"profzone/eden-framework/pkg/enumeration"
	str "profzone/eden-framework/pkg/strings"
	"sort"
	"strconv"
	"strings"
)

type EnumScanner struct {
	Enums map[string]Enum
}

func NewEnumScanner() *EnumScanner {
	return &EnumScanner{
		Enums: make(map[string]Enum),
	}
}

func (s *EnumScanner) HasOffset(typeFullName string) bool {
	pkgPath, typeName := importer.GetPackagePathAndDecl(typeFullName)
	if typeName == "" {
		logrus.Panic("typeFullName must have path and typeName")
	}
	pkgs, err := packages.Load(nil, pkgPath)
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

func (s *EnumScanner) Enum(typeFullName string) Enum {
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
		if strings.HasPrefix(name, strings.ToUpper(typeName)+"_") {
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

type Enum enumeration.Enum

func (enum Enum) Sort() Enum {
	sort.Slice(enum, func(i, j int) bool {
		switch enum[i].Value.(type) {
		case string:
			return enum[i].Value.(string) < enum[j].Value.(string)
		case int64:
			return enum[i].Value.(int64) < enum[j].Value.(int64)
		case float64:
			return enum[i].Value.(float64) < enum[j].Value.(float64)
		}
		return true
	})
	return enum
}

func (enum Enum) Labels() (labels []string) {
	for _, e := range enum {
		labels = append(labels, e.Label)
	}
	return
}

func (enum Enum) Vals() (vals []interface{}) {
	for _, e := range enum {
		vals = append(vals, e.Val)
	}
	return
}

func (enum Enum) Values() (values []interface{}) {
	for _, e := range enum {
		values = append(values, e.Value)
	}
	return
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
