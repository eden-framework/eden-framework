package scanner

import (
	"fmt"
	"github.com/eden-framework/courier/httpx"
	"github.com/eden-framework/courier/status_error"
	"github.com/eden-framework/courier/transport_http"
	"github.com/eden-framework/packagex"
	"github.com/eden-framework/reflectx"
	"github.com/go-courier/oas"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/types"
	"net/http"
	"reflect"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
)

type OperatorScanner struct {
	*DefinitionScanner
	*StatusErrScanner
	pkg       *packagex.Package
	operators map[*types.TypeName]*Operator
}

func NewOperatorScanner(pkg *packagex.Package) *OperatorScanner {
	return &OperatorScanner{
		pkg:               pkg,
		DefinitionScanner: NewDefinitionScanner(pkg),
		StatusErrScanner:  NewStatusErrScanner(pkg),
	}
}

func (scanner *OperatorScanner) Operator(typeName *types.TypeName) *Operator {
	if typeName == nil {
		return nil
	}

	if operator, ok := scanner.operators[typeName]; ok {
		return operator
	}

	logrus.Debugf("scanning Operator `%s.%s`", typeName.Pkg().Path(), typeName.Name())

	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Errorf("scan Operator `%s` failed, panic: %s; calltrace: %s", fullTypeName(typeName), fmt.Sprint(e), string(debug.Stack())))
		}
	}()

	if typeStruct, ok := typeName.Type().Underlying().(*types.Struct); ok {
		operator := &Operator{}

		operator.Tag = scanner.tagFrom(typeName.Pkg().Path())

		scanner.scanRouteMeta(operator, typeName)
		scanner.scanParameterOrRequestBody(operator, typeStruct)
		scanner.scanReturns(operator, typeName)

		// cached scanned
		if scanner.operators == nil {
			scanner.operators = map[*types.TypeName]*Operator{}
		}

		scanner.operators[typeName] = operator

		return operator
	}

	return nil
}

func (scanner *OperatorScanner) singleReturnOf(typeName *types.TypeName, name string) (string, bool) {
	if typeName == nil {
		return "", false
	}

	for _, typ := range []types.Type{
		typeName.Type(),
		types.NewPointer(typeName.Type()),
	} {
		method, ok := reflectx.FromTType(typ).MethodByName(name)
		if ok {
			results, n := scanner.pkg.FuncResultsOf(method.(*reflectx.TMethod).Func)
			if n == 1 {
				for _, v := range results[0] {
					if v.Value != nil {
						s, err := strconv.Unquote(v.Value.ExactString())
						if err != nil {
							panic(fmt.Errorf("%s: %s", err, v.Value))
						}
						return s, true
					}
				}
			}
		}
	}

	return "", false
}

func (scanner *OperatorScanner) tagFrom(pkgPath string) string {
	tag := strings.TrimPrefix(pkgPath, scanner.pkg.PkgPath)
	return strings.TrimPrefix(tag, "/")
}

func (scanner *OperatorScanner) scanRouteMeta(op *Operator, typeName *types.TypeName) {
	typeStruct := typeName.Type().Underlying().(*types.Struct)

	op.ID = typeName.Name()

	for i := 0; i < typeStruct.NumFields(); i++ {
		f := typeStruct.Field(i)
		tags := reflect.StructTag(typeStruct.Tag(i))

		if f.Anonymous() && strings.Contains(f.Type().String(), pkgImportPathHttpx+".Method") {
			if path, ok := tags.Lookup("path"); ok {
				vs := strings.Split(path, ",")
				op.Path = vs[0]

				if len(vs) > 0 {
					for i := range vs {
						switch vs[i] {
						case "deprecated":
							op.Deprecated = true
							break
						}
					}
				}
			}

			if basePath, ok := tags.Lookup("basePath"); ok {
				op.BasePath = basePath
			}

			if summary, ok := tags.Lookup("summary"); ok {
				op.Summary = summary
			}

			break
		}
	}

	lines := scanner.pkg.CommentsOf(scanner.pkg.IdentOf(typeName))
	comments := strings.Split(lines, "\n")

	for i := range comments {
		if strings.Index(comments[i], "@deprecated") != -1 {
			op.Deprecated = true
		}
	}

	if op.Summary == "" {
		comments = filterMarkedLines(comments)

		if comments[0] != "" {
			op.Summary = comments[0]
			if len(comments) > 1 {
				op.Description = strings.Join(comments[1:], "\n")
			}
		}
	}

	if method, ok := scanner.singleReturnOf(typeName, "Method"); ok {
		op.Method = method
	}

	if path, ok := scanner.singleReturnOf(typeName, "Path"); ok {
		op.Path = path
	}

	if bathPath, ok := scanner.singleReturnOf(typeName, "BasePath"); ok {
		op.BasePath = bathPath
	}
}

