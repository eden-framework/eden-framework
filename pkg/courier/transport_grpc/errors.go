package transport_grpc

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	httpCodes = map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           http.StatusRequestTimeout,
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusRequestTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.Unauthenticated:    http.StatusUnauthorized,
		codes.ResourceExhausted:  http.StatusForbidden,
		codes.FailedPrecondition: http.StatusPreconditionFailed,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
	}
	statusCodes = map[int]codes.Code{
		http.StatusOK:                  codes.OK,
		http.StatusRequestTimeout:      codes.DeadlineExceeded,
		http.StatusInternalServerError: codes.Internal,
		http.StatusBadRequest:          codes.InvalidArgument,
		http.StatusNotFound:            codes.NotFound,
		http.StatusConflict:            codes.AlreadyExists,
		http.StatusForbidden:           codes.PermissionDenied,
		http.StatusUnauthorized:        codes.Unauthenticated,
		http.StatusPreconditionFailed:  codes.FailedPrecondition,
		http.StatusNotImplemented:      codes.Unimplemented,
		http.StatusServiceUnavailable:  codes.Unavailable,
	}
)

func HTTPStatusFromCode(code codes.Code) int {
	if status, ok := httpCodes[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func CodeFromHTTPStatus(status int) codes.Code {
	if code, ok := statusCodes[status]; ok {
		return code
	}
	return codes.Unknown
}
