package client

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/profzone/eden-framework/pkg/courier/transport_http"
)

func NewInsecureTLSTransport(rootCA []byte) transport_http.TransportWrapper {
	return func(rt http.RoundTripper) http.RoundTripper {
		if httpRt, ok := rt.(*http.Transport); ok {
			if httpRt.TLSClientConfig == nil {
				httpRt.TLSClientConfig = &tls.Config{}
			}
			httpRt.TLSClientConfig.RootCAs = rootCertPool(rootCA)
			return httpRt
		}
		return rt
	}
}

func rootCertPool(caData []byte) *x509.CertPool {
	if len(caData) == 0 {
		return nil
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return certPool
}
