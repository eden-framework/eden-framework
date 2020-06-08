package openapi_scanner

import (
	"github.com/go-courier/oas"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"github.com/profzone/eden-framework/pkg/courier/transport_http"
	"github.com/profzone/eden-framework/pkg/reflectx"
	"go/constant"
	"go/types"
	"reflect"
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
