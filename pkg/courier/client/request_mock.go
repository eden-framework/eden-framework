package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eden-framework/eden-framework/pkg/courier"
	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
)

func MetadataWithMocker(mocker *Mocker) courier.Metadata {
	m := courier.Metadata{}
	m.Set(httpx.HeaderRequestID, mocker.RequestID())
	return m
}

type MockRequest struct {
	*MockData
	courier.Result
}

func (mock *MockRequest) Do() courier.Result {
	r := courier.Result{}
	r.Unmarshal = json.Unmarshal
	r.Data = mock.MockData.Data
	r.Meta = mock.MockData.Meta
	if mock.MockData.Error != nil {
		statusErr := &status_error.StatusError{}
		err := json.Unmarshal(mock.MockData.Error, &statusErr)
		if err == nil {
			r.Err = statusErr
		}
	}
	return r
}

func ParseMockID(service string, requestID string) (mock *Mocker, err error) {
	requestIDs := strings.Split(requestID, ";")
	mock = Mock(service)
	for _, requestID := range requestIDs {
		prefix := service + ":"
		if strings.HasPrefix(requestID, prefix) {
			mock, err = mock.From(strings.Replace(requestID, prefix, "", 1))
			return
		}
	}
	return nil, fmt.Errorf("no mock")
}

func Mock(service string) *Mocker {
	return &Mocker{
		Service: service,
		Mocks:   map[string]*MockData{},
	}
}

type MockData struct {
	Data  []byte           `json:"data,omitempty"`
	Error []byte           `json:"error,omitempty"`
	Meta  courier.Metadata `json:"metadata,omitempty"`
}

func (data *MockData) String() string {
	if data.Data != nil {
		return fmt.Sprintf("return: %s", data.Data)
	}
	if data.Error != nil {
		return fmt.Sprintf("error: %s", data.Error)
	}
	return ""
}

type Mocker struct {
	Service string
	Mocks   map[string]*MockData
}

func (mocker Mocker) From(mock string) (*Mocker, error) {
	pair := strings.Split(mock, ":")
	if len(pair) != 2 {
		return nil, fmt.Errorf("invalid request id")
	}
	data, err := base64.StdEncoding.DecodeString(pair[1])
	if err != nil {
		return nil, err
	}
	m := &MockData{}
	errForUnmarshal := json.Unmarshal(data, &m)
	if errForUnmarshal != nil {
		return nil, errForUnmarshal
	}
	return mocker.For(pair[0], m), nil
}

func (mocker Mocker) For(methodID string, m *MockData) *Mocker {
	mocker.Mocks[methodID] = m
	return &mocker
}

func (mocker *Mocker) RequestID() string {
	buf := new(bytes.Buffer)
	for id, m := range mocker.Mocks {
		data, _ := json.Marshal(m)
		buf.WriteString(fmt.Sprintf("%s:%s:%s;", mocker.Service, id, base64.StdEncoding.EncodeToString(data)))
	}
	return buf.String()
}
