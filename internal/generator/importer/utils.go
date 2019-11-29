package importer

import "strings"

func GetPackagePathAndDecl(path string) (importPath, decl string) {
	slash := strings.LastIndex(path, "/")
	dot := strings.LastIndex(path, ".")
	if dot > slash {
		return path[0:dot], path[dot+1:]
	}

	return path, ""
}

func ParseVendor(path string) string {
	paths := strings.Split(path, "/vendor/")
	return paths[len(paths)-1]
}
