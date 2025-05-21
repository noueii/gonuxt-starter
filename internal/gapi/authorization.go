package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/noueii/gonuxt-starter/internal/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	fmt.Println(values)
	if len(values) == 0 {
		grpcMeta := server.getMetadata(ctx)
		token := grpcMeta.AccessToken

		values = []string{token}
		fmt.Println(values)
		if len(values) == 0 {
			return nil, fmt.Errorf("missing authorization header")
		}
	}

	authHeader := values[0]

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]

	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	if !hasPermission(payload.Role, accessibleRoles) {
		return nil, fmt.Errorf("permission denied")
	}

	return payload, nil

}

func hasPermission(role string, accessibleRoles []string) bool {
	for _, r := range accessibleRoles {
		if role == r {
			return true
		}
	}

	return false
}
