package transport_http

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
)

func GetParams(path string, url string) (params httprouter.Params, err error) {
	pathArr := strings.Split(httprouter.CleanPath(path), "/")
	urlArr := strings.Split(httprouter.CleanPath(url), "/")

	if len(pathArr) != len(urlArr) {
		return nil, fmt.Errorf("url %s is not match path %s", url, path)
	}

	for i, p := range pathArr {
		if strings.HasPrefix(p, ":") {
			params = append(params, httprouter.Param{
				Key:   p[1:],
				Value: urlArr[i],
			})
		}
	}

	return params, nil
}

func GetClientIP(r *http.Request) string {
	clientIP := r.Header.Get(httpx.HeaderForwardedFor)

	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)

	if len(clientIP) > 0 {
		return clientIP
	}

	clientIP = strings.TrimSpace(r.Header.Get(httpx.HeaderRealIP))

	if len(clientIP) > 0 {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