func (scanner *OperatorScanner) scanReturns(op *Operator, typeName *types.TypeName) {
	for _, typ := range []types.Type{
		typeName.Type(),
		types.NewPointer(typeName.Type()),
	} {
		method, ok := reflectx.FromTType(typ).MethodByName("Output")
		if ok {
			results, n := scanner.pkg.FuncResultsOf(method.(*reflectx.TMethod).Func)
			if n == 2 {
				for _, v := range results[0] {
					if v.Type != nil {
						if v.Type.String() != types.Typ[types.UntypedNil].String() {
							if op.SuccessType != nil && op.SuccessType.String() != v.Type.String() {
								logrus.Warnf(fmt.Sprintf("%s success result must be same struct, but got %v, already set %v", op.ID, v.Type, op.SuccessType))
							}
							op.SuccessType = v.Type
							op.SuccessStatus, op.SuccessResponse = scanner.getResponse(v.Type, v.Expr)
						}
					}
				}
			}

			if scanner.StatusErrScanner.StatusErrType != nil {
				op.StatusErrors = scanner.StatusErrScanner.StatusErrorsInFunc(method.(*reflectx.TMethod).Func)
				op.StatusErrorSchema = scanner.DefinitionScanner.GetSchemaByType(scanner.StatusErrScanner.StatusErrType)
			}
		}
	}
}

func (scanner *OperatorScanner) firstValueOfFunc(named *types.Named, name string) (interface{}, bool) {
	method, ok := reflectx.FromTType(types.NewPointer(named)).MethodByName(name)
	if ok {
		results, n := scanner.pkg.FuncResultsOf(method.(*reflectx.TMethod).Func)
		if n == 1 {
			for _, r := range results[0] {
				if r.IsValue() {
					if v := valueOf(r.Value); v != nil {
						return v, true
					}
				}
			}
			return nil, true
		}
	}
	return nil, false
}

func (scanner *OperatorScanner) getResponse(tpe types.Type, expr ast.Expr) (statusCode int, response *oas.Response) {
	response = &oas.Response{}

	if tpe.String() == "error" {
		statusCode = http.StatusNoContent
		return
	}

	var contentType string

	if pointer, ok := tpe.(*types.Pointer); ok {
		tpe = pointer.Elem()
	}

	if named, ok := tpe.(*types.Named); ok {
		if v, ok := scanner.firstValueOfFunc(named, "ContentType"); ok {
			if s, ok := v.(string); ok {
				contentType = s
			}
			if contentType == "" {
				contentType = "*"
			}
		}
		if v, ok := scanner.firstValueOfFunc(named, "StatusCode"); ok {
			if i, ok := v.(int64); ok {
				statusCode = int(i)
			}
		}
	}

	if contentType == "" {
		contentType = httpx.MIME_JSON
	}

	response.AddContent(contentType, oas.NewMediaTypeWithSchema(scanner.DefinitionScanner.GetSchemaByType(tpe)))

	return
}

