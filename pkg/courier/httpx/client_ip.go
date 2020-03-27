package httpx

import (
	"strings"
)

func GetClientIPByHeaderForwardedFor(headerForwardedFor string) string {
	if index := strings.IndexByte(headerForwardedFor, ','); index >= 0 {
		return headerForwardedFor[0:index]
	}
	return strings.TrimSpace(headerForwardedFor)
}

func GetClientIPByHeaderRealIP(headerRealIP string) string {
	return strings.TrimSpace(headerRealIP)
}
