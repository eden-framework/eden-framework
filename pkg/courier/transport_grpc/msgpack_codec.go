package transport_grpc

import (
	"reflect"

	"github.com/vmihailenco/msgpack"
)

type MsgPackCodec struct {
}

func (c *MsgPackCodec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func MsgPackUnmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (c *MsgPackCodec) Unmarshal(data []byte, v interface{}) error {
	// delay Unmarshal
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if _, ok := rv.Interface().([]byte); ok {
		rv.SetBytes(data)
	}
	return nil
}

func (c *MsgPackCodec) String() string {
	return "MsgPackCodec"
}
