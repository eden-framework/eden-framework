package scanner

import (
	"fmt"
	"github.com/eden-framework/packagex"
	"github.com/eden-framework/reflectx"
	str "github.com/eden-framework/strings"
	"go/types"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/go-courier/oas"
	"github.com/sirupsen/logrus"
)

func NewDefinitionScanner(pkg *packagex.Package) *DefinitionScanner {
	return &DefinitionScanner{
		enumScanner:       NewEnumScanner(pkg),
		pkg:               pkg,
		ioWriterInterface: packagex.NewPackage(pkg.Pkg("io")).TypeName("Writer").Type().Underlying().(*types.Interface),
	}
}

type DefinitionScanner struct {
	enumInterfaceType *types.Interface
	pkg               *packagex.Package
	enumScanner       *EnumScanner
	definitions       map[*types.TypeName]*oas.Schema
	schemas           map[string]*oas.Schema
	ioWriterInterface *types.Interface
}

func addExtension(s *oas.Schema, key string, v interface{}) {
	if s == nil {
		return
	}
	if len(s.AllOf) > 0 {
		s.AllOf[len(s.AllOf)-1].AddExtension(key, v)
	} else {
		s.AddExtension(key, v)
	}
}

func setMetaFromDoc(s *oas.Schema, doc string) {
	if s == nil {
		return
	}

	lines := strings.Split(doc, "\n")

	for i := range lines {
		if strings.Index(lines[i], "@deprecated") != -1 {
			s.Deprecated = true
		}
	}

	description := dropMarkedLines(lines)

	if len(s.AllOf) > 0 {
		s.AllOf[len(s.AllOf)-1].Description = description
	} else {
		s.Description = description
	}
}

func (scanner *DefinitionScanner) BindSchemas(openapi *oas.OpenAPI) {
	openapi.Components.Schemas = scanner.schemas
}

func (scanner *DefinitionScanner) Def(typeName *types.TypeName) *oas.Schema {
	if s, ok := scanner.definitions[typeName]; ok {
		return s
	}

	logrus.Debugf("scanning Type `%s.%s`", typeName.Pkg().Path(), typeName.Name())

	if typeName.IsAlias() {
		typeName = typeName.Type().(*types.Named).Obj()
	}

	doc := scanner.pkg.CommentsOf(scanner.pkg.IdentOf(typeName.Type().(*types.Named).Obj()))

	// register empty before scan
	// to avoid cycle
	scanner.setDef(typeName, &oas.Schema{})

	if doc, fmtName := parseStrfmt(doc); fmtName != "" {
		s := oas.NewSchema(oas.TypeString, fmtName)
		setMetaFromDoc(s, doc)
		return scanner.setDef(typeName, s)
	}

	if doc, typ := parseType(doc); typ != "" {
		s := oas.NewSchema(oas.Type(typ), "")
		setMetaFromDoc(s, doc)
		return scanner.setDef(typeName, s)
	}

	if reflectx.FromTType(types.NewPointer(typeName.Type())).Implements(reflectx.FromTType(scanner.ioWriterInterface)) {
		return scanner.setDef(typeName, oas.Binary())
	}

	if typeName.Pkg() != nil {
		if typeName.Pkg().Path() == "time" && typeName.Name() == "Time" {
			return scanner.setDef(typeName, oas.DateTime())
		}
	}

	doc, hasEnum := ParseEnum(doc)
	if hasEnum {
		enumOptions := scanner.enumScanner.Enum(typeName)
		if enumOptions == nil {
			panic(fmt.Errorf("missing enum option but annotated by openapi:enum"))
		}
		s := oas.String()
		for _, e := range enumOptions {
			s.Enum = append(s.Enum, e.Value)
		}
		s.AddExtension(XEnumOptions, enumOptions)
		return scanner.setDef(typeName, s)
	}

	s := scanner.GetSchemaByType(typeName.Type().Underlying())

	setMetaFromDoc(s, doc)

	return scanner.setDef(typeName, s)
}

func (scanner *DefinitionScanner) isInternal(typeName *types.TypeName) bool {
	return strings.HasPrefix(typeName.Pkg().Path(), strings.TrimSuffix(scanner.pkg.PkgPath, "/cmd"))
}

func (scanner *DefinitionScanner) typeUniqueName(typeName *types.TypeName, isExist func(name string) bool) (string, bool) {
	typePkgPath := typeName.Pkg().Path()
	name := typeName.Name()

	if scanner.isInternal(typeName) {
		pathParts := strings.Split(typePkgPath, "/")
		count := 1
		for isExist(name) {
			name = str.ToUpperCamelCase(pathParts[len(pathParts)-count]) + name
			count++
		}
		return name, true
	}

	return str.ToUpperCamelCase(typePkgPath) + name, false
}

