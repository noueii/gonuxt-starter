package gapi

import (
	"context"

	"github.com/noueii/gonuxt-starter/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthorized: %s", err)
	}

	newToken, payload, err := server.tokenMaker.CreateToken(payload.Username, payload.Role, server.config.TokenAccessDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error refreshing token: %s", err)
	}

	resp := &pb.RefreshTokenResponse{
		AccessToken:          newToken,
		AccessTokenExpiresAt: timestamppb.New(payload.ExpiresAt),
	}
	return resp, nil
}
