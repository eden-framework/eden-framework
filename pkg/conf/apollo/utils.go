package apollo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"
)

var (
	internalIp string
)

//ips
func getInternal() string {
	if internalIp != "" {
		return internalIp
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				internalIp = ipnet.IP.To4().String()
				return internalIp
			}
		}
	}
	return ""
}

func getCurrentTimeMillis() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
}

func url2PathWithQuery(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	pathWithQuery := u.Path

	if len(u.RawQuery) > 0 {
		pathWithQuery += "?" + u.RawQuery
	}
	return pathWithQuery
}

func getSignature(timestamp, path, secret string) string {
	source := fmt.Sprintf("%s\n%s", timestamp, path)
	hasher := hmac.New(sha1.New, []byte(secret))
	hasher.Write([]byte(source))
	hash := hasher.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}
