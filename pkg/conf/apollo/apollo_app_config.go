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
	InputConf interface{}
}

func (ac *ApolloAppConfig) Unmarshal(rawBytes []byte) error {
	if err := json.Unmarshal(rawBytes, ac); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(ac.Configurations.Content), ac.InputConf); err != nil {
		return err
	}

	return nil
}
