package importer

import "strings"

func getPackagePathAndDecl(path string) (importPath, decl string) {
	slash := strings.LastIndex(path, "/")
	dot := strings.LastIndex(path, ".")
	if dot > slash {
		return path[0:dot], path[dot+1:]
	}

	return path, ""
}

func parseVendor(path string) string {
	paths := strings.Split(path, "/vendor/")
	return paths[len(paths)-1]
}
