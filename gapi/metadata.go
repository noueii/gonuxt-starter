package gapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent    string
	ClientIP     string
	AccessToken  string
	RefreshToken string
}

type httpContextKey string

const (
	grpcGatewayUserAgentHeader                = "grpcgateway-user-agent"
	userAgentHeader                           = "user-agent"
	xForwardedForHeader                       = "x-forwarded-for"
	ResponseWriterKey          httpContextKey = "http-response-writer"
	cookie                                    = "grpcgateway-cookie"
)

const (
	cookieAccessToken  = "access_token"
	cookieRefreshToken = "refresh_token"
)

func (server *Server) getMetadata(ctx context.Context) *Metadata {
	md := &Metadata{}

	if metadataFromContext, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := metadataFromContext.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}

		if userAgents := metadataFromContext.Get(userAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}

		if clientIPs := metadataFromContext.Get(xForwardedForHeader); len(clientIPs) > 0 {
			md.ClientIP = clientIPs[0]
		}

		if cookies := metadataFromContext.Get(cookie); len(cookies) > 0 {
			fmt.Println(cookies)
			parsed := parseCookieHeader(cookies[0])

			md.AccessToken = "bearer " + parsed[cookieAccessToken]
			md.RefreshToken = "bearer " + parsed[cookieRefreshToken]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		md.ClientIP = p.Addr.String()
	}

	return md

}

func GetResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	rw, ok := ctx.Value(ResponseWriterKey).(http.ResponseWriter)
	return rw, ok
}

func parseCookieHeader(header string) map[string]string {
	result := make(map[string]string)

	cookiePairs := strings.Split(header, ";")
	fmt.Println(cookiePairs)
	for _, pair := range cookiePairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}
