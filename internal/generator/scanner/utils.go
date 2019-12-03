package scanner

import (
	"go/ast"
	"regexp"
	"strings"
)

var (
	rxEnum   = regexp.MustCompile(`api:enum`)
	rxStrFmt = regexp.MustCompile(`api:stringFormat\s+(\S+)([\s\S]+)?$`)
)

func ParseEnum(doc string) (string, bool) {
	if rxEnum.MatchString(doc) {
		return strings.TrimSpace(strings.Replace(doc, "api:enum", "", -1)), true
	}
	return doc, false
}

func ParseType(typeExpr ast.Expr) (keyType, pkgName string) {
	switch typeExpr.(type) {
	case *ast.Ident:
		keyType = typeExpr.(*ast.Ident).Name
	case *ast.StarExpr:
		starExpr := typeExpr.(*ast.StarExpr)
		keyType, pkgName = ParseType(starExpr.X)
	case *ast.SelectorExpr:
		selectorExpr := typeExpr.(*ast.SelectorExpr)
		pkgName, _ = ParseType(selectorExpr.X)
		keyType = selectorExpr.Sel.Name
	case *ast.ArrayType:
		arrayType := typeExpr.(*ast.ArrayType)
		keyType, pkgName = ParseType(arrayType.Elt)
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
