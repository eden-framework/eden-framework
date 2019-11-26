package files

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"profzone/eden-framework/internal/generator/api"
	"profzone/eden-framework/internal/generator/importer"
	str "profzone/eden-framework/pkg/strings"
	"strings"
)

type ClientFile struct {
	ClientName  string
	PackageName string
	Name        string
	Importer    *importer.PackageImporter
	a           *api.Api
}

func NewClientFile(name string, a *api.Api) *ClientFile {
	return &ClientFile{
		Name:        str.ToLowerLinkCase(name),
		PackageName: str.ToLowerSnakeCase("client-" + name),
		ClientName:  str.ToUpperCamelCase("client-" + name),
		Importer:    importer.NewPackageImporter(""),
		a:           a,
	}
}

func (c *ClientFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", c.PackageName))
	return
}

func (c *ClientFile) WriteTypeInterface(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("type %sInterface interface {\n", c.ClientName))
	if err != nil {
		return err
	}
	for groupName, group := range c.a.Operators {
		for methodName, method := range group.Methods {
			req, resp := make([]string, 0), make([]string, 0)
			var typeName string
			for _, modelName := range method.Inputs {
				model, ok := c.a.Models[modelName]
				if !ok {
					logrus.Panic(fmt.Errorf("%s not exist in model definations", modelName))
				}
				if model.NeedAlias {
					typeName = c.Importer.Use(modelName)
				} else {
					typeName = model.Name
				}
				req = append(req, fmt.Sprintf("%s *%s", str.ToLowerCamelCase(model.Name), typeName))
			}
			for _, modelName := range method.Outputs {
				model, ok := c.a.Models[modelName]
				if !ok {
					logrus.Panic(fmt.Errorf("%s not exist in model definations", modelName))
				}
				if model.NeedAlias {
					typeName = c.Importer.Use(modelName)
				} else {
					typeName = model.Name
				}
				resp = append(resp, fmt.Sprintf("%s *%s", str.ToLowerCamelCase(model.Name), typeName))
			}
			resp = append(resp, "err error")
			methodString := fmt.Sprintf("%s(%s) (%s)\n", str.ToUpperCamelCase(groupName+methodName), strings.Join(req, ", "), strings.Join(resp, ", "))
			_, err = io.WriteString(w, methodString)
			if err != nil {
				return err
			}
		}
	}

	_, err = io.WriteString(w, "}\n")
	return err
}

func (c *ClientFile) String() string {
	buf := bytes.NewBuffer([]byte{})

	err := c.WritePackage(buf)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteTypeInterface(buf)
	if err != nil {
		logrus.Panic(err)
	}

	return buf.String()
}
