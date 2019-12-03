package files

import (
	"bytes"
	"fmt"
	"github.com/profzone/eden-framework/internal/generator/api"
	"github.com/profzone/eden-framework/internal/generator/importer"
	str "github.com/profzone/eden-framework/pkg/strings"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
	"strconv"
)

type Enum struct {
	PackagePath string
	PackageName string
	Name        string
	Options     api.Enum
	Importer    *importer.PackageImporter
	HasOffset   bool
}

func NewEnum(packagePath, packageName, name string, options api.Enum, hasOffset bool) *Enum {
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
	`+e.Importer.Use("github.com/profzone/eden-framework/pkg/enumeration.RegisterEnums")+"(\""+e.Name+"\", map[string]string{\n"))

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

func (e *Enum) WriteParseFromLabelStringFunc(w io.Writer) (err error) {
	funcStr := fmt.Sprintf(`func Parse%sFromLabelString(s string) (%s, error) {
	switch s {
	case %s:
		return %s, nil
`, e.Name, e.Name, strconv.Quote(""), e.ConstUnknown())

	for _, option := range e.Options {
		funcStr += fmt.Sprintf(`case %s:
		return %s, nil
`, fmt.Sprintf(`"%v"`, option.Label), e.ConstKey(option.Value))
	}

	funcStr += fmt.Sprintf(`}
	return %s, %s
}

`, e.ConstUnknown(), e.InvalidErrorString())
	_, err = io.WriteString(w, funcStr)
	return
}

func (e *Enum) WriteEnumDescriptor(w io.Writer) (err error) {
	contentStr := fmt.Sprintf(`func (%s) EnumType() string {
	return %s
}

`, e.Name, strconv.Quote(e.Name))

	contentStr += fmt.Sprintf(`func (%s) Enums() map[int][]string {
	return map[int][]string{
`, e.Name)

	for _, option := range e.Options {
		contentStr += fmt.Sprintf("int(%s): {%s, %s},\n", e.ConstKey(option.Value), strconv.Quote(fmt.Sprintf("%s", option.Value)), strconv.Quote(option.Label))
	}

	contentStr += "}\n}\n\n"

	_, err = io.WriteString(w, contentStr)
	return
}

func (e *Enum) WriteStringer(w io.Writer) (err error) {
	contentStr := fmt.Sprintf(`func (v %s) String() string {
	switch v {
	case %s:
		return ""
`, e.Name, e.ConstUnknown())

	for _, option := range e.Options {
		contentStr += fmt.Sprintf(`case %s:
	return %s
`, e.ConstKey(option.Value), fmt.Sprintf(`"%v"`, option.Value))
	}

	contentStr += `}
	return "UNKNOWN"
}

`

	_, err = io.WriteString(w, contentStr)
	return
}

func (e *Enum) WriteLabeler(w io.Writer) (err error) {
	contentStr := fmt.Sprintf(`func (v %s) Label() string {
	switch v {
	case %s:
		return ""
`, e.Name, e.ConstUnknown())

	for _, option := range e.Options {
		contentStr += fmt.Sprintf(`case %s:
	return %s
`, e.ConstKey(option.Value), strconv.Quote(option.Label))
	}

	contentStr += `}
	return "UNKNOWN"
}

`

	_, err = io.WriteString(w, contentStr)
	return
}

func (e *Enum) WriteTextMarshalerAndUnmarshaler(w io.Writer) (err error) {
	contentStr := fmt.Sprintf(`var _ interface {
	%s
	%s
} = (*%s)(nil)

func (v %s) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, %s
	}
	return []byte(str), nil
}

func (v *%s) UnmarshalText(data []byte) (err error) {
	*v, err = Parse%sFromString(string(%s(data)))
	return
}

`,
		e.Importer.Use("encoding.TextMarshaler"),
		e.Importer.Use("encoding.TextUnmarshaler"),
		e.Name, e.Name, e.InvalidErrorString(), e.Name, e.Name,
		e.Importer.Use("bytes.ToUpper"))

	_, err = io.WriteString(w, contentStr)
	return
}

func (e *Enum) WriteScannerAndValuer(w io.Writer) (err error) {
	if !e.HasOffset {
		return
	}

	contentStr := fmt.Sprintf(`var _ interface {
	%s
	%s
} = (*%s)(nil)

func (v *%s) Scan(src interface{}) error {
	integer, err := %s(src, %s)
	if err != nil {
		return err
	}
	*v = %s(integer - %s)
	return nil
}

func (v %s) Value() (%s, error) {
	return int64(v) + %s, nil
}

`,
		e.Importer.Use("database/sql.Scanner"),
		e.Importer.Use("database/sql/driver.Valuer"),
		e.Name, e.Name,
		e.Importer.Use("github.com/profzone/eden-framework/pkg/enumeration.AsInt64"),
		e.ConstOffset(), e.Name, e.ConstOffset(), e.Name,
		e.Importer.Use("database/sql/driver.Value"),
		e.ConstOffset())

	_, err = io.WriteString(w, contentStr)
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
	err = e.WriteParseFromLabelStringFunc(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteEnumDescriptor(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteStringer(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteLabeler(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteTextMarshalerAndUnmarshaler(w)
	if err != nil {
		logrus.Panic(err)
	}
	err = e.WriteScannerAndValuer(w)
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

func (e *Enum) WriteEnumDefinition(w io.Writer) (err error) {
	contentStr := fmt.Sprintf(`// api:enum
type %s uint

const (
%s
`, e.Name, e.ConstPrefix()+"_UNKNOWN "+e.Name+" = iota")

	sort.Slice(e.Options, func(i, j int) bool {
		return e.Options[i].Val.(float64) < e.Options[j].Val.(float64)
	})

	var index = 1
	for _, enum := range e.Options {
		val := int(enum.Val.(float64))
		if val > index {
			contentStr += `(

const (
`
			contentStr += fmt.Sprintf("%s__%s %s = iota + %d // %s\n", e.ConstPrefix(), enum.Value.(string), e.Name, val, enum.Label)
			index = val + 1
		} else {
			contentStr += fmt.Sprintf("%s__%s // %s\n", e.ConstPrefix(), enum.Value.(string), enum.Label)
			index++
		}
	}
	contentStr += ")\n\n"

	_, err = io.WriteString(w, contentStr)
	return
}
