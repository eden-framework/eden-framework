package importer

type PackageImporter struct {
	PackagePath string
	pkgs        map[string]*Package
}

func NewPackageImporter(packagePath string) *PackageImporter {
	return &PackageImporter{
		PackagePath: packagePath,
		pkgs:        make(map[string]*Package),
	}
}

func (i *PackageImporter) Import(importPath string, useAlias bool) *Package {
	importPath = parseVendor(importPath)
	p, ok := i.pkgs[importPath]
	if !ok {

	}

	return p
}
