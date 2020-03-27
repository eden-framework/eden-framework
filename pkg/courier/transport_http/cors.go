package transport_http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/profzone/eden-framework/pkg/courier/httpx"
)

func setCORS(headers *http.Header) {
	headers.Set("Access-Control-Allow-Credentials", "false")
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Access-Control-Allow-Methods", strings.Join([]string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodHead,
		http.MethodDelete,
		http.MethodPatch,
	}, ","))
	headers.Set("Access-Control-Allow-Headers", strings.Join([]string{
		"Origin",
		httpx.HeaderContentType,
		"Content-Length",
		"Authorization",
		"AppToken",
		"AccessKey",
	}, ","))
	headers.Set("Access-Control-Max-Age", strconv.FormatInt(int64(12*time.Hour/time.Second), 10))
}
