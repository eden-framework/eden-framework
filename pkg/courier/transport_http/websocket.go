package transport_http

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/profzone/eden-framework/pkg/strings"
)

type IWebSocket interface {
	SubscribeOn(conn *websocket.Conn)
}

type TypeAndListener struct {
	Listener Listener
	Type     reflect.Type
}

func typeIndirect(tpe reflect.Type) reflect.Type {
	for tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}
	return tpe
}

type Listener func(v interface{}, ws *WSClient) error

type Listeners map[string]TypeAndListener

func (listeners Listeners) On(v interface{}, listener Listener) Listeners {
	tpe := typeIndirect(reflect.TypeOf(v))
	listeners[tpe.Name()] = TypeAndListener{
		Listener: listener,
		Type:     tpe,
	}
	return listeners
}

func (listeners Listeners) SubscribeOn(conn *websocket.Conn) {
	defer conn.Close()
	ws := ConnWS(conn)
	logger := logrus.WithField("remote_addr", conn.RemoteAddr().String())

	for !ws.Closed {
		msg := Msg{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Errorf("%s", err.Error())
			}
			break
		}
		logrus.Debugf(
			`{"type":"%s","payload": %s}`,
			color.RedString(msg.Type),
			color.GreenString("%s", strconv.Quote(string(msg.Payload))),
		)

		if typeAndListener, exists := listeners[msg.Type]; exists {
			m := reflect.New(typeAndListener.Type).Interface()
			errForUnmarshal := PayloadUnmarshal(msg.Payload, m)
			if errForUnmarshal != nil {
				logger.Errorf("%s", errForUnmarshal)
				break
			}
			errForOnMsg := typeAndListener.Listener(m, ws)
			if errForOnMsg != nil {
				logger.Errorf("%s", errForOnMsg)
				break
			}
		}
	}
}

type Msg struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

func ConnWS(conn *websocket.Conn) *WSClient {
	return &WSClient{conn: conn}
}

type WSClient struct {
	conn   *websocket.Conn
	Closed bool
}

func (c *WSClient) Close() error {
	defer c.conn.Close()
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (c WSClient) Send(payload interface{}) (err error) {
	msg := &Msg{}
	msg.Type = typeIndirect(reflect.TypeOf(payload)).Name()
	data, errForMarshal := PayloadMarshal(payload)
	if errForMarshal != nil {
		return errForMarshal
	}
	msg.Payload = data
	return c.conn.WriteJSON(msg)
}

func PayloadMarshal(payload interface{}) ([]byte, error) {
	tpe := typeIndirect(reflect.TypeOf(payload))
	if tpe.Kind() == reflect.Struct || tpe.Kind() == reflect.Array || tpe.Kind() == reflect.Slice {
		return json.Marshal(payload)
	}
	if marshaler, ok := payload.(json.Marshaler); ok {
		return marshaler.MarshalJSON()
	}
	return []byte(fmt.Sprintf("%v", payload)), nil
}

func PayloadUnmarshal(data []byte, payload interface{}) error {
	tpe := typeIndirect(reflect.TypeOf(payload))
	if tpe.Kind() == reflect.Struct || tpe.Kind() == reflect.Array || tpe.Kind() == reflect.Slice {
		return json.Unmarshal(data, payload)
	}
	if unmarshaler, ok := payload.(json.Unmarshaler); ok {
		return unmarshaler.UnmarshalJSON(data)
	}
	rv := reflect.ValueOf(payload)
	return str.ConvertFromStr(string(data), rv)
}
