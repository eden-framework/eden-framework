package transport_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"github.com/profzone/eden-framework/pkg/courier/status_error"
	"github.com/profzone/eden-framework/pkg/courier/transport_http/transform"
	"github.com/profzone/eden-framework/pkg/duration"
)

type TransportWrapper func(rt http.RoundTripper) http.RoundTripper

type HttpRequest struct {
	BaseURL       string
	Method        string
	URI           string
	ID            string
	Metadata      courier.Metadata
	Timeout       time.Duration
	Req           interface{}
	WrapTransport TransportWrapper
}

func (httpRequest *HttpRequest) Do() (result courier.Result) {
	result = courier.Result{}
	request, err := transform.NewRequest(httpRequest.Method, httpRequest.BaseURL+httpRequest.URI, httpRequest.Req, httpRequest.Metadata)
	if err != nil {
		result.Err = status_error.InvalidRequestParams.StatusError().WithDesc(err.Error())
		return
	}

	d := duration.NewDuration()
	defer func() {
		logger := d.ToLogger().WithFields(logrus.Fields{
			"request":  httpRequest.ID,
			"method":   httpRequest.Method,
			"url":      request.URL.String(),
			"metadata": httpRequest.Metadata,
		})

		if result.Err == nil {
			logger.Debugf("success")
		} else {
			logger.Warnf("do http request failed %s", result.Err)
		}
	}()

	httpClient := GetShortConnClient(httpRequest.Timeout)

	if httpRequest.WrapTransport != nil {
		httpClient.Transport = httpRequest.WrapTransport(httpClient.Transport)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		result.Err = status_error.RequestTimeout.StatusError().WithDesc(err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	result.Data = data
	result.Meta = courier.Metadata(resp.Header)

	contentType := resp.Header.Get(httpx.HeaderContentType)

	if strings.Contains(contentType, httpx.MIMEJSON) {
		result.Unmarshal = json.Unmarshal
		// todo add more structUnmarshal
	}

	if !HttpStatusOK(resp.StatusCode) {
		statusError := &status_error.StatusError{}
		err := json.Unmarshal(result.Data, statusError)
		if err != nil {
			msg := fmt.Sprintf("[%d] %s %s", resp.StatusCode, request.Method, request.URL)
			result.Err = status_error.HttpRequestFailed.StatusError().WithDesc(msg)
			return
		}
		result.Err = statusError
	}
	return
}

func HttpStatusOK(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}

func GetShortConnClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 0,
			}).Dial,
			DisableKeepAlives: true,
		},
	}
}

func GetLongConnClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 0,
			}).Dial,
			DisableKeepAlives: false,
		},
	}
}
