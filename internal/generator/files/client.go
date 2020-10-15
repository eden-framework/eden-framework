package files

import (
	"bytes"
	"fmt"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/courier/transport_http/transform"
	"github.com/eden-framework/eden-framework/internal/generator/importer"
	"github.com/eden-framework/eden-framework/internal/generator/scanner"
	str "github.com/eden-framework/strings"
	"github.com/go-courier/oas"
	"github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
)

type Op interface {
	ID() string
	Method() string
	Path() string
	HasRequest() bool
	WriteReqType(w io.Writer, ipt *importer.PackageImporter) error
	WriteRespBodyType(w io.Writer, ipt *importer.PackageImporter) error
}

type Operation struct {
	serviceName string
	method      string
	path        string
	*oas.Operation
	components oas.Components
}

func NewOperation(serviceName string, method string, path string, operation *oas.Operation, components oas.Components) *Operation {
	return &Operation{
		serviceName: serviceName,
		method:      method,
		path:        path,
		Operation:   operation,
		components:  components,
	}
}

func (o *Operation) ID() string {
	return o.Operation.OperationId
}

func (o *Operation) Method() string {
	return o.method
}

func (o *Operation) Path() string {
	return PathFromSwaggerPath(o.path)
}

func (o *Operation) HasRequest() bool {
	return len(o.Operation.Parameters) > 0 || o.RequestBody != nil
}

