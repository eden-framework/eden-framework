package scanner

import (
	"fmt"
	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
	"github.com/eden-framework/eden-framework/pkg/packagex"
	"github.com/eden-framework/eden-framework/pkg/reflectx"
	"go/ast"
	"go/types"
	"sort"
	"strconv"
	"strings"
)

func NewStatusErrorScanner(pkg *packagex.Package) *StatusErrorScanner {
	return &StatusErrorScanner{
		pkg: pkg,
	}
}

type StatusErrorScanner struct {
	pkg          *packagex.Package
	StatusErrors map[*types.TypeName][]*status_error.StatusError
}

func sortedStatusErrList(list []*status_error.StatusError) []*status_error.StatusError {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Code < list[j].Code
	})
	return list
}

func (scanner *StatusErrorScanner) StatusError(typeName *types.TypeName) []*status_error.StatusError {
	if typeName == nil {
		return nil
	}

	if statusErrs, ok := scanner.StatusErrors[typeName]; ok {
		return sortedStatusErrList(statusErrs)
	}

	if !strings.Contains(typeName.Type().Underlying().String(), "int") {
		panic(fmt.Errorf("status error type underlying must be an int or uint, but got %s", typeName.String()))
	}

	pkgInfo := scanner.pkg.Pkg(typeName.Pkg().Path())
	if pkgInfo == nil {
		return nil
	}

	var serviceCode int64 = 0

	method, ok := reflectx.FromTType(typeName.Type()).MethodByName("ServiceCode")
	if ok {
		results, n := scanner.pkg.FuncResultsOf(method.(*reflectx.TMethod).Func)
		if n == 1 {
			ret := results[0][0]
			if ret.IsValue() {
				if i, err := strconv.ParseInt(ret.Value.String(), 10, 64); err == nil {
					serviceCode = i
				}
			}
		}
	}

	for ident, def := range pkgInfo.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type() != typeName.Type() {
			continue
		}

		key := typeConst.Name()
		code, _ := strconv.ParseInt(typeConst.Val().String(), 10, 64)

		msg, canBeTalkError := ParseStatusErrMsg(ident.Obj.Decl.(*ast.ValueSpec).Doc.Text())

		scanner.addStatusError(typeName, key, code+serviceCode, msg, canBeTalkError)
	}

	return sortedStatusErrList(scanner.StatusErrors[typeName])
}

func ParseStatusErrMsg(s string) (string, bool) {
	firstLine := strings.Split(strings.TrimSpace(s), "\n")[0]

	prefix := "@errTalk "
	if strings.HasPrefix(firstLine, prefix) {
		return firstLine[len(prefix):], true
	}
	return firstLine, false
}

func (scanner *StatusErrorScanner) addStatusError(
	typeName *types.TypeName,
	key string, code int64, msg string, canBeTalkError bool,
) {
	if scanner.StatusErrors == nil {
		scanner.StatusErrors = map[*types.TypeName][]*status_error.StatusError{}
	}

	statusErr := status_error.NewStatusError(key, code, msg)
	if canBeTalkError {
		statusErr = statusErr.WithErrTalk()
	}
	scanner.StatusErrors[typeName] = append(scanner.StatusErrors[typeName], statusErr)
}
