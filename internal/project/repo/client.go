package repo

import (
	"github.com/eden-framework/eden-framework/pkg/courier/client"
	"net/http"
	"time"
)

type Client struct {
	tagUri string
	*client.Client
}

func NewClient(mode, host string, port int16) *Client {
	cli := &Client{
		tagUri: "/repos/profzone/eden-framework/tags",
		Client: &client.Client{
			Host:    host,
			Mode:    mode,
			Port:    port,
			Timeout: time.Minute,
		},
	}
	cli.MarshalDefaults(cli.Client)
	return cli
}

func (c *Client) GetTags() (TagsResponse, error) {
	var resp TagsResponse
	request := c.Request("", http.MethodGet, c.tagUri, nil)
	err := request.Do().Into(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
