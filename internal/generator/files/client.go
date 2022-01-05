package files

import (
	"bytes"
	"fmt"
	"gitee.com/eden-framework/eden-framework/internal/generator/importer"
	"gitee.com/eden-framework/eden-framework/internal/generator/operator"
	"gitee.com/eden-framework/eden-framework/internal/generator/scanner"
	str "gitee.com/eden-framework/strings"
	"github.com/go-courier/oas"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
)

type ClientFile struct {
	ClientName  string
	PackageName string
	Name        string
	Importer    *importer.PackageImporter
	a           *oas.OpenAPI
	ops         map[string]operator.Op
}

func NewClientFile(name string, a *oas.OpenAPI) *ClientFile {
	f := &ClientFile{
		Name:        str.ToLowerLinkCase(name),
		PackageName: str.ToLowerSnakeCase("client-" + name),
		ClientName:  str.ToUpperCamelCase("client-" + name),
		Importer:    importer.NewPackageImporter(""),
		a:           a,
		ops:         make(map[string]operator.Op),
	}

	for path, pathItem := range a.Paths.Paths {
		for method, op := range pathItem.Operations.Operations {
			f.AddOp(operator.NewOperation(name, strings.ToUpper(string(method)), path, op, a.Components, op.Extensions))
		}
	}

	return f
}

func (c *ClientFile) AddOp(op operator.Op) {
	if op != nil {
		c.ops[op.ID()] = op
	}
}

func (c *ClientFile) WritePackage(w io.Writer) (err error) {
	_, err = io.WriteString(w, fmt.Sprintf("package %s\n\n", c.PackageName))
	return
}

func (c *ClientFile) WriteImports(w io.Writer) (err error) {
	_, err = io.WriteString(w, c.Importer.String())
	return
}

