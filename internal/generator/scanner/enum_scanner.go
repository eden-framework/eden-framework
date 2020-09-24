package scanner

import (
	"fmt"
	"github.com/eden-framework/eden-framework/pkg/enumeration"
	"github.com/eden-framework/eden-framework/pkg/packagex"
	str "github.com/eden-framework/eden-framework/pkg/strings"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"
)

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
