package scanner

import (
	"github.com/eden-framework/eden-framework/pkg/courier"
	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
	"github.com/eden-framework/eden-framework/pkg/courier/transport_http"
	"github.com/eden-framework/eden-framework/pkg/reflectx"
	"github.com/go-courier/oas"
	"go/ast"
	"go/constant"
	"go/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	XID           = "x-id"
	XGoVendorType = `x-go-vendor-type`
	XGoStarLevel  = `x-go-star-level`
	XGoFieldName  = `x-go-field-name`

	XTagValidate = `x-tag-validate`
	XTagMime     = `x-tag-mime`
	XTagJSON     = `x-tag-json`
	XTagXML      = `x-tag-xml`
	XTagName     = `x-tag-name`

	XEnumOptions = `x-enum-options`
	XStatusErrs  = `x-status-errors`
)

var (
	pkgImportPathHttpTransport = reflectx.ImportGoPath(reflect.TypeOf(transport_http.HttpRouteMeta{}).PkgPath())
	pkgImportPathHttpx         = reflectx.ImportGoPath(reflect.TypeOf(httpx.MethodGet{}).PkgPath())
	pkgImportPathCourier       = reflectx.ImportGoPath(reflect.TypeOf(courier.Router{}).PkgPath())
)

var (
	rxEnum   = regexp.MustCompile(`api:enum`)
	rxStrFmt = regexp.MustCompile(`api:stringFormat\s+(\S+)([\s\S]+)?$`)
)

var positionOrders = map[oas.Position]string{
	"path":   "1",
	"header": "2",
	"query":  "3",
	"cookie": "4",
}

func valueOf(v constant.Value) interface{} {
	if v == nil {
		return nil
	}

	switch v.Kind() {
	case constant.Float:
		v, _ := strconv.ParseFloat(v.String(), 10)
		return v
	case constant.Bool:
		v, _ := strconv.ParseBool(v.String())
		return v
	case constant.String:
		v, _ := strconv.Unquote(v.String())
		return v
	case constant.Int:
		v, _ := strconv.ParseInt(v.String(), 10, 64)
		return v
	}

	return nil
}

func isRouterType(typ types.Type) bool {
	return strings.HasSuffix(typ.String(), pkgImportPathCourier+".Router")
}

func isFromHttpTransport(typ types.Type) bool {
	return strings.Contains(typ.String(), pkgImportPathHttpTransport+".")
}

func filterMarkedLines(comments []string) []string {
	lines := make([]string, 0)
	for _, line := range comments {
		if !strings.HasPrefix(line, "@") {
			lines = append(lines, line)
		}
	}
	return lines
}

func tagValueAndFlagsByTagString(tagString string) (string, map[string]bool) {
	valueAndFlags := strings.Split(tagString, ",")
	v := valueAndFlags[0]
	tagFlags := map[string]bool{}
	if len(valueAndFlags) > 1 {
		for _, flag := range valueAndFlags[1:] {
			tagFlags[flag] = true
		}
	}
	return v, tagFlags
}

func dropMarkedLines(lines []string) string {
	return strings.Join(filterMarkedLines(lines), "\n")
}

func fullTypeName(typeName *types.TypeName) string {
	pkg := typeName.Pkg()
	if pkg != nil {
		return pkg.Path() + "." + typeName.Name()
	}
	return typeName.Name()
}

func ParseEnum(doc string) (string, bool) {
	if rxEnum.MatchString(doc) {
		return strings.TrimSpace(strings.Replace(doc, "api:enum", "", -1)), true
	}
	return doc, false
}

func ParseType(typeExpr ast.Expr) (keyType, pkgName string, pointer bool) {
	switch typeExpr.(type) {
	case *ast.Ident:
		keyType = typeExpr.(*ast.Ident).Name
	case *ast.StarExpr:
		starExpr := typeExpr.(*ast.StarExpr)
		keyType, pkgName, _ = ParseType(starExpr.X)
		keyType = "*" + keyType
		pointer = true
	case *ast.SelectorExpr:
		selectorExpr := typeExpr.(*ast.SelectorExpr)
		pkgName, _, pointer = ParseType(selectorExpr.X)
		keyType = selectorExpr.Sel.Name
	case *ast.ArrayType:
		arrayType := typeExpr.(*ast.ArrayType)
		keyType, pkgName, pointer = ParseType(arrayType.Elt)
		keyType = "[]" + keyType
	}

	return
}

func ParseStringFormat(doc string) (string, string) {
	matched := rxStrFmt.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

func RetrievePackageName(path string) string {
	path = strings.Trim(path, "\"")
	paths := strings.Split(path, "/")
	return strings.Replace(paths[len(paths)-1], "-", "_", -1)
}
