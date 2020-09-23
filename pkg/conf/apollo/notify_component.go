package apollo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

type apolloNotify struct {
	NotificationId int64  `json:"notificationId"`
	NamespaceName  string `json:"namespaceName"`
}

type notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int64  `json:"notificationId"`
}

func NewNotifyComponent(namespace string, worker chan bool) *NotifyComponent {
	notifyComponent := new(NotifyComponent)
	notifyComponent.notifications = make(map[string]int64)
	notifyComponent.worker = worker

	notifyComponent.setNotify(namespace, -1)
	return notifyComponent
}

type NotifyComponent struct {
	notifications map[string]int64
	worker        chan bool
}

func (nc *NotifyComponent) setNotify(namespaceName string, notificationId int64) {
	nc.notifications[namespaceName] = notificationId
}

func (nc *NotifyComponent) getNotifies() string {
	notificationArr := make([]*notification, 0)
	for namespaceName, notificationId := range nc.notifications {
		notificationArr = append(notificationArr,
			&notification{
				NamespaceName:  namespaceName,
				NotificationId: notificationId,
			})
	}

	dstBytes, err := json.Marshal(notificationArr)
	if err != nil {
		logrus.Errorf("Marshal failed![%+v]", err)
		return ""
	}

	logrus.Debugf("notification %+v", string(dstBytes))
	return string(dstBytes)
}

func (nc *NotifyComponent) getNotifyUrl(apollo *Apollo) string {
	return fmt.Sprintf("%snotifications/v2?appId=%s&cluster=%s&notifications=%s",
		"http://"+apollo.Host+"/",
		url.QueryEscape(apollo.AppId),
		url.QueryEscape(apollo.Cluster),
		url.QueryEscape(nc.getNotifies()))
}

func (nc *NotifyComponent) dealNotifyRespByHttpCode(httpCode int, respBody []byte, apollo *Apollo) error {
	switch httpCode {
	case http.StatusOK:
		notifies := make([]*apolloNotify, 0)
		err := json.Unmarshal(respBody, &notifies)
		if err != nil {
			logrus.Error("Unmarshal apollo notify Fail,Error:", err)
			return err
		}

		for _, notify := range notifies {
			if notify.NamespaceName == "" {
				continue
			}

			nc.setNotify(notify.NamespaceName, notify.NotificationId)
		}

		if err := apollo.fetchConfigFromServer(); err != nil {
			return err
		}

		// notify for refresh config
		logrus.Infof("dealNotifyRespByHttpCode:%+v", httpCode)
		nc.worker <- true
	case http.StatusNotModified:
		return nil
	default:
		logrus.Errorf("Notify Error![code:%+v], [resp:%+v]", httpCode, string(respBody))
		return nil
	}

	return nil

}

var (
	// the interval for long pool
	PoolInterval = 1 * time.Second
	// the timeout for long pool
	LongPollTimeout = 1*time.Minute + 30*time.Second
)

func (nc *NotifyComponent) acceptNotify(apollo *Apollo) {
	timer := time.NewTimer(PoolInterval)
	httpUtil := NewHttpUtil(apollo.SecretKey, apollo.AppId, http.MethodGet, nc.getNotifyUrl(apollo), LongPollTimeout, nil)
	for {
		select {
		case <-timer.C:
			httpUtil.Url = nc.getNotifyUrl(apollo)
			respBody, httpCode, err := httpUtil.Request()
			if err != nil {
				logrus.Errorf("fetchNotifyFromServer failed![%+v]", err)
				timer.Reset(PoolInterval)
				continue
			}

			if httpCode/http.StatusOK != 1 {
				logrus.Errorf("fetchNotifyFromServer failed![httpCode:%+v, resp:%+v]", httpCode, string(respBody))
				timer.Reset(PoolInterval)
				continue
			}

			nc.dealNotifyRespByHttpCode(httpCode, respBody, apollo)
			timer.Reset(PoolInterval)
		}
	}
}
