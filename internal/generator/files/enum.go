package files

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"profzone/eden-framework/internal/generator/importer"
	"profzone/eden-framework/internal/generator/scanner"
	str "profzone/eden-framework/pkg/strings"
	"strconv"
)

type Enum struct {
	PackagePath string
	PackageName string
	Name        string
	Options     scanner.Enum
	Importer    *importer.PackageImporter
	HasOffset   bool
}

func NewEnum(packagePath, packageName, name string, options scanner.Enum, hasOffset bool) *Enum {
	return &Enum{
		PackagePath: packagePath,
		PackageName: packageName,
		Name:        name,
		Options:     options,
		Importer:    importer.NewPackageImporter(packagePath),
		HasOffset:   hasOffset,
	}
}

func (e *Enum) ConstPrefix() string {
	return str.ToUpperSnakeCase(e.Name)
}

func (e *Enum) ConstOffset() string {
	return e.ConstPrefix() + "_OFFSET"
}

func (e *Enum) ConstUnknown() string {
	return e.ConstPrefix() + "_UNKNOWN"
}

func (e *Enum) InvalidErrorString() string {
	return fmt.Sprintf("Invalid%s", e.Name)
}

func (e *Enum) ConstKey(key interface{}) string {
	return fmt.Sprintf("%s__%v", e.ConstPrefix(), key)
}

func (e *Enum) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", e.PackageName))
	return
}

func (e *Enum) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, e.Importer.String())
	return
}

func (e *Enum) WriteVars(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("var %s = %s(\"invalid %s\")\n\n", e.InvalidErrorString(), e.Importer.Use("errors.New"), e.Name))
	return
}

func (e *Enum) WriteInitFunc(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf(`func init() {
	`+e.Importer.Use("profzone/eden-framework/pkg/enumeration.RegisterEnums")+"(\""+e.Name+"\", map[string]string{\n"))

	for _, option := range e.Options {
		_, err = io.WriteString(w, strconv.Quote(fmt.Sprintf("%v", option.Value))+":"+strconv.Quote(option.Label)+",\n")
		if err != nil {
			return
		}
	}

	_, err = io.WriteString(w, "})\n}\n\n")
	return
}

func (e *Enum) WriteParseFromStringFunc(w io.Writer) (err error) {
	funcStr := fmt.Sprintf(`func Parse%sFromString(s string) (%s, error) {
	switch s {
	case %s:
		return %s, nil
`, e.Name, e.Name, strconv.Quote(""), e.ConstUnknown())

	for _, option := range e.Options {
		funcStr += fmt.Sprintf(`case %s:
		return %s, nil
`, fmt.Sprintf(`"%v"`, option.Value), e.ConstKey(option.Value))
	}

	funcStr += fmt.Sprintf(`}
	return %s, %s
}

`, e.ConstUnknown(), e.InvalidErrorString())
	_, err = io.WriteString(w, funcStr)
	return
}

func (e *Enum) WriteAll() string {
	w := bytes.NewBuffer([]byte{})

	err := e.WriteVars(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteInitFunc(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteParseFromStringFunc(w)
	if err != nil {
		logrus.Panic(err)
	}

	return w.String()
}

func (e *Enum) String() string {
	buf := bytes.NewBuffer([]byte{})

	content := e.WriteAll()

	err := e.WritePackage(buf)
	if err != nil {
		logrus.Panic(err)
	}

	err = e.WriteImports(buf)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = io.WriteString(buf, content)
	if err != nil {
		logrus.Panic(err)
	}

	return buf.String()
}
