package apollo

import (
	"encoding/json"
)

type ApolloAppConfig struct {
	AppId          string `json:"appId"`
	Cluster        string `json:"cluster"`
	NamespaceName  string `json:"namespaceName"`
	ReleaseKey     string `json:"releaseKey"`
	Configurations struct {
		Content string `json:"content"`
	} `json:"configurations"`
	InputConf []interface{}
}

func (ac *ApolloAppConfig) Unmarshal(rawBytes []byte) error {
	if err := json.Unmarshal(rawBytes, ac); err != nil {
		return err
	}

	for _, c := range ac.InputConf {
		if err := json.Unmarshal([]byte(ac.Configurations.Content), c); err != nil {
			return err
		}
	}

	return nil
}
