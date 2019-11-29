package scanner

import (
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
