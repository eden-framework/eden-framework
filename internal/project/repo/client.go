package repo

import (
	"github.com/eden-framework/eden-framework/pkg/courier/client"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	tagUri      string
	repoListUri string

	*client.Client
}

func NewClient(mode, host string, port int16) *Client {
	cli := &Client{
		tagUri:      "/repos/eden-framework/eden-framework/tags",
		repoListUri: "/orgs/eden-framework/repos",
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

func (c *Client) GetPlugins() (RepoResponse, error) {
	var result RepoResponse
	request := c.Request("", http.MethodGet, c.repoListUri, nil)
	err := request.Do().Into(&result)
	if err != nil {
		return nil, err
	}

	var resp RepoResponse
	for _, r := range result {
		if strings.HasPrefix(r.Name, "plugin-") {
			resp = append(resp, r)
		}
	}

	return resp, nil
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
