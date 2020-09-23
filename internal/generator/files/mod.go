package files

import (
	"bytes"
	"fmt"
)

type ModRequired struct {
	pkgName string
	version string
}

func NewModRequired(pkg, ver string) ModRequired {
	return ModRequired{
		pkgName: pkg,
		version: ver,
	}
}

func (r ModRequired) String() string {
	return fmt.Sprintf("%s %s", r.pkgName, r.version)
}

type ModReplace struct {
	fromPkgName    string
	fromPkgVersion string
	toPkgName      string
	toPkgVersion   string
}

func NewModReplace(fromPkg, fromVer, toPkg, toVer string) ModReplace {
	return ModReplace{
		fromPkgName:    fromPkg,
		fromPkgVersion: fromVer,
		toPkgName:      toPkg,
		toPkgVersion:   toVer,
	}
}

func (r ModReplace) String() string {
	return fmt.Sprintf("%s %s => %s %s", r.fromPkgName, r.fromPkgVersion, r.toPkgName, r.toPkgVersion)
}

type ModFile struct {
	moduleName string
	version    string

	required []ModRequired
	replaces []ModReplace
}

func NewModFile(name, version string) *ModFile {
	return &ModFile{
		moduleName: name,
		version:    version,
		required:   make([]ModRequired, 0),
	}
}

func (m *ModFile) AddRequired(pkg, ver string) *ModFile {
	m.required = append(m.required, NewModRequired(pkg, ver))
	return m
}

func (m *ModFile) AddReplace(fromPkg, fromVer, toPkg, toVer string) *ModFile {
	m.replaces = append(m.replaces, NewModReplace(fromPkg, fromVer, toPkg, toVer))
	return m
}

func (m *ModFile) requiredString() string {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteString("require (\n")
	for _, r := range m.required {
		buf.WriteString("\t" + r.String() + "\n")
	}
	buf.WriteString(")\n")

	return buf.String()
}

func (m *ModFile) replacesString() string {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteString("replace (\n")
	for _, r := range m.replaces {
		buf.WriteString("\t" + r.String() + "\n")
	}
	buf.WriteString(")\n")

	return buf.String()
}

func (m *ModFile) String() string {
	return fmt.Sprintf(`module %s

go %s

%s

%s
`, m.moduleName, m.version, m.replacesString(), m.requiredString())
}
