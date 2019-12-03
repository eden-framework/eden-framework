package files

import (
	"bytes"
	"fmt"
	"github.com/profzone/eden-framework/internal/generator/api"
	"github.com/profzone/eden-framework/internal/generator/importer"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"io"
)

type TypesFile struct {
	a           *api.Api
	PackageName string
	Importer    *importer.PackageImporter
}

func NewTypesFile(serviceName string, a *api.Api) *TypesFile {
	return &TypesFile{
		a:           a,
		PackageName: str.ToLowerSnakeCase("client-" + serviceName),
		Importer:    importer.NewPackageImporter(""),
	}
}

func (f *TypesFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", f.PackageName))
	return
}

func (f *TypesFile) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, f.Importer.String())
	return
}

func (f *TypesFile) WriteAll() string {
	buf := bytes.NewBuffer([]byte{})

	return buf.String()
}

func (f *TypesFile) String() string {
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
