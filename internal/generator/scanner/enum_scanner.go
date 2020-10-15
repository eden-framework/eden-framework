package scanner

import (
	"fmt"
	"github.com/eden-framework/enumeration"
	"github.com/eden-framework/packagex"
	str "github.com/eden-framework/strings"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"
)

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

func NewEnumScanner(pkg *packagex.Package) *EnumScanner {
	return &EnumScanner{
		pkg: pkg,
	}
}

type EnumScanner struct {
	pkg     *packagex.Package
	EnumSet map[*types.TypeName][]enumeration.EnumOption
}

func sortEnumOptions(enumOptions []enumeration.EnumOption) []enumeration.EnumOption {
	sort.Slice(enumOptions, func(i, j int) bool {
		return enumOptions[i].Val > enumOptions[j].Val
	})
	return enumOptions
}

func (scanner *EnumScanner) Enum(typeName *types.TypeName) []enumeration.EnumOption {
	if typeName == nil {
		return nil
	}

	if enumOptions, ok := scanner.EnumSet[typeName]; ok {
		return sortEnumOptions(enumOptions)
	}

	if !strings.Contains(typeName.Type().Underlying().String(), "int") {
		panic(fmt.Errorf("enumeration type underlying must be an int or uint, but got %s", typeName.String()))
	}

	pkgInfo := scanner.pkg.Pkg(typeName.Pkg().Path())
	if pkgInfo == nil {
		return nil
	}

	typeNameString := typeName.Name()

	for ident, def := range pkgInfo.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type() != typeName.Type() {
			continue
		}
		name := typeConst.Name()

		if strings.HasPrefix(name, "_") {
			continue
		}

		val := typeConst.Val()
		label := strings.TrimSpace(ident.Obj.Decl.(*ast.ValueSpec).Comment.Text())

		if strings.HasPrefix(name, str.ToUpperSnakeCase(typeNameString)) {
			var values = strings.SplitN(name, "__", 2)
			if len(values) == 2 {
				firstChar := values[1][0]
				if '0' <= firstChar && firstChar <= '9' {
					panic(fmt.Errorf("not support enum key starts with number, please fix `%s`", name))
				}
				intVal, _ := strconv.ParseInt(val.String(), 10, 64)
				scanner.addEnum(typeName, int(intVal), values[1], label)
			}
		}
	}

	return sortEnumOptions(scanner.EnumSet[typeName])
}

func (scanner *EnumScanner) addEnum(typeName *types.TypeName, constValue int, value string, label string) {
	if scanner.EnumSet == nil {
		scanner.EnumSet = map[*types.TypeName][]enumeration.EnumOption{}
	}

	if label == "" {
		label = value
	}

	scanner.EnumSet[typeName] = append(scanner.EnumSet[typeName], enumeration.EnumOption{
		Val:   constValue,
		Value: value,
		Label: label,
	})
}

func (scanner *EnumScanner) HasOffset(typeName *types.TypeName) bool {
	pkgInfo := scanner.pkg.PkgInfoOf(typeName)
	if pkgInfo == nil {
		return false
	}
	for _, def := range pkgInfo.Defs {
		if typeConst, ok := def.(*types.Const); ok {
			if typeConst.Name() == str.ToUpperSnakeCase(typeName.Name())+"_OFFSET" {
				return true
			}
		}
	}
	return false
}
