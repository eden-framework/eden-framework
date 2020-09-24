package transport_http

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/eden-framework/eden-framework/pkg/courier"
	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
	"github.com/eden-framework/eden-framework/pkg/courier/transport_http/transform"
)

func setContextType(w http.ResponseWriter, s string) {
	w.Header().Set(httpx.HeaderContentType, s+";charset=utf-8")
}

func encodeHttpResponse(_ context.Context, w http.ResponseWriter, r *http.Request, response interface{}) error {
	if redirect, ok := response.(httpx.IRedirect); ok {
		http.Redirect(w, r, redirect.Redirect(r.Host), redirect.Status())
		return nil
	}

	if hasMeta, ok := response.(courier.IMeta); ok {
		meta := hasMeta.Meta()
		for k := range meta {
			w.Header().Set(k, meta.Get(k))
		}
	}

	code := http.StatusOK
	if r.Method == http.MethodPost {
		code = http.StatusCreated
	}

	if sc, ok := response.(IStatus); ok {
		code = sc.Status()
	}

	if response == nil {
		code = http.StatusNoContent
	}

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	tpe := reflect.TypeOf(response)
	for tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}

	contentType := ""

	if withContentType, ok := response.(IContentType); ok {
		contentType = withContentType.ContentType()
		setContextType(w, contentType)
	}

	if reader, ok := response.(io.Reader); ok {
		b, errForRead := ioutil.ReadAll(reader)
		if errForRead != nil {
			return errForRead
		}
		w.WriteHeader(code)
		w.Write(b)
		return nil
	}

	if contentType == "" {
		contentType = httpx.MIME_JSON
		setContextType(w, contentType)
	}

	dataList, err := transform.ContentMarshal(reflect.ValueOf(response), "body", contentType)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	w.Write(dataList[0])
	return nil
}

func encodeHttpError(c context.Context, w http.ResponseWriter, r *http.Request, err error) (finalErr error) {
	if redirect, ok := err.(httpx.IRedirect); ok {
		http.Redirect(w, r, redirect.Redirect(r.Host), redirect.Status())
		return nil
	}

	finalStatusErr := status_error.FromError(err).WithSource(c.Value(ContextKeyServerName).(string))
	encodeHttpResponse(c, w, r, finalStatusErr)
	return finalStatusErr
}
