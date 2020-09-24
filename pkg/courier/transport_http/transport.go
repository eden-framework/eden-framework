package transport_http

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	logContext "github.com/eden-framework/eden-framework/pkg/context"
	"github.com/eden-framework/eden-framework/pkg/courier"
	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
	"github.com/eden-framework/eden-framework/pkg/duration"
)

func CreateHttpHandler(s *ServeHTTP, ops ...courier.IOperator) httprouter.Handle {
	operatorMetas := courier.ToOperatorMetaList(ops...)

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var err error
		ctx := r.Context()
		ctx = ContextWithServiceName(ctx, s.Name)
		ctx = ContextWithOperators(ctx, ops...)
		ctx = ContextWithRequest(ctx, r)

		codeWriter := &ResponseWriter{
			ResponseWriter: w,
			Code:           0,
			written:        -1,
		}

		w.Header().Set("X-Reversion", ProjectRef)

		reqID := r.Header.Get(httpx.HeaderRequestID)

		if reqID == "" {
			reqID = uuid.New().String()
		}

		logContext.SetLogID(reqID)
		defer logContext.Close()

		d := duration.NewDuration()

		defer func() {
			if r.URL.Path == "/healthz" {
				return
			}
			fields := logrus.Fields{
				"tag":        "access",
				"log_id":     reqID,
				"remote_ip":  GetClientIP(r),
				"method":     r.Method,
				"pathname":   r.URL.Path,
				"user_agent": r.Header.Get(httpx.HeaderUserAgent),
			}

			fields["status"] = codeWriter.Code
			fields["request_time"] = d.Get()

			logger := logrus.WithFields(fields)

			if err != nil {
				if codeWriter.Code >= http.StatusInternalServerError {
					logger.Errorf(err.Error())
				} else {
					logger.Warnf(err.Error())
				}
			} else {
				logger.Infof("")
			}
		}()

		opDecode := createHttpRequestDecoder(r, &params)

		for _, opMeta := range operatorMetas {
			op, decodeErr := courier.NewOperatorBy(opMeta.Type, opMeta.Operator, opDecode)
			if decodeErr != nil {
				err = encodeHttpError(ctx, codeWriter, r, decodeErr)
				return
			}

			response, endpointErr := op.Output(ctx)
			if canCookie, ok := op.(httpx.ICookie); ok {
				cookie := canCookie.Cookies()
				if cookie != nil {
					http.SetCookie(w, cookie)
				}
			}

			if endpointErr != nil {
				err = encodeHttpError(ctx, codeWriter, r, endpointErr)
				return
			}

			if !opMeta.IsLast {
				// set result in context with key of operator name
				ctx = context.WithValue(ctx, opMeta.ContextKey, response)
				continue
			}

			if ws, ok := response.(IWebSocket); ok {
				conn, errForUpgrade := (&websocket.Upgrader{}).Upgrade(codeWriter, r, nil)
				if errForUpgrade != nil {
					err = errForUpgrade
					return
				}
				ws.SubscribeOn(conn)
				return
			}

			encodeErr := encodeHttpResponse(ctx, codeWriter, r, response)
			if encodeErr != nil {
				err = encodeHttpError(ctx, codeWriter, r, encodeErr)
				return
			}
		}
	}
}

var ProjectRef = os.Getenv("PROJECT_REF")

var (
	ContextKeyServerName = uuid.New().String()
	ContextKeyRequest    = uuid.New().String()
	ContextKeyOperators  = uuid.New().String()
)

func ContextWithServiceName(ctx context.Context, serverName string) context.Context {
	return context.WithValue(ctx, ContextKeyServerName, serverName)
}

func ContextWithOperators(ctx context.Context, ops ...courier.IOperator) context.Context {
	return context.WithValue(ctx, ContextKeyOperators, ops)
}

func GetOperators(ctx context.Context) []courier.IOperator {
	return ctx.Value(ContextKeyOperators).([]courier.IOperator)
}

func ContextWithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, ContextKeyRequest, req)
}

func GetRequest(ctx context.Context) *http.Request {
	return ctx.Value(ContextKeyRequest).(*http.Request)
}

type ResponseWriter struct {
	http.ResponseWriter
	Code    int
	written int64
}

func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.Code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	w.written += int64(n)
	return n, err
}
