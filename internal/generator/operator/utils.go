package operator

import (
	"gitee.com/eden-framework/eden-framework/internal/generator/importer"
	"github.com/go-courier/oas"
	"regexp"
	"strings"
)

func RefName(ref oas.Refer) string {
	parts := strings.Split(ref.RefString(), "/")
	return parts[len(parts)-1]
}

func BasicType(schemaType string, format string, ipt *importer.PackageImporter) string {
	switch format {
	case "binary":
		return ipt.Use("mime/multipart.FileHeader")
	case "byte", "int", "int8", "int16", "int32", "int64", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64":
		return format
	case "float":
		return "float32"
	case "double":
		return "float64"
	default:
		switch schemaType {
		case "boolean":
			return "bool"
		default:
			return "string"
		}
	}
}

func PathFromSwaggerPath(str string) string {
	r := regexp.MustCompile(`/\{([^/\\}]+)\}`)
	result := r.ReplaceAllString(str, "/:$1")
	return result
}

func RequestOf(id string) string {
	return id + "Request"
}

func ResponseOf(id string) string {
	return id + "Response"
}

func mayComposedFieldSchema(schema *oas.Schema) *oas.Schema {
	// for named field
	if schema.AllOf != nil && len(schema.AllOf) == 2 && schema.AllOf[len(schema.AllOf)-1].Type == "" {
		nextSchema := &oas.Schema{
			Reference:    schema.AllOf[0].Reference,
			SchemaObject: schema.AllOf[1].SchemaObject,
		}

		for k, v := range schema.AllOf[1].SpecExtensions.Extensions {
			nextSchema.AddExtension(k, v)
		}

		for k, v := range schema.SpecExtensions.Extensions {
			nextSchema.AddExtension(k, v)
		}

		return nextSchema
	}

	return schema
}
