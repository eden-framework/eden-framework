package generator

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"profzone/eden-framework/internal/generator/api"
)

type ClientGenerator struct {
	Api         api.Api
	ServiceName string
}

func (c *ClientGenerator) Load(path string) {
	url, err := url.Parse(path)
	if err != nil {
		logrus.Panic(err)
	}
	if url.Scheme == "" {
		// 本地文件
		err = c.loadLocalFile(url.Path)
	} else if url.Scheme == "http" || url.Scheme == "https" {
		// 远程文件
		err = c.loadRemoteFile(url.Path)
	} else {
		logrus.Panicf("unsupported scheme %s", url.Scheme)
	}
	if err != nil {
		logrus.Panic(err)
	}

	if c.ServiceName == "" {
		c.ServiceName = "client-" + c.Api.ServiceName
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
	packageName := ToLowerSnakeCase(c.ServiceName)

	outputs.Add(GeneratedSuffix(path.Join(packageName, "client.go")), "package "+packageName)
	outputs.Add(GeneratedSuffix(path.Join(packageName, "types.go")), "package "+packageName)
	outputs.Add(GeneratedSuffix(path.Join(packageName, "enums.go")), "package "+packageName)

	return outputs
}
