package client

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/profzone/eden-framework/pkg/conf"
	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"github.com/profzone/eden-framework/pkg/courier/transport_grpc"
	"github.com/profzone/eden-framework/pkg/courier/transport_http"
	"github.com/profzone/eden-framework/pkg/env"
	"github.com/profzone/eden-framework/pkg/log/context"
)

type Client struct {
	Name string
	// used in service
	Service       string
	Version       string
	Host          string `conf:"upstream" validate:"@hostname"`
	Mode          string
	Port          int16
	Timeout       time.Duration
	WrapTransport transport_http.TransportWrapper `json:"-"`
}

func (c Client) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		// todo make switch in docker or expose
		"Host":    conf.RancherInternal("service-"+c.Name, "service-"+c.Name),
		"Timeout": 5 * time.Second,
		"Port":    80,
	}
}

func (Client) MarshalDefaults(v interface{}) {
	if client, ok := v.(*Client); ok {
		if client.Service == "" {
			client.Service = os.Getenv("PROJECT_NAME")
		}
		if client.Version == "" {
			client.Version = os.Getenv("PROJECT_REF")
		}
		if client.Mode == "" {
			client.Mode = "http"
		}
		if client.Host == "" {
			client.Host = fmt.Sprintf("service-%s.staging.g7pay.net", client.Name)
		}
		if client.Port == 0 {
			client.Port = 80
		}
		if client.Timeout == 0 {
			client.Timeout = 5 * time.Second
		}
	}
}

func (c Client) GetBaseURL(protocol string) (url string) {
	url = c.Host
	if protocol != "" {
		url = fmt.Sprintf("%s://%s", protocol, c.Host)
	}
	if c.Port > 0 {
		url = fmt.Sprintf("%s:%d", url, c.Port)
	}
	return
}

func (c *Client) Request(id, httpMethod, uri string, req interface{}, metas ...courier.Metadata) IRequest {
	requestID := context.GetLogID()
	metadata := courier.MetadataMerge(metas...)

	if !env.IsOnline() {
		if requestIDInMeta := metadata.Get(httpx.HeaderRequestID); requestIDInMeta != "" {
			requestID = requestIDInMeta
		}
		mocker, err := ParseMockID(c.Service, requestID)
		if err == nil {
			if m, exists := mocker.Mocks[id]; exists {
				logrus.Errorf("mocking %s with %s", id, m)

				return &MockRequest{
					MockData: m,
				}
			}
		}
	}

	if metadata.Has(courier.VersionSwitchKey) {
		requestID = courier.ModifyRequestIDWithVersionSwitch(requestID, metadata.Get(courier.VersionSwitchKey))
	} else {
		if _, v, exists := courier.ParseVersionSwitch(requestID); exists {
			metadata.Set(courier.VersionSwitchKey, v)
		}
	}

	if requestID == "" {
		requestID = uuid.New().String()
	}

	metadata.Add(httpx.HeaderRequestID, requestID)
	metadata.Add(httpx.HeaderUserAgent, c.Service+" "+c.Version)

	switch strings.ToLower(c.Mode) {
	case "grpc":
		serverName, method := parseID(id)
		return &transport_grpc.GRPCRequest{
			BaseURL:    c.GetBaseURL(""),
			ServerName: serverName,
			Method:     method,
			Timeout:    c.Timeout,
			RequestID:  requestID,
			Req:        req,
			Metadata:   metadata,
		}
	default:
		return &transport_http.HttpRequest{
			BaseURL:       c.GetBaseURL(c.Mode),
			Method:        httpMethod,
			URI:           uri,
			ID:            id,
			Timeout:       c.Timeout,
			WrapTransport: c.WrapTransport,
			Req:           req,
			Metadata:      metadata,
		}
	}
}

func parseID(id string) (serverName string, method string) {
	values := strings.Split(id, ".")
	if len(values) == 2 {
		serverName = strings.ToLower(strings.Replace(values[0], "Client", "", -1))
		method = values[1]
	}
	return
}
