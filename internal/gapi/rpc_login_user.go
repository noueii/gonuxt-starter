package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.db.GetUserByEmail(ctx, req.GetEmail())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	if !user.HashedPassword.Valid {
		return nil, status.Errorf(codes.Internal, "CRITICAL: Missing user password. Contact support. %s", err)
	}

	if err = util.CheckPassword(req.GetPassword(), user.HashedPassword.String); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, user.Name, user.Role, server.config.TokenAccessDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, user.Name, user.Role, server.config.TokenRefreshDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	metadata := server.getMetadata(ctx)

	session, err := server.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	resp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt),
		User:                  convertUser(user),
	}

	if w, ok := GetResponseWriter(ctx); ok {
		fmt.Println("FOUND RESPONSE WRITER")
		secure := false
		sameSite := http.SameSiteLaxMode

		if server.config.Environment == util.Production {

			fmt.Printf("PRODUCTION ENV, %s\n", server.config.Environment)
			sameSite = http.SameSiteLaxMode
			secure = true
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
			Expires:  accessPayload.ExpiresAt,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/v1/token/refresh",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
			Expires:  refreshPayload.ExpiresAt,
		})
	}

	return resp, nil
}
