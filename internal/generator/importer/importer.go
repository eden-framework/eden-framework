package importer

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
	"io"
	str "profzone/eden-framework/pkg/strings"
	"strings"
)

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
	importPath = ParseVendor(importPath)
	p, ok := i.pkgs[importPath]
	if !ok {
		pkgs, err := packages.Load(nil, importPath)
		if err != nil {
			logrus.Panic(err)
		}
		p = &Package{
			Package: pkgs[0],
		}
		if useAlias {
			p.Alias = str.ToLowerSnakeCase(importPath)
		}
		i.pkgs[importPath] = p
	}

	return p
}

func (i *PackageImporter) Use(importPath string, subs ...string) string {
	importPath, decl := getPackagePathAndDecl(strings.Join(append([]string{importPath}, subs...), "/"))
	p := i.Import(importPath, true)

	if decl != "" {
		if importPath == i.PackagePath {
			return decl
		}

		return fmt.Sprintf("%s.%s", p.GetSelector(), decl)
	}

	return p.GetSelector()
}

func (i *PackageImporter) WriteToImports(w io.Writer) {
	if len(i.pkgs) > 0 {
		for _, importPkg := range i.pkgs {
			io.WriteString(w, importPkg.String()+"\n")
		}
	}
}

func (i *PackageImporter) String() string {
	buf := &bytes.Buffer{}
	if len(i.pkgs) > 0 {
		buf.WriteString("import (\n")
		i.WriteToImports(buf)
		buf.WriteString(")\n\n")
	}
	return buf.String()
}