func (scanner *DefinitionScanner) reformatSchemas() {
	typeNameList := make([]*types.TypeName, 0)

	for typeName := range scanner.definitions {
		v := typeName
		typeNameList = append(typeNameList, v)
	}

	sort.Slice(typeNameList, func(i, j int) bool {
		return scanner.isInternal(typeNameList[i]) && fullTypeName(typeNameList[i]) < fullTypeName(typeNameList[j])
	})

	schemas := map[string]*oas.Schema{}

	for _, typeName := range typeNameList {
		name, _ := scanner.typeUniqueName(typeName, func(name string) bool {
			_, exists := schemas[name]
			return exists
		})

		s := scanner.definitions[typeName]
		addExtension(s, XID, name)
		addExtension(s, XGoVendorType, fullTypeName(typeName))
		schemas[name] = s
	}

	scanner.schemas = schemas
}

func (scanner *DefinitionScanner) setDef(typeName *types.TypeName, schema *oas.Schema) *oas.Schema {
	if scanner.definitions == nil {
		scanner.definitions = map[*types.TypeName]*oas.Schema{}
	}
	scanner.definitions[typeName] = schema
	scanner.reformatSchemas()
	return schema
}

func NewSchemaRefer(s *oas.Schema) *SchemaRefer {
	return &SchemaRefer{
		Schema: s,
	}
}

type SchemaRefer struct {
	*oas.Schema
}

func (r SchemaRefer) RefString() string {
	s := r.Schema
	if r.Schema.AllOf != nil {
		s = r.AllOf[len(r.Schema.AllOf)-1]
	}
	return oas.NewComponentRefer("schemas", s.Extensions[XID].(string)).RefString()
}

func (scanner *DefinitionScanner) GetSchemaByType(typ types.Type) *oas.Schema {
	switch t := typ.(type) {
	case *types.Named:
		if t.String() == "mime/multipart.FileHeader" {
			return oas.Binary()
		}
		return oas.RefSchemaByRefer(NewSchemaRefer(scanner.Def(t.Obj())))
	case *types.Interface:
		return &oas.Schema{}
	case *types.Basic:
		typeName, format := getSchemaTypeFromBasicType(reflectx.FromTType(t).Kind().String())
		if typeName != "" {
			return oas.NewSchema(typeName, format)
		}
	case *types.Pointer:
		count := 1
		elem := t.Elem()

		for {
			if p, ok := elem.(*types.Pointer); ok {
				elem = p.Elem()
				count++
			} else {
				break
			}
		}

		s := scanner.GetSchemaByType(elem)
		markPointer(s, count)
		return s
	case *types.Map:
		keySchema := scanner.GetSchemaByType(t.Key())
		if keySchema != nil && len(keySchema.Type) > 0 && keySchema.Type != "string" {
			panic(fmt.Errorf("only support map[string]interface{}"))
		}
		return oas.KeyValueOf(keySchema, scanner.GetSchemaByType(t.Elem()))
	case *types.Slice:
		return oas.ItemsOf(scanner.GetSchemaByType(t.Elem()))
	case *types.Array:
		length := uint64(t.Len())
		s := oas.ItemsOf(scanner.GetSchemaByType(t.Elem()))
		s.MaxItems = &length
		s.MinItems = &length
		return s
	case *types.Struct:
		err := (StructFieldUniqueChecker{}).Check(t, false)
		if err != nil {
			panic(fmt.Errorf("type %s: %s", typ, err))
		}

		structSchema := oas.ObjectOf(nil)
		schemas := make([]*oas.Schema, 0)

		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)

			if !field.Exported() {
				continue
			}

			structFieldType := field.Type()

			tags := reflect.StructTag(t.Tag(i))

			tagValueForName := tags.Get("json")
			if tagValueForName == "" {
				tagValueForName = tags.Get("name")
			}

			name, flags := tagValueAndFlagsByTagString(tagValueForName)
			if name == "-" {
				continue
			}

			if name == "" && field.Anonymous() {
				if field.Type().String() == "bytes.Buffer" {
					structSchema = oas.Binary()
					break
				}
				s := scanner.GetSchemaByType(structFieldType)
				if s != nil {
					schemas = append(schemas, s)
				}
				continue
			}

			if name == "" {
				name = field.Name()
			}

			required := true
			if hasOmitempty, ok := flags["omitempty"]; ok {
				required = !hasOmitempty
			}

			structSchema.SetProperty(
				name,
				scanner.propSchemaByField(field.Name(), structFieldType, tags, name, flags, scanner.pkg.CommentsOf(scanner.pkg.IdentOf(field))),
				required,
			)
		}

		if len(schemas) > 0 {
			return oas.AllOf(append(schemas, structSchema)...)
		}

		return structSchema
	}
	return nil
}

