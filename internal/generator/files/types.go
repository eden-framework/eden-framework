package files

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/eden-framework/internal/generator/importer"
	"github.com/eden-framework/eden-framework/internal/generator/operator"
	str "github.com/eden-framework/strings"
	"github.com/go-courier/oas"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
)

type TypesFile struct {
	a               *oas.OpenAPI
	PackageName     string
	Importer        *importer.PackageImporter
	serviceName     string
	definitionNames []string
}

func NewTypesFile(serviceName string, a *oas.OpenAPI) *TypesFile {
	names := make([]string, 0)
	for name := range a.Components.Schemas {
		names = append(names, name)
	}
	sort.Strings(names)

	f := &TypesFile{
		a:               a,
		PackageName:     str.ToLowerSnakeCase("client-" + serviceName),
		Importer:        importer.NewPackageImporter(""),
		serviceName:     serviceName,
		definitionNames: names,
	}

	return f
}

func (f *TypesFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", f.PackageName))
	return
}

func (f *TypesFile) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, f.Importer.String())
	return
}

func (f *TypesFile) WriteDefinition(w io.Writer) (err error) {
	for _, name := range f.definitionNames {
		schema := f.a.Components.Schemas[name]
		typ, alias := operator.NewTypeGenerator(f.serviceName, f.Importer).Type(schema)
		op := " "
		if alias {
			op = " = "
		}
		contentStr := `
type ` + name + op + typ + `
`
		_, err = io.WriteString(w, contentStr)
		if err != nil {
			return
		}
	}
	return
}

func (f *TypesFile) WriteAll() string {
	buf := bytes.NewBuffer([]byte{})

	err := f.WriteDefinition(buf)
	if err != nil {
		logrus.Panic(err)
	}

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
