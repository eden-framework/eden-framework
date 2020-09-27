package files

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/eden-framework/internal/generator/api"
	"github.com/eden-framework/eden-framework/internal/generator/importer"
	str "github.com/eden-framework/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"io"
)

type ClientEnumsFile struct {
	PackageName string
	Importer    *importer.PackageImporter
	a           *api.Api
}

func NewClientEnumsFile(outputPath, serviceName string, a *api.Api) *ClientEnumsFile {
	pkgName := str.ToLowerSnakeCase("client-" + serviceName)
	return &ClientEnumsFile{
		PackageName: pkgName,
		Importer:    importer.NewPackageImporter(""),
		a:           a,
	}
}

func (f *ClientEnumsFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", f.PackageName))
	return
}

func (f *ClientEnumsFile) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, f.Importer.String())
	return
}

func (f *ClientEnumsFile) WriteAll() string {
	buf := bytes.NewBuffer([]byte{})

	for enumTypeFullName, enum := range f.a.Enums {
		_, decl := importer.GetPackagePathAndDecl(enumTypeFullName)
		enum := NewEnum("", f.PackageName, decl, enum, false)
		err := enum.WriteEnumDefinition(buf)
		if err != nil {
			logrus.Panic(err)
		}

		content := enum.WriteAll()
		_, err = io.WriteString(buf, content)
		if err != nil {
			logrus.Panic(err)
		}

		f.Importer.Merge(enum.Importer)
	}

	return buf.String()
}

func (f *ClientEnumsFile) String() string {
	buf := bytes.NewBuffer([]byte{})

	content := f.WriteAll()

	err := f.WritePackage(buf)
	if err != nil {
		logrus.Panic(err)
	}

	err = f.WriteImports(buf)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = io.WriteString(buf, content)
	if err != nil {
		logrus.Panic(err)
	}

	return buf.String()
}
