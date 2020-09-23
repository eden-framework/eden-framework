package apollo

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

func NewApolloConfig(namespace string, defaultBaseConf ApolloBaseConfig, conf ...interface{}) *Apollo {
	apollo := &Apollo{
		ApolloBaseConfig: defaultBaseConf,
	}
	apollo.Cluster = defaultBaseConf.Cluster
	apollo.NamespaceName = namespace

	apollo.Config = new(ApolloAppConfig)
	apollo.Config.InputConf = conf

	return apollo
}

type ApolloBaseConfig struct {
	AppId            string `json:"appId"`
	Host             string `json:"host"`
	BackupConfigPath string `json:"backupConfigPath"`
	Cluster          string `json:"cluster"`
	SecretKey        string `json:"-"`
}

type Apollo struct {
	ApolloBaseConfig
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	Config        *ApolloAppConfig
}

func (a *Apollo) Start(worker chan bool) error {
	err := a.fetchConfigFromServer()
	if err != nil {
		configContent, err := loadConfigFile(a.assembleConfigFullPath())
		if err != nil {
			panic(fmt.Sprintf("load config file[%+v] failed![%+v]", a.assembleConfigFullPath(), err))
		}

		if len(configContent) == 0 {
			panic(fmt.Sprintf("config file[%+v] is empty!", a.assembleConfigFullPath()))
		}

		if err := a.assignToConf(configContent); err != nil {
			logrus.Errorf("assignToConfAndSaveTofile Unmarshal Fail![%+v]", err)
			return err
		}
	}

	a.printReleaseKey()

	//start accept apollo notify
	go NewNotifyComponent(a.NamespaceName, worker).acceptNotify(a)

	logrus.Info("apollo start finished!")
	return nil
}

func (a *Apollo) fetchConfigFromServer() error {
	httpUtil := NewHttpUtil(a.SecretKey, a.AppId, http.MethodGet, a.getConfigUrl(), 0, nil)
	respBody, httpCode, err := httpUtil.Request()
	if err != nil {
		logrus.Errorf("fetchConfigFromServer failed![%+v]", err)
		return err
	}

	if !httpUtil.Is2xxCode(httpCode) {
		if httpCode == http.StatusNotModified {
			// config not changed.
			return nil
		} else {
			logrus.Errorf("fetchConfigFromServer failed![httpCode:%+v, resp:%+v]", httpCode, string(respBody))
			return fmt.Errorf("fetchConfigFromServer failed![httpCode:%+v, resp:%+v]", httpCode, string(respBody))
		}
	}

	return a.assignToConfAndSaveTofile(respBody)
}

func (a *Apollo) assembleConfigFullPath() string {
	if a.ApolloBaseConfig.BackupConfigPath != "" {
		return fmt.Sprintf("%s/%s/%s/%s", a.ApolloBaseConfig.BackupConfigPath, a.AppId, a.Cluster, a.NamespaceName)
	} else {
		return fmt.Sprintf("%s/%s/%s", a.AppId, a.Cluster, a.NamespaceName)
	}
}

func (a *Apollo) assignToConf(respBody []byte) error {
	if err := a.Config.Unmarshal(respBody); err != nil {
		logrus.Errorf("assignToConf Unmarshal Fail![%+v]", err)
		return err
	}

	return nil
}

func (a *Apollo) assignToConfAndSaveTofile(respBody []byte) error {
	if err := a.assignToConf(respBody); err != nil {
		logrus.Errorf("assignToConfAndSaveTofile Unmarshal Fail![%+v]", err)
		return err
	}

	go writeConfigFile(respBody, a.assembleConfigFullPath())
	return nil
}

func (a *Apollo) getConfigUrl() string {
	return fmt.Sprintf("%sconfigs/%s/%s/%s?releaseKey=%s&ip=%s",
		"http://"+a.Host+"/",
		url.QueryEscape(a.AppId),
		url.QueryEscape(a.Cluster),
		url.QueryEscape(a.NamespaceName),
		url.QueryEscape(a.Config.ReleaseKey),
		getInternal())
}

func (a *Apollo) printReleaseKey() {
	logrus.Infof("apollo config release key: %s", a.Config.ReleaseKey)
}
