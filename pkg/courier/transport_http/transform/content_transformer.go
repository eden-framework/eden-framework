package transform

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"

	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func init() {
	RegisterContentTransformer(&ContentTransformer{
		Key:         "json",
		ContentType: httpx.MIMEJSON,
		Marshal:     json.Marshal,
		Unmarshal: func(data []byte, v interface{}) error {
			d := json.NewDecoder(bytes.NewBuffer(data))

			err := d.Decode(v)
			if err != nil {
				switch err.(type) {
				case *json.UnmarshalTypeError:
					return err
				case *json.SyntaxError:
					return err
				default:
					offset := reflect.ValueOf(d).Elem().Field(2 /*d*/).Field(1 /*off*/).Int()
					return &json.UnmarshalTypeError{
						Offset: int64(offset),
						Value:  err.Error(),
					}
				}
				return err
			}
			return nil
		},
	})
	RegisterContentTransformer(&ContentTransformer{
		Key:         "xml",
		ContentType: httpx.MIMEXML,
		Marshal:     xml.Marshal,
		Unmarshal:   xml.Unmarshal,
	})
}

var contentTransformers = map[string]*ContentTransformer{}

func RegisterContentTransformer(contentTransformer *ContentTransformer) {
	contentTransformers[contentTransformer.Key] = contentTransformer
	contentTransformers[contentTransformer.ContentType] = contentTransformer
}

func GetContentTransformer(keyOrContentType string) *ContentTransformer {
	return contentTransformers[keyOrContentType]
}

type ContentTransformer struct {
	Key         string
	ContentType string
	Marshal     func(v interface{}) ([]byte, error)
	Unmarshal   func(data []byte, v interface{}) error
}
