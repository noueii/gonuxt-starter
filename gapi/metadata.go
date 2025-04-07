package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
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
	}

	if p, ok := peer.FromContext(ctx); ok {
		md.ClientIP = p.Addr.String()
	}

	return md

}
