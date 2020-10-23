package files

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/courier/status_error"
	"github.com/eden-framework/eden-framework/internal/generator/importer"
	"github.com/sirupsen/logrus"
	"io"
)

type ErrCodeFile struct {
	PackageName string
	Importer    *importer.PackageImporter

	errCodes status_error.StatusErrorCodeMap
}

func NewErrCodeFile(pkgName string, errCodes status_error.StatusErrorCodeMap) *ErrCodeFile {
	f := &ErrCodeFile{
		PackageName: pkgName,
		Importer:    importer.NewPackageImporter(""),
		errCodes:    errCodes,
	}

	return f
}

func (c *ErrCodeFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", c.PackageName))
	return
}

func (c *ErrCodeFile) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, c.Importer.String())
	return
}

func (c *ErrCodeFile) WriteInit(w io.Writer) (err error) {
	_, err = io.WriteString(w, `
func init() {
`)
	if err != nil {
		return err
	}

	for _, code := range c.errCodes {
		_, err = io.WriteString(w, fmt.Sprintf("%s.StatusErrorCodes.Register(\"%s\", %d, \"%s\", \"\", %v)\n",
			c.Importer.UseWithoutAlias("github.com/eden-framework/courier/status_error", ""),
			code.Key,
			code.Code,
			code.Msg,
			code.CanBeErrorTalk))
		if err != nil {
			return
		}
	}
	_, err = io.WriteString(w, `
}
`)
	return err
}

func (c *ErrCodeFile) WriteAll() string {
	w := bytes.NewBuffer([]byte{})

	err := c.WriteInit(w)
	if err != nil {
		logrus.Panic(err)
	}

	return w.String()
}

func (c *ErrCodeFile) String() string {
	buf := bytes.NewBuffer([]byte{})

	content := c.WriteAll()

	err := c.WritePackage(buf)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteImports(buf)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = io.WriteString(buf, content)
	if err != nil {
		logrus.Panic(err)
	}

	return buf.String()
}
