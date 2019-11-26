package scanner

import "strings"

func RetrievePackageName(path string) string {
	path = strings.Trim(path, "\"")
	paths := strings.Split(path, "/")
	return strings.Replace(paths[len(paths)-1], "-", "_", -1)
}
