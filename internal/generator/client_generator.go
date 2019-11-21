package generator

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"profzone/eden-framework/internal/generator/api"
)

type ClientGenerator struct {
	Api api.Api
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
	return Outputs{}
}