func (scanner *DefinitionScanner) propSchemaByField(
	fieldName string,
	fieldType types.Type,
	tags reflect.StructTag,
	name string,
	flags map[string]bool,
	desc string,
) *oas.Schema {
	propSchema := scanner.GetSchemaByType(fieldType)

	refSchema := (*oas.Schema)(nil)

	if propSchema.Refer != nil {
		refSchema = propSchema
		propSchema = &oas.Schema{}
		propSchema.Extensions = refSchema.Extensions
	}

	defaultValue := tags.Get("default")
	//validate, hasValidate := tags.Lookup("validate")

	if flags != nil && flags["string"] {
		propSchema.Type = oas.TypeString
	}

	if defaultValue != "" {
		propSchema.Default = defaultValue
	}

	//if hasValidate {
	//	if err := BindSchemaValidationByValidateBytes(propSchema, fieldType, []byte(validate)); err != nil {
	//		panic(err)
	//	}
	//}

	setMetaFromDoc(propSchema, desc)
	propSchema.AddExtension(XGoFieldName, fieldName)

	tagKeys := map[string]string{
		"name":     XTagName,
		"mime":     XTagMime,
		"json":     XTagJSON,
		"xml":      XTagXML,
		"validate": XTagValidate,
	}

	for k, extKey := range tagKeys {
		if v, ok := tags.Lookup(k); ok {
			propSchema.AddExtension(extKey, v)
		}
	}

	if refSchema != nil {
		return oas.AllOf(
			refSchema,
			propSchema,
		)
	}

	return propSchema
}

type StructFieldUniqueChecker map[string]*types.Var

func (checker StructFieldUniqueChecker) Check(structType *types.Struct, anonymous bool) error {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if !field.Exported() {
			continue
		}
		if field.Anonymous() {
			if named, ok := field.Type().(*types.Named); ok {
				if st, ok := named.Underlying().(*types.Struct); ok {
					if err := checker.Check(st, true); err != nil {
						return err
					}
				}
			}
			continue
		}
		if anonymous {
			if _, ok := checker[field.Name()]; ok {
				return fmt.Errorf("%s.%s already defined in other anonymous field", structType.String(), field.Name())
			}
			checker[field.Name()] = field
		}
	}
	return nil
}

type VendorExtensible interface {
	AddExtension(key string, value interface{})
}

func markPointer(vendorExtensible VendorExtensible, count int) {
	vendorExtensible.AddExtension(XGoStarLevel, count)
}

var (
	reStrFmt = regexp.MustCompile(`open-?api:strfmt\s+(\S+)([\s\S]+)?$`)
	reType   = regexp.MustCompile(`open-?api:type\s+(\S+)([\s\S]+)?$`)
)

func parseStrfmt(doc string) (string, string) {
	matched := reStrFmt.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

func parseType(doc string) (string, string) {
	matched := reType.FindAllStringSubmatch(doc, -1)
	if len(matched) > 0 {
		return strings.TrimSpace(matched[0][2]), matched[0][1]
	}
	return doc, ""
}

var basicTypeToSchemaType = map[string][2]string{
	"invalid": {"null", ""},

	"bool":    {"boolean", ""},
	"error":   {"string", "string"},
	"float32": {"number", "float"},
	"float64": {"number", "double"},

	"int":   {"integer", "int32"},
	"int8":  {"integer", "int8"},
	"int16": {"integer", "int16"},
	"int32": {"integer", "int32"},
	"int64": {"integer", "int64"},

	"rune": {"integer", "int32"},

	"uint":   {"integer", "uint32"},
	"uint8":  {"integer", "uint8"},
	"uint16": {"integer", "uint16"},
	"uint32": {"integer", "uint32"},
	"uint64": {"integer", "uint64"},

	"byte": {"integer", "uint8"},

	"string": {"string", ""},
}

func getSchemaTypeFromBasicType(basicTypeName string) (typ oas.Type, format string) {
	if schemaTypeAndFormat, ok := basicTypeToSchemaType[basicTypeName]; ok {
		return oas.Type(schemaTypeAndFormat[0]), schemaTypeAndFormat[1]
	}
	panic(fmt.Errorf("unsupported type %q", basicTypeName))
}
