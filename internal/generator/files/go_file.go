package files

import (
	"bytes"
	"fmt"
	"github.com/profzone/eden-framework/internal/generator/importer"
	"github.com/sirupsen/logrus"
	"html/template"
)

type GoFile struct {
	PackageName string
	*importer.PackageImporter

	buf *bytes.Buffer
}

func NewGoFile(pkgName string) *GoFile {
	return &GoFile{
		PackageName:     pkgName,
		PackageImporter: importer.NewPackageImporter(""),
		buf:             bytes.NewBuffer([]byte{}),
	}
}

func (f *GoFile) WithBlock(tpl string) *GoFile {
	f.write(tpl)
	return f
}

func (f *GoFile) write(tpl string) {
	t, err := template.New(f.PackageName).Parse(tpl)
	if err != nil {
		logrus.Panicf("template parse failed: %v", err)
	}

	err = t.Execute(f.buf, f)
	if err != nil {
		logrus.Panicf("template execute failed: %v", err)
	}
}

func (f *GoFile) String() string {
	return fmt.Sprintf(`package %s

%s

%s
`, f.PackageName, f.PackageImporter.String(), f.buf.String())
}
