package scanner

import "strings"

func RetrievePackageName(path string) string {
	paths := strings.Split(path, "/")
	return strings.Replace(paths[len(paths)-1], "-", "_", -1)
}
