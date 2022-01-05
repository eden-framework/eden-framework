package importer

import (
	"bytes"
	"fmt"
	str "gitee.com/eden-framework/strings"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
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

func (i *PackageImporter) mustGetPackage(packageList []string) {
	for _, p := range packageList {
		cmd := exec.Command("go", "get", p)
		logrus.Infof("command executing: %s", cmd.String())
		err := cmd.Run()
		if err != nil {
			logrus.Panicf("go get err: %v\n", err)
		}
	}
}

func (i *PackageImporter) AddImport(importPath string, p *Package) {
	if _, ok := i.pkgs[importPath]; ok {
		return
	}
	i.mustGetPackage([]string{importPath})
	i.pkgs[importPath] = p
}

func (i *PackageImporter) Import(pkgNamePattern, searchPath string, useAlias bool) *Package {
	pkgNamePattern = ParseVendor(pkgNamePattern)
	p, ok := i.pkgs[pkgNamePattern]
	if !ok {
		pkgName := RetrievePackageName(pkgNamePattern)
		p = &Package{
			Name:    pkgName,
			PkgPath: pkgNamePattern,
		}
		if useAlias {
			p.Alias = str.ToLowerSnakeCase(pkgNamePattern)
		}
		if searchPath == "" {
			i.mustGetPackage([]string{pkgNamePattern})
		}
		i.pkgs[pkgNamePattern] = p
	}

	return p
}

func (i *PackageImporter) Use(pkgName string, subs ...string) string {
	pkgName, decl := GetPackagePathAndDecl(strings.Join(append([]string{pkgName}, subs...), "/"))
	p := i.Import(pkgName, "", true)

	if decl != "" {
		if pkgName == i.PackagePath {
			return decl
		}

		return fmt.Sprintf("%s.%s", p.GetSelector(), decl)
	}

	return p.GetSelector()
}

func (i *PackageImporter) UseWithoutAlias(pkgName, searchPath string, subs ...string) string {
	pkgName, decl := GetPackagePathAndDecl(strings.Join(append([]string{pkgName}, subs...), "/"))
	p := i.Import(pkgName, searchPath, false)

	if decl != "" {
		if pkgName == i.PackagePath {
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

func (i *PackageImporter) Merge(target *PackageImporter) {
	for importPath, pkg := range target.pkgs {
		i.AddImport(importPath, pkg)
	}
}
