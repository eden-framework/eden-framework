package apollo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func NewHttpUtil(secretKey, appID, method, url string, timeout time.Duration, requestBody []byte) *HttpUtil {
	httpUtil := new(HttpUtil)
	httpUtil.Method = method
	httpUtil.Url = url
	httpUtil.RequestBody = requestBody
	httpUtil.Timeout = timeout
	httpUtil.appID = appID
	httpUtil.secretKey = secretKey
	return httpUtil
}

type HttpUtil struct {
	Method      string
	Url         string
	RequestBody []byte
	Timeout     time.Duration

	appID     string
	secretKey string
}

func (hu *HttpUtil) Request() ([]byte, int, error) {
	defer PrintDuration(map[string]interface{}{
		"Request": hu.Url,
	})()

	if hu.Timeout == 0 {
		// set default timeout
		hu.Timeout = 5 * time.Second
	}

	req, err := http.NewRequest(hu.Method, hu.Url, bytes.NewBuffer(hu.RequestBody))
	if err != nil {
		logrus.Warningf("NewRequest fail![err:%v]", err)
		return []byte(""), -1, err
	}

	req.Header.Add("Content-type", "application/json;charset=UTF-8")

	if hu.secretKey != "" {
		var timestamp = getCurrentTimeMillis()
		req.Header.Add("Authorization", fmt.Sprintf("Apollo %s:%s", hu.appID, getSignature(timestamp, url2PathWithQuery(hu.Url), hu.secretKey)))
		req.Header.Add("Timestamp", timestamp)
	}

	client := &http.Client{
		Timeout: hu.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Do Request fail![err:%v], request[%s], resp[%+v]", err, string(hu.RequestBody), resp)
		return []byte(""), -1, err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Warningf("ReadAll failed![err:%v]", err)
		return []byte(""), -1, err
	}

	if resp.StatusCode/http.StatusOK != 1 {
		if resp.StatusCode/http.StatusInternalServerError == 1 {
			logrus.Errorf("Request return fail![err:%v], request[%s], resp[%+v]", err, string(hu.RequestBody), resp)
		} else {
			logrus.Warnf("Request return fail![err:%v], request[%s], resp[%+v]", err, string(hu.RequestBody), resp)
		}
	} else {
		logrus.Debugf("apollo result[%s], request[%s], httpCoce[%d]", string(respData), string(hu.RequestBody),
			resp.StatusCode)
	}

	return respData, resp.StatusCode, nil
}

func (hu *HttpUtil) Is2xxCode(httpCode int) bool {
	return httpCode/http.StatusContinue == 2
}
