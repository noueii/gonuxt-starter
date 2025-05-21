package gapi

import (
	"context"

	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) VerifyToken(ctx context.Context, _ *emptypb.Empty) (*pb.VerifyTokenResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.UserRole, util.AdminRole})

	if err != nil {
		return nil, unauthenticatedError(err)
	}

	resp := &pb.VerifyTokenResponse{
		UserId: authPayload.ID.String(),
		Role:   authPayload.Role,
	}

	return resp, nil

}
