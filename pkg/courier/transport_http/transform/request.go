package transform

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"github.com/profzone/eden-framework/pkg/courier/status_error"
)

func NewRequest(method string, uri string, v interface{}, metadatas ...courier.Metadata) (*http.Request, error) {
	if v == nil {
		return NewHttpRequestFromParameterGroup(method, uri, nil, metadatas...)
	}
	return NewHttpRequestFromParameterGroup(method, uri, ParameterGroupFromValue(v), metadatas...)
}

func CloneRequestBody(req *http.Request) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

func NewHttpRequestFromParameterGroup(method string, uri string, m *ParameterGroup, metadatas ...courier.Metadata) (req *http.Request, err error) {
	body := &bytes.Buffer{}
	multipartWriter := (*multipart.Writer)(nil)

	if m != nil {
		if m.FormData.Len() > 0 {
			postForm := (*url.Values)(nil)

			for _, parameterMeta := range m.FormData {
				name := parameterMeta.Name
				dataList, errForMarshal := parameterMeta.Marshal()

				if errForMarshal != nil {
					err = errForMarshal
					return
				}

				if parameterMeta.IsMultipart() {
					if multipartWriter == nil {
						multipartWriter = multipart.NewWriter(body)
					}

					if fileHeaders, ok := parameterMeta.Value.Interface().([]*multipart.FileHeader); ok {
						for _, fileHeader := range fileHeaders {
							if errForAppend := appendFile(multipartWriter, name, fileHeader, nil); errForAppend != nil {
								err = errForAppend
								return
							}
						}
						continue
					}

					if fileHeader, ok := parameterMeta.Value.Interface().(*multipart.FileHeader); ok {
						if errForAppend := appendFile(multipartWriter, name, fileHeader, nil); errForAppend != nil {
							err = errForAppend
							return
						}
						continue
					}

					for _, v := range dataList {
						err = multipartWriter.WriteField(name, string(v))
						if err != nil {
							return
						}
					}
				} else {
					if postForm == nil {
						postForm = &url.Values{}
					}
					for _, v := range dataList {
						postForm.Add(name, string(v))
					}
				}
			}
			if postForm != nil {
				body.WriteString(postForm.Encode())
			}
		} else if m.Body != nil {
			dataList, errForMarshal := m.Body.Marshal()
			if errForMarshal != nil {
				err = errForMarshal
				return
			}
			if len(dataList) > 0 {
				body.Write(dataList[0])
			}
		}
	}
	if multipartWriter != nil {
		multipartWriter.Close()
	}

	u, errForParseUrl := url.Parse(uri)
	if errForParseUrl != nil {
		err = errForParseUrl
		return
	}

	if u.Path == "" {
		u.Path = "/"
	}

	req, err = http.NewRequest(
		method,
		u.String(),
		body,
	)

	if err != nil {
		return
	}

	req.Close = true

	if metadatas != nil {
		for key, values := range courier.MetadataMerge(metadatas...) {
			for _, v := range values {
				req.Header.Add(key, v)
			}
		}
	}

	if m == nil {
		return
	}

	for _, parameterMeta := range m.Parameters {
		name := parameterMeta.Name
		dataList, errForMarshal := parameterMeta.Marshal()
		if errForMarshal != nil {
			err = errForMarshal
			return
		}

		if len(dataList) == 0 {
			continue
		}

		switch parameterMeta.In {
		case "query":
			// todo use explode when service updated
			query := url.Values{}
			query.Add(parameterMeta.Name, string(bytes.Join(dataList, []byte(","))))

			if req.URL.RawQuery == "" {
				req.URL.RawQuery = query.Encode()
			} else {
				req.URL.RawQuery = req.URL.RawQuery + "&" + query.Encode()
			}
		case "cookie":
			for _, v := range dataList {
				req.AddCookie(&http.Cookie{
					Name:  name,
					Value: string(v),
				})
			}
		case "path":
			paramKey := ":" + name
			if !strings.Contains(req.URL.Path, paramKey) {
				err = fmt.Errorf("uri %s need path parameter %s, but uri has no %s", uri, name, paramKey)
				return
			}
			value := string(bytes.Join(dataList, []byte(",")))
			req.URL.Path = strings.Replace(req.URL.Path, paramKey, value, -1)
		case "header":
			value := string(bytes.Join(dataList, []byte(",")))
			req.Header.Add(name, value)
		}
	}

	if m.FormData.Len() > 0 {
		if multipartWriter != nil {
			req.Header.Add(httpx.HeaderContentType, multipartWriter.FormDataContentType())
		} else {
			req.Header.Add(httpx.HeaderContentType, httpx.MIME_POST_URLENCODED+"; param=value")
		}
	}

	return
}

func appendFile(multipartWriter *multipart.Writer, field string, fileHeader *multipart.FileHeader, content []byte) (err error) {
	if fileHeader == nil {
		return
	}
	filePart, errForCreate := multipartWriter.CreateFormFile(field, fileHeader.Filename)
	if errForCreate != nil {
		err = errForCreate
		return
	}

	if len(content) == 0 {
		if file, errForOpen := fileHeader.Open(); errForOpen != nil {
			err = errForOpen
			return
		} else if _, errForCopy := io.Copy(filePart, file); errForCopy != nil {
			err = errForCopy
			return
		}
		return nil
	}

	if _, errForWrite := filePart.Write(content); errForWrite != nil {
		err = errForWrite
		return
	}
	return nil
}

