package validatetpl

import (
	"net/url"
	"strings"
)

const (
	HTTP_URL_TYPE_ERROR   = "url类型错误"
	HTTP_URL_VALUE_ERROR  = "非法的http url"
	HTTP_URL_SCHEME_ERROR = "scheme错误"
	HTTP_URL_TOO_LONG     = "url超过了256个字节"
)

const (
	MAX_URL_LEN = 256
)

func ValidateHttpUrl(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, HTTP_URL_TYPE_ERROR
	}

	if len(s) > MAX_URL_LEN {
		return false, HTTP_URL_TOO_LONG
	}

	if u, err := url.Parse(s); err != nil {
		return false, HTTP_URL_VALUE_ERROR
	} else {
		scheme := strings.ToLower(u.Scheme)
		if scheme != "http" && scheme != "https" {
			return false, HTTP_URL_SCHEME_ERROR
		}
	}
	return true, ""
}

func ValidateHttpUrlOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, HTTP_URL_TYPE_ERROR
	}

	if s == "" {
		return true, ""
	}

	return ValidateHttpUrl(v)
}
