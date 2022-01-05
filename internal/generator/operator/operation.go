package operator

import (
	"fmt"
	"gitee.com/eden-framework/courier/httpx"
	"gitee.com/eden-framework/courier/transport_http/transform"
	"gitee.com/eden-framework/eden-framework/internal/generator/importer"
	"gitee.com/eden-framework/eden-framework/internal/generator/scanner"
	str "gitee.com/eden-framework/strings"
	"github.com/go-courier/oas"
	"io"
)

type Op interface {
	ID() string
	Method() string
	Path() string
	HasRequest() bool
	Annotation() map[string]Annotation
	WriteReqType(w io.Writer, ipt *importer.PackageImporter) error
	WriteRespBodyType(w io.Writer, ipt *importer.PackageImporter) error

	CanRevert() bool
	RevertIDField() string
	RevertTarget() string
}

type Operation struct {
	serviceName string
	method      string
	path        string
	*oas.Operation
	components oas.Components
	annotation map[string]Annotation
}

func NewOperation(serviceName string, method string, path string, operation *oas.Operation, components oas.Components, extension map[string]interface{}) *Operation {
	op := &Operation{
		serviceName: serviceName,
		method:      method,
		path:        path,
		Operation:   operation,
		components:  components,
		annotation:  make(map[string]Annotation),
	}
	clientName := str.ToUpperCamelCase("client-" + serviceName)
	for extensionKey, annotate := range extension {
		a := getAnnotation(extensionKey)
		if a == nil {
			panic(fmt.Sprintf("[NewOperation] getAnnotation is nil which extensionKey is %s", extensionKey))
		}
		a.SetArgs(annotate.(string), clientName)
		op.annotation[extensionKey] = a
	}

	return op
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

func (o *Operation) Annotation() map[string]Annotation {
	return o.annotation
}

func (o *Operation) HasRequest() bool {
	return len(o.Operation.Parameters) > 0 || o.RequestBody != nil
}

func (o *Operation) CanRevert() bool {
	return len(o.Operation.Parameters) == 1 && o.Operation.Parameters[0].Schema.Format == "uint64"
}

func (o *Operation) RevertIDField() string {
	if o.CanRevert() {
		return str.ToUpperCamelCase(o.Operation.Parameters[0].Name)
	}
	return ""
}

func (o *Operation) RevertTarget() string {
	if o.Extensions[scanner.XAnnotationRevert] != nil {
		return o.Extensions[scanner.XAnnotationRevert].(string)
	}
	return ""
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
