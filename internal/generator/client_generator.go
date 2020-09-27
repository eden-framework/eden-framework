package generator

import (
	"encoding/json"
	"github.com/eden-framework/eden-framework/internal/generator/api"
	"github.com/eden-framework/eden-framework/internal/generator/files"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type ClientGenerator struct {
	Api         api.Api
	serviceName string
}

func NewClientGenerator(pkgName string) *ClientGenerator {
	return &ClientGenerator{
		serviceName: pkgName,
	}
}

func (c *ClientGenerator) Load(path string) {
	urlPattern, err := url.Parse(path)
	if err != nil {
		logrus.Panic(err)
	}
	if urlPattern.Scheme == "" {
		// 本地文件
		err = c.loadLocalFile(urlPattern.Path)
	} else if urlPattern.Scheme == "http" || urlPattern.Scheme == "https" {
		// 远程文件
		err = c.loadRemoteFile(urlPattern.Path)
	} else {
		logrus.Panicf("unsupported scheme %s", urlPattern.Scheme)
	}
	if err != nil {
		logrus.Panic(err)
	}

	if c.serviceName == "" {
		c.serviceName = c.Api.ServiceName
	}
}

func (c *ClientGenerator) loadLocalFile(path string) (err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c.Api)
	return
}

func (c *ClientGenerator) loadRemoteFile(path string) (err error) {
	resp, err := http.Get(path)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c.Api)
	return
}

func (c *ClientGenerator) Pick() {
}

func (c *ClientGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	clientFile := files.NewClientFile(c.serviceName, &c.Api)
	clientEnumsFile := files.NewClientEnumsFile(outputPath, c.serviceName, &c.Api)
	typesFile := files.NewTypesFile(c.serviceName, &c.Api)

	outputs.Add(GeneratedSuffix(path.Join(outputPath, clientFile.PackageName, "client.go")), clientFile.String())
	outputs.Add(GeneratedSuffix(path.Join(outputPath, clientFile.PackageName, "types.go")), typesFile.String())
	outputs.Add(GeneratedSuffix(path.Join(outputPath, clientFile.PackageName, "enums.go")), clientEnumsFile.String())

	return outputs
}