func (scanner *OperatorScanner) scanParameterOrRequestBody(op *Operator, typeStruct *types.Struct) {
	reflectx.EachField(reflectx.FromTType(typeStruct), "name", func(field reflectx.StructField, fieldDisplayName string, omitempty bool) bool {
		location, _ := tagValueAndFlagsByTagString(field.Tag().Get("in"))

		if location == "" {
			panic(fmt.Errorf("missing tag `in` for %s of %s", field.Name(), op.ID))
		}

		name, flags := tagValueAndFlagsByTagString(field.Tag().Get("name"))

		schema := scanner.DefinitionScanner.propSchemaByField(
			field.Name(),
			field.Type().(*reflectx.TType).Type,
			field.Tag(),
			name,
			flags,
			scanner.pkg.CommentsOf(scanner.pkg.IdentOf(field.(*reflectx.TStructField).Var)),
		)

		//transformer, err := transform.TransformerMgrDefault.NewTransformer(nil, field.Type(), transform.TransformerOption{
		//	MIME: field.Tag().Get("mime"),
		//})

		//if err != nil {
		//	panic(err)
		//}

		switch location {
		case "body":
			reqBody := oas.NewRequestBody("", true)
			// TODO
			reqBody.AddContent("application/json", oas.NewMediaTypeWithSchema(schema))
			op.SetRequestBody(reqBody)
		case "query":
			op.AddNonBodyParameter(oas.QueryParameter(fieldDisplayName, schema, !omitempty))
		case "cookie":
			op.AddNonBodyParameter(oas.CookieParameter(fieldDisplayName, schema, !omitempty))
		case "header":
			op.AddNonBodyParameter(oas.HeaderParameter(fieldDisplayName, schema, !omitempty))
		case "path":
			op.AddNonBodyParameter(oas.PathParameter(fieldDisplayName, schema))
		}

		return true
	}, "in")
}

type Operator struct {
	transport_http.RouteMeta

	Tag         string
	Description string

	NonBodyParameters map[string]*oas.Parameter
	RequestBody       *oas.RequestBody

	StatusErrors      []*status_error.StatusError
	StatusErrorSchema *oas.Schema

	SuccessStatus   int
	SuccessType     types.Type
	SuccessResponse *oas.Response
}

func (operator *Operator) AddNonBodyParameter(parameter *oas.Parameter) {
	if operator.NonBodyParameters == nil {
		operator.NonBodyParameters = map[string]*oas.Parameter{}
	}
	operator.NonBodyParameters[parameter.Name] = parameter
}

func (operator *Operator) SetRequestBody(requestBody *oas.RequestBody) {
	operator.RequestBody = requestBody
}

func (operator *Operator) BindOperation(method string, operation *oas.Operation, last bool) {
	parameterNames := map[string]bool{}
	for _, parameter := range operation.Parameters {
		parameterNames[parameter.Name] = true
	}

	for _, parameter := range operator.NonBodyParameters {
		if !parameterNames[parameter.Name] {
			operation.Parameters = append(operation.Parameters, parameter)
		}
	}

	if operator.RequestBody != nil {
		operation.SetRequestBody(operator.RequestBody)
	}

	for _, statusError := range operator.StatusErrors {
		statusErrorList := make([]string, 0)

		if operation.Responses.Responses != nil {
			if resp, ok := operation.Responses.Responses[statusError.Status()]; ok {
				if resp.Extensions != nil {
					if v, ok := resp.Extensions[XStatusErrs]; ok {
						if list, ok := v.([]string); ok {
							statusErrorList = append(statusErrorList, list...)
						}
					}
				}
			}
		}

		statusErrorList = append(statusErrorList, statusError.String())

		sort.Strings(statusErrorList)

		resp := oas.NewResponse("")
		resp.AddExtension(XStatusErrs, statusErrorList)
		resp.AddContent(httpx.MIME_JSON, oas.NewMediaTypeWithSchema(operator.StatusErrorSchema))
		operation.AddResponse(statusError.Status(), resp)
	}

	if last {
		operation.OperationId = operator.ID
		operation.Deprecated = operator.Deprecated
		operation.Summary = operator.Summary
		operation.Description = operator.Description

		if operator.Tag != "" {
			operation.Tags = []string{operator.Tag}
		}

		if operator.SuccessType == nil {
			operation.Responses.AddResponse(http.StatusNoContent, &oas.Response{})
		} else {
			status := operator.SuccessStatus
			if status == 0 {
				status = http.StatusOK
				if method == http.MethodPost {
					status = http.StatusCreated
				}
			}
			if status >= http.StatusMultipleChoices && status < http.StatusBadRequest {
				operator.SuccessResponse = oas.NewResponse(operator.SuccessResponse.Description)
			}
			operation.Responses.AddResponse(status, operator.SuccessResponse)
		}
	}

	// sort all parameters by postion and name
	if len(operation.Parameters) > 0 {
		sort.Slice(operation.Parameters, func(i, j int) bool {
			return positionOrders[operation.Parameters[i].In]+operation.Parameters[i].Name <
				positionOrders[operation.Parameters[j].In]+operation.Parameters[j].Name
		})
	}
}
