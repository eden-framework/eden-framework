package repo

import (
	"fmt"
	"gitee.com/eden-framework/courier/client"
	"github.com/profzone/envconfig"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	tagUri      string
	repoListUri string

	*client.Client
}

func NewClient(mode, host string, port int16, org string) *Client {
	cli := &Client{
		tagUri:      "/api/v5/repos/" + org + "/%s/tags",
		repoListUri: "/api/v5/orgs/" + org + "/repos",
		Client: &client.Client{
			Host:    host,
			Mode:    mode,
			Port:    port,
			Timeout: envconfig.Duration(time.Minute),
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

func (c *Client) GetTags(repoName string) (TagsResponse, error) {
	var resp TagsResponse
	request := c.Request("", http.MethodGet, fmt.Sprintf(c.tagUri, repoName), nil)
	err := request.Do().Into(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