func (o *Operation) WriteReqType(w io.Writer, ipt *importer.PackageImporter) error {
	_, err := io.WriteString(w, `struct {
`)

	for _, parameter := range o.Parameters {
		schema := mayComposedFieldSchema(parameter.Schema)

		fieldName := str.ToUpperCamelCase(parameter.Name)
		if parameter.Extensions[scanner.XGoFieldName] != nil {
			fieldName = parameter.Extensions[scanner.XGoFieldName].(string)
		}

		field := NewField(fieldName)
		field.AddTag("in", string(parameter.In))
		field.AddTag("name", parameter.Name)

		field.Comment = parameter.Description

		if parameter.Extensions[scanner.XTagValidate] != nil {
			field.AddTag("validate", fmt.Sprintf("%s", parameter.Extensions[scanner.XTagValidate]))
		}

		if !parameter.Required {
			if schema != nil {
				d := fmt.Sprintf("%v", schema.Default)
				if schema.Default != nil && d != "" {
					field.AddTag("default", d)
				}
			}
			field.AddTag("name", parameter.Name, "omitempty")
		}

		if schema != nil {
			field.Type, _ = NewTypeGenerator(o.serviceName, ipt).Type(schema)
		}

		_, err = io.WriteString(w, field.String())
		if err != nil {
			return err
		}
	}

	if o.RequestBody != nil {
		field := NewField("Body")
		if jsonMedia, ok := o.RequestBody.Content[httpx.MIME_JSON]; ok && jsonMedia.Schema != nil {
			field.Type, _ = NewTypeGenerator(o.serviceName, ipt).Type(jsonMedia.Schema)
			field.Comment = jsonMedia.Schema.Description
			field.AddTag("in", "body")
			field.AddTag("fmt", transform.GetContentTransformer(httpx.MIME_JSON).Key)
		}
		if formMedia, ok := o.RequestBody.Content[httpx.MIME_MULTIPART_FORM_DATA]; ok && formMedia.Schema != nil {
			field.Type, _ = NewTypeGenerator(o.serviceName, ipt).Type(formMedia.Schema)
			field.Comment = formMedia.Schema.Description
			field.AddTag("in", "formData,multipart")
		}
		if formMedia, ok := o.RequestBody.Content[httpx.MIME_POST_URLENCODED]; ok && formMedia.Schema != nil {
			field.Type, _ = NewTypeGenerator(o.serviceName, ipt).Type(formMedia.Schema)
			field.Comment = formMedia.Schema.Description
			field.AddTag("in", "formData")
		}
		_, err = io.WriteString(w, field.String())
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(w, `
}
`)
	return err
}

func (o *Operation) WriteRespBodyType(w io.Writer, ipt *importer.PackageImporter) error {
	respBodySchema := o.respBodySchema()
	if respBodySchema == nil {
		_, err := io.WriteString(w, `[]byte`)
		return err
	}
	tpe, _ := NewTypeGenerator(o.serviceName, ipt).Type(respBodySchema)
	_, err := io.WriteString(w, tpe)
	return err
}

func (o *Operation) respBodySchema() (schema *oas.Schema) {
	if o.Responses.Responses == nil {
		return nil
	}

	for code, resp := range o.Responses.Responses {
		if resp.Refer != nil && o.components.Responses != nil {
			if presetResponse, ok := o.components.Responses[RefName(resp.Refer)]; ok {
				resp = presetResponse
			}
		}

		if code >= 200 && code < 300 {
			if resp.Content[httpx.MIME_JSON] != nil {
				schema = resp.Content[httpx.MIME_JSON].Schema
				return
			}
		}
	}

	return
}

type ClientFile struct {
	ClientName  string
	PackageName string
	Name        string
	Importer    *importer.PackageImporter
	a           *oas.OpenAPI
	ops         map[string]Op
}

func NewClientFile(name string, a *oas.OpenAPI) *ClientFile {
	f := &ClientFile{
		Name:        str.ToLowerLinkCase(name),
		PackageName: str.ToLowerSnakeCase("client-" + name),
		ClientName:  str.ToUpperCamelCase("client-" + name),
		Importer:    importer.NewPackageImporter(""),
		a:           a,
		ops:         make(map[string]Op),
	}

	for path, pathItem := range a.Paths.Paths {
		for method, op := range pathItem.Operations.Operations {
			f.AddOp(NewOperation(name, strings.ToUpper(string(method)), path, op, a.Components))
		}
	}

	return f
}

func (c *ClientFile) AddOp(op Op) {
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

		reqVar := ""
		reqType := ""
		reqTypeInParams := ""

		if op.HasRequest() {
			reqVar = "req"
			reqType = RequestOf(op.ID())
			reqTypeInParams = reqType + ", "
		}

		interfaceMethod := op.ID() + `(` + reqVar + ` ` + reqTypeInParams + `metas... ` + c.Importer.Use("github.com/eden-framework/courier.Metadata") + `) (resp *` + ResponseOf(op.ID()) + `, err error)
`

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
	_, err = io.WriteString(w, `
type `+c.ClientName+` struct {
	`+c.Importer.Use(scanner.PkgImportPathClient+".Client")+`
}

func (`+c.ClientName+`) MarshalDefaults(v interface{}) {
	if cl, ok := v.(* `+c.ClientName+`); ok {
		cl.Name = "`+c.Name+`"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c  `+c.ClientName+`) Init() {
	c.CheckService()
}

func (c  `+c.ClientName+`) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := `+c.Importer.Use("github.com/eden-framework/courier/status_error.FromError")+`(err)
	if statusErr.Code == int64(`+c.Importer.Use("github.com/eden-framework/courier/status_error.RequestTimeout")+`) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}
`)
	return
}

func (c *ClientFile) WriteOperations(w io.Writer) (err error) {
	keys := make([]string, 0)
	for key := range c.ops {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		op := c.ops[key]

		reqVar := ""
		reqType := ""
		reqTypeInParams := ""
		reqVarInUse := "nil"

		if op.HasRequest() {
			reqVar = "req"
			reqVarInUse = reqVar
			reqType = RequestOf(op.ID())
			reqTypeInParams = reqType + ", "

			_, err = io.WriteString(w, `
type `+reqType+" ")
			if err != nil {
				return
			}

			err = op.WriteReqType(w, c.Importer)
			if err != nil {
				return
			}
		}

		interfaceMethod := op.ID() + `(` + reqVar + ` ` + reqTypeInParams + `metas... ` + c.Importer.Use("github.com/eden-framework/courier.Metadata") + `) (resp *` + ResponseOf(op.ID()) + `, err error)`

		_, err = io.WriteString(w, `
func (c `+c.ClientName+`) `+interfaceMethod+` {
	resp = &`+ResponseOf(op.ID())+`{}
	resp.Meta = `+c.Importer.Use("github.com/eden-framework/courier.Metadata")+`{}

	err = c.Request(c.Name + ".`+op.ID()+`", "`+op.Method()+`", "`+op.Path()+`", `+reqVarInUse+`, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}
`)
		if err != nil {
			return
		}

		_, err = io.WriteString(w, `
type `+ResponseOf(op.ID())+`  struct {
	Meta `+c.Importer.Use("github.com/eden-framework/courier.Metadata")+`
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
