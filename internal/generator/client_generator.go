package generator

import (
	"bytes"
	"encoding/json"
	"github.com/eden-framework/courier/status_error"
	"github.com/eden-framework/eden-framework/internal/generator/files"
	"github.com/go-courier/oas"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
)

type ClientGenerator struct {
	Api              *oas.OpenAPI
	serviceName      string
	statusErrCodeMap status_error.StatusErrorCodeMap
}

func NewClientGenerator(pkgName string) *ClientGenerator {
	return &ClientGenerator{
		serviceName: pkgName,
	}
}

func (c *ClientGenerator) Load(path string) {
	if c.serviceName == "" {
		panic("[ClientGenerator] must specify a service name")
	}

	urlPattern, err := url.Parse(path)
	if err != nil {
		logrus.Panic(err)
	}
	if urlPattern.Scheme == "" {
		// 本地文件
		c.Api, c.statusErrCodeMap, err = c.loadLocalFile(urlPattern.Path)
	} else if urlPattern.Scheme == "http" || urlPattern.Scheme == "https" {
		// 远程文件
		c.Api, c.statusErrCodeMap, err = c.loadRemoteFile(urlPattern.String())
	} else {
		logrus.Panicf("unsupported scheme %s", urlPattern.Scheme)
	}
	if err != nil {
		logrus.Panic(err)
	}
}

func (c *ClientGenerator) loadLocalFile(path string) (api *oas.OpenAPI, statusErrCodeMap status_error.StatusErrorCodeMap, err error) {
	result, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	regexp.MustCompile("@httpError[^;]+;").ReplaceAllFunc(result, func(i []byte) []byte {
		v := bytes.Replace(i, []byte(`\"`), []byte(`"`), -1)
		s := status_error.ParseString(string(v))
		statusErrCodeMap[s.Code] = *s
		return i
	})

	api = new(oas.OpenAPI)
	err = json.Unmarshal(result, &api)
	return
}

func (c *ClientGenerator) loadRemoteFile(path string) (api *oas.OpenAPI, statusErrCodeMap status_error.StatusErrorCodeMap, err error) {
	resp, err := http.Get(path)
	if err != nil {
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	regexp.MustCompile("@httpError[^;]+;").ReplaceAllFunc(result, func(i []byte) []byte {
		v := bytes.Replace(i, []byte(`\"`), []byte(`"`), -1)
		s := status_error.ParseString(string(v))
		statusErrCodeMap[s.Code] = *s
		return i
	})

	api = new(oas.OpenAPI)
	err = json.Unmarshal(result, &api)
	return
}

func (c *ClientGenerator) Pick() {
}

func (c *ClientGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	clientFile := files.NewClientFile(c.serviceName, c.Api)
	typesFile := files.NewTypesFile(c.serviceName, c.Api)
	clientEnumsFile := files.NewClientEnumsFile(outputPath, c.serviceName, c.Api)

	outputs.Add(GeneratedSuffix(path.Join(outputPath, clientFile.PackageName, "client.go")), clientFile.String())
	outputs.Add(GeneratedSuffix(path.Join(outputPath, clientFile.PackageName, "types.go")), typesFile.String())
	outputs.Add(GeneratedSuffix(path.Join(outputPath, clientFile.PackageName, "enums.go")), clientEnumsFile.String())

	return outputs
}