func NewFileHeader(fieldName string, filename string, content []byte) (fileHeader *multipart.FileHeader, err error) {
	buffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(buffer)
	appendFile(multipartWriter, fieldName, &multipart.FileHeader{
		Filename: filename,
	}, content)
	multipartWriter.Close()
	reader := multipart.NewReader(buffer, multipartWriter.Boundary())
	form, errForReader := reader.ReadForm(int64(buffer.Len()))
	if errForReader != nil {
		err = errForReader
		return
	}
	fileHeader = form.File[fieldName][0]
	return
}

type IParameterValuesGetter interface {
	Initial() error
	Param(name string) string
	Query(name string) []string
	Header(name string) []string
	Cookie(name string) []string
	FormValue(name string) []string
	FormFile(name string) ([]*multipart.FileHeader, error)
	Body() io.Reader
}

func NewParameterValuesGetter(r *http.Request) *ParameterValuesGetter {
	return &ParameterValuesGetter{Request: r}
}

type ParameterValuesGetter struct {
	Request *http.Request
	query   url.Values
	now     time.Time
	cookies []*http.Cookie
}

func (getter *ParameterValuesGetter) Initial() error {
	getter.now = time.Now()
	getter.query = getter.Request.URL.Query()
	return nil
}

func (getter *ParameterValuesGetter) Param(name string) string {
	return ""
}

func (getter *ParameterValuesGetter) Query(name string) []string {
	return getter.query[name]
}

func (getter *ParameterValuesGetter) Header(name string) []string {
	return getter.Request.Header[textproto.CanonicalMIMEHeaderKey(name)]
}

func (getter *ParameterValuesGetter) Cookie(name string) []string {
	values := make([]string, 0)
	if len(getter.cookies) == 0 {
		getter.cookies = getter.Request.Cookies()
	}
	for _, c := range getter.cookies {
		if c.Name == name {
			if c.Expires.IsZero() {
				values = append(values, c.Value)
			} else if c.Expires.After(getter.now) {
				values = append(values, c.Value)
			}
		}
	}
	return values
}

func (getter *ParameterValuesGetter) FormValue(name string) []string {
	if getter.Request.Form == nil {
		// just parse form
		getter.Request.FormValue("")
	}
	return getter.Request.Form[name]
}

func (getter *ParameterValuesGetter) FormFile(name string) ([]*multipart.FileHeader, error) {
	if getter.Request.MultipartForm == nil {
		// just parse form
		_, fileHeader, err := getter.Request.FormFile(name)
		if err != nil || fileHeader == nil {
			return nil, err
		}
	}
	return getter.Request.MultipartForm.File[name], nil
}

func (getter *ParameterValuesGetter) Body() io.Reader {
	return getter.Request.Body
}

func MarshalParameters(group *ParameterGroup, getter IParameterValuesGetter) error {
	parameterErrors := ParameterErrors{}

	if err := getter.Initial(); err != nil {
		parameterErrors.Merge(err)
		return parameterErrors.Err()
	}

	for _, parameterMeta := range group.Parameters.List() {
		switch parameterMeta.In {
		case "path":
			parameterErrors.Merge(parameterMeta.UnmarshalStringAndValidate(strings.Split(getter.Param(parameterMeta.Name), ",")...))
		case "header":
			parameterErrors.Merge(parameterMeta.UnmarshalStringAndValidate(getter.Header(parameterMeta.Name)...))
		case "query":
			parameterErrors.Merge(parameterMeta.UnmarshalStringAndValidate(getter.Query(parameterMeta.Name)...))
		case "cookie":
			parameterErrors.Merge(parameterMeta.UnmarshalStringAndValidate(getter.Cookie(parameterMeta.Name)...))
		}
	}

	for _, parameterMeta := range group.FormData {
		if _, ok := parameterMeta.Value.Interface().([]*multipart.FileHeader); ok {
			fileHeaders, err := getter.FormFile(parameterMeta.Name)
			if err != nil || len(fileHeaders) == 0 {
				return status_error.ReadFormFileFailed
			}
			parameterMeta.Value.Set(reflect.ValueOf(fileHeaders))
			continue
		}
		if _, ok := parameterMeta.Value.Interface().(*multipart.FileHeader); ok {
			fileHeaders, err := getter.FormFile(parameterMeta.Name)
			if err != nil || len(fileHeaders) == 0 {
				return status_error.ReadFormFileFailed
			}
			if parameterMeta.Value.CanSet() {
				parameterMeta.Value.Set(reflect.ValueOf(fileHeaders[0]))
			}
			continue
		}
		parameterErrors.Merge(parameterMeta.UnmarshalStringAndValidate(getter.FormValue(parameterMeta.Name)...))
	}

	if parameterErrors.StatusError == nil {
		isValid, errFields := group.ValidateNoBodyByHook()
		if !isValid {
			parameterErrors.Merge(status_error.InvalidField.StatusError().WithErrorFields(errFields...))
		}
	}

	if group.Body != nil {
		parameterErrors.Merge(group.Body.UnmarshalFromReader(getter.Body()))
	}

	return parameterErrors.Err()
}