func (c *ClientFile) WriteTypeInterface(w io.Writer) (err error) {
	keys := make([]string, 0)
	for key := range c.ops {
		if key == "Swagger" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	_, err = io.WriteString(w, `type `+c.ClientName+`Interface interface {
`)
	if err != nil {
		return
	}

	for _, key := range keys {
		op := c.ops[key]

		var interfaceMethod string
		if op.HasRequest() {
			if annotate := op.Annotation(); annotate[scanner.XAnnotationRevert] != nil && op.CanRevert() {
				interfaceMethod = annotate[scanner.XAnnotationRevert].Run(operator.CmdGenerateInterface, op) + "\n"
				c.Importer.Merge(annotate[scanner.XAnnotationRevert].Importer())
			} else {
				requestStr := fmt.Sprintf("%s %s, ", "req", operator.RequestOf(op.ID()))
				interfaceMethod = op.ID() + `(` + requestStr + `metas... ` + c.Importer.Use("gitee.com/eden-framework/courier.Metadata") + `) (resp *` + operator.ResponseOf(op.ID()) + `, err error)
`
			}
		} else {
			interfaceMethod = op.ID() + `(metas... ` + c.Importer.Use("gitee.com/eden-framework/courier.Metadata") + `) (resp *` + operator.ResponseOf(op.ID()) + `, err error)
`
		}

		_, err = io.WriteString(w, interfaceMethod)
		if err != nil {
			return
		}
	}

	_, err = io.WriteString(w, `}
`)
	return
}

func (c *ClientFile) WriteTypeInstance(w io.Writer) (err error) {
	keys := make([]string, 0)
	for key := range c.ops {
		if key == "Swagger" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	_, err = io.WriteString(w, `
type `+c.ClientName+` struct {
	*`+c.Importer.Use(scanner.PkgImportPathClient+".Client")+`
}

func (c *`+c.ClientName+`) MarshalDefaults() {
	c.Name = "`+c.Name+`"
	c.Client.MarshalDefaults(c.Client)
}

func (c  *`+c.ClientName+`) Init() {
	c.MarshalDefaults()
	c.CheckService()
`)
	if err != nil {
		return
	}

	for _, key := range keys {
		op := c.ops[key]
		if annotate := op.Annotation(); annotate[scanner.XAnnotationRevert] != nil && op.CanRevert() {
			if target := op.RevertTarget(); target != "" {
				_, err = io.WriteString(w, c.Importer.Use("gitee.com/eden-framework/revert.RegisterRevertFunc")+`("`+target+`", c.`+op.ID()+`)
`)
			}
		}
	}
	_, err = io.WriteString(w, `}

func (c  `+c.ClientName+`) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := `+c.Importer.Use("gitee.com/eden-framework/courier/status_error.FromError")+`(err)
	if statusErr.Code == int64(`+c.Importer.Use("gitee.com/eden-framework/courier/status_error.RequestTimeout")+`) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}
`)
	return
}

func (c *ClientFile) WriteOperations(w io.Writer) (err error) {
	keys := make([]string, 0)
	for key := range c.ops {
		if key == "Swagger" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		op := c.ops[key]

		var operationBody string
		var extensionBody string
		if op.HasRequest() {
			_, err = io.WriteString(w, `
type `+operator.RequestOf(op.ID())+" ")
			if err != nil {
				return
			}

			err = op.WriteReqType(w, c.Importer)
			if err != nil {
				return
			}

			if annotate := op.Annotation(); annotate[scanner.XAnnotationRevert] != nil && op.CanRevert() {
				operationBody = annotate[scanner.XAnnotationRevert].Run(operator.CmdGenerateImplement, op)
				if target := op.RevertTarget(); target != "" {
					targetOp := c.ops[target]
					if targetOp != nil {
						extensionBody = annotate[scanner.XAnnotationRevert].Run(operator.CmdGenerateGetRevertID, targetOp)
					}
				}

				c.Importer.Merge(annotate[scanner.XAnnotationRevert].Importer())
			} else {
				requestStr := fmt.Sprintf("%s %s, ", "req", operator.RequestOf(op.ID()))
				interfaceMethod := op.ID() + `(` + requestStr + `metas... ` + c.Importer.Use("gitee.com/eden-framework/courier.Metadata") + `) (resp *` + operator.ResponseOf(op.ID()) + `, err error)`
				operationBody = `
func (c ` + c.ClientName + `) ` + interfaceMethod + ` {
	resp = &` + operator.ResponseOf(op.ID()) + `{}
	resp.Meta = ` + c.Importer.Use("gitee.com/eden-framework/courier.Metadata") + `{}

	err = c.Request(c.Name + ".` + op.ID() + `", "` + op.Method() + `", "` + op.Path() + `", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}
`
			}
		} else {
			interfaceMethod := op.ID() + `(metas... ` + c.Importer.Use("gitee.com/eden-framework/courier.Metadata") + `) (resp *` + operator.ResponseOf(op.ID()) + `, err error)`
			operationBody = `
func (c ` + c.ClientName + `) ` + interfaceMethod + ` {
	resp = &` + operator.ResponseOf(op.ID()) + `{}
	resp.Meta = ` + c.Importer.Use("gitee.com/eden-framework/courier.Metadata") + `{}

	err = c.Request(c.Name + ".` + op.ID() + `", "` + op.Method() + `", "` + op.Path() + `", nil, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}
`
		}

		_, err = io.WriteString(w, operationBody)
		if err != nil {
			return
		}

		_, err = io.WriteString(w, extensionBody)
		if err != nil {
			return
		}

		_, err = io.WriteString(w, `
type `+operator.ResponseOf(op.ID())+`  struct {
	Meta `+c.Importer.Use("gitee.com/eden-framework/courier.Metadata")+`
	Body `)
		if err != nil {
			return
		}

		err = op.WriteRespBodyType(w, c.Importer)
		if err != nil {
			return
		}

		_, err = io.WriteString(w, `
}
`)
		if err != nil {
			return
		}
	}

	return nil
}

func (c *ClientFile) WriteAll() string {
	w := bytes.NewBuffer([]byte{})

	err := c.WriteTypeInterface(w)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteTypeInstance(w)
	if err != nil {
		logrus.Panic(err)
	}

	err = c.WriteOperations(w)
	if err != nil {
		logrus.Panic(err)
	}

	return w.String()
}

func (c *ClientFile) String() string {
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
