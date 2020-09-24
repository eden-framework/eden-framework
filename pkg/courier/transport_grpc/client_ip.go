package transport_grpc

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
)

func ClientIP(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values, ok := md[httpx.HeaderForwardedFor]; ok {
			clientIP := httpx.GetClientIPByHeaderForwardedFor(strings.Join(values, ""))
			if clientIP != "" {
				return clientIP
			}
		}

		if values, ok := md[httpx.HeaderRealIP]; ok {
			clientIP := httpx.GetClientIPByHeaderRealIP(strings.Join(values, ""))
			if clientIP != "" {
				return clientIP
			}
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		if clientIP, _, err := net.SplitHostPort(strings.TrimSpace(p.Addr.String())); err == nil {
			return clientIP
		}
	}

	return ""
}
