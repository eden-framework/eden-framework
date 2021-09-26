package importer

import (
	"strconv"
)

type Package struct {
	Name    string
	Alias   string
	PkgPath string
}

func (p *Package) GetSelector() string {
	if p.Alias != "" {
		return p.Alias
	}

	return p.Name
}

func (p *Package) String() string {
	if p.Alias != "" {
		pkgName := RetrievePackageName(p.PkgPath)
		if p.Alias != pkgName {
			return p.Alias + " " + strconv.Quote(p.PkgPath)
		}
	}

	return strconv.Quote(p.PkgPath)
}
