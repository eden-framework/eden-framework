package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eden-framework/eden-framework/internal/generator/importer"
	"github.com/eden-framework/eden-framework/internal/generator/scanner"
	"github.com/eden-framework/enumeration"
	"github.com/eden-framework/packagex"
	str "github.com/eden-framework/strings"
	"github.com/go-courier/oas"
	"io"
	"sort"
	"strings"
)

func NewTypeGenerator(serviceName string, ipt *importer.PackageImporter) *TypeGenerator {
	return &TypeGenerator{
		ServiceName: serviceName,
		Importer:    ipt,
	}
}

type TypeGenerator struct {
	ServiceName string
	Importer    *importer.PackageImporter
}

func (g *TypeGenerator) PrefixType(tpe string) string {
	return str.ToUpperCamelCase(g.ServiceName) + tpe
}

func (g *TypeGenerator) Type(schema *oas.Schema) (string, bool) {
	pointer := ""
	if schema.Extensions[scanner.XGoStarLevel] != nil {
		pointer = strings.Repeat("*", int(schema.Extensions[scanner.XGoStarLevel].(float64)))
	}
	tpe, alias := g.TypeIndirect(schema)
	return pointer + tpe, alias
}

func (g *TypeGenerator) TypeIndirect(schema *oas.Schema) (string, bool) {
	if schema == nil {
		return "interface{}", false
	}

	if schema.Refer != nil {
		return RefName(schema.Refer), true
	}

	if schema.Extensions[scanner.XGoVendorType] != nil {
		if schema.Type == "array" && schema.Items != nil && schema.Items.Format == "uint64" {
			panic("modify it")
			return g.Importer.Use("github.com/johnnyeven/libtools/httplib.Uint64List"), true
		}

		typeFullName := fmt.Sprint(schema.Extensions[scanner.XGoVendorType])
		isInCommonLib := strings.Contains(typeFullName, "github.com/eden-framework")

		if schema.Type == "string" || schema.Type == "boolean" || schema.Enum != nil || isInCommonLib {

			pkgImportName, typeName := packagex.GetPkgImportPathAndExpose(typeFullName)

			if schema.Extensions[scanner.XEnumOptions] != nil {
				enums := enumeration.Enum{}
				result, _ := json.Marshal(schema.Extensions[scanner.XEnumOptions])
				err := json.Unmarshal(result, &enums)
				if err != nil {
					panic(fmt.Sprintf("the enum options can not convert to enumeration.Enum: %s", string(result)))
				}

				typeName = g.PrefixType(typeName)
				RegisterEnum(g.ServiceName, typeName, enums...)
				return typeName, true
			}

			if schema.Type == "boolean" {
				typeName = "Bool"
				pkgImportName = "github.com/eden-framework/enumeration"
				isInCommonLib = true
			}

			return g.Importer.Use(fmt.Sprintf("%s.%s", pkgImportName, typeName)), true
		}
	}

	if len(schema.AllOf) > 0 {
		buf := &bytes.Buffer{}
		buf.WriteString(`struct {
`)

		for _, subSchema := range schema.AllOf {
			if subSchema.Refer != nil {
				field := NewField(RefName(subSchema.Refer))
				buf.WriteString(field.String())
			}

			if subSchema.Properties != nil {
				g.WriteFields(buf, subSchema)
			}
		}

		buf.WriteString(`}`)
		return buf.String(), false
	}

	if schema.Type == "object" {
		if schema.AdditionalProperties != nil {
			tpe, _ := g.Type(schema.AdditionalProperties.Schema)
			return fmt.Sprintf("map[string]%s", tpe), false
		}

		buf := &bytes.Buffer{}
		buf.WriteString(`struct {
`)

		g.WriteFields(buf, schema)

		buf.WriteString(`}`)
		return buf.String(), false
	}

	if schema.Type == "array" {
		if schema.Items != nil {
			tpe, _ := g.Type(schema.Items)
			return fmt.Sprintf("[]%s", tpe), false
		}
	}

	return BasicType(string(schema.Type), schema.Format, g.Importer), false
}

func (g *TypeGenerator) WriteFields(w io.Writer, schema *oas.Schema) error {
	if schema.Properties == nil {
		return nil
	}

	fieldNames := make([]string, 0)
	for fieldName := range schema.Properties {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	for _, fieldName := range fieldNames {
		propSchema := mayComposedFieldSchema(schema.Properties[fieldName])

		_, err := io.WriteString(w, g.FieldFrom(fieldName, propSchema, schema.Required...).String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *TypeGenerator) FieldFrom(name string, propSchema *oas.Schema, requiredFields ...string) *Field {
	isRequired := str.StringIncludes(requiredFields, name)

	fieldName := name
	if propSchema.Extensions[scanner.XGoFieldName] != nil {
		fieldName = propSchema.Extensions[scanner.XGoFieldName].(string)
	}

	field := NewField(fieldName)
	field.Comment = propSchema.Description

	if propSchema.Extensions[scanner.XTagJSON] != nil {
		tagName := fmt.Sprintf("%s", propSchema.Extensions[scanner.XTagJSON])
		flags := make([]string, 0)
		if !isRequired && !strings.Contains(tagName, "omitempty") {
			flags = append(flags, "omitempty")
		}
		field.AddTag("json", tagName, flags...)
	}

	if propSchema.Extensions[scanner.XTagName] != nil {
		tagName := fmt.Sprintf("%s", propSchema.Extensions[scanner.XTagName])
		flags := make([]string, 0)
		if !isRequired && !strings.Contains(tagName, "omitempty") {
			flags = append(flags, "omitempty")
		}
		field.AddTag("name", tagName, flags...)
	}

	if propSchema.Extensions[scanner.XTagXML] != nil {
		field.AddTag("xml", fmt.Sprintf("%s", propSchema.Extensions[scanner.XTagXML]))
	}

	if propSchema.Extensions[scanner.XTagValidate] != nil {
		field.AddTag("validate", fmt.Sprintf("%s", propSchema.Extensions[scanner.XTagValidate]))
	}

	if propSchema.Default != nil {
		field.AddTag("default", fmt.Sprintf("%v", propSchema.Default))
	}

	field.Type, _ = g.Type(propSchema)

	return field
}
