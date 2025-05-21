package gapi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RefreshToken(ctx context.Context, _ *emptypb.Empty) (*pb.RefreshTokenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		grpcMeta := server.getMetadata(ctx)
		token := grpcMeta.RefreshToken

		values = []string{token}

		if len(values) == 0 {
			return nil, fmt.Errorf("missing authorization header")
		}
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	fmt.Println("Passed verification")

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	refreshToken := fields[1]

	payload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	session, err := server.db.GetSessionById(ctx, payload.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "Session not found: %s", err)
		}

		return nil, status.Errorf(codes.Internal, "Server error: %s", err)
	}

	if session.IsRevoked.Bool {
		return nil, status.Error(codes.Unauthenticated, "Blocked session")
	}

	if session.Email != payload.Email {
		return nil, status.Error(codes.Unauthenticated, "incorrect session user")
	}

	if session.RefreshToken != refreshToken {
		return nil, status.Error(codes.Unauthenticated, "session token missmatch")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, status.Error(codes.Unauthenticated, "session expired")
	}

	newToken, payload, err := server.tokenMaker.CreateToken(payload.UserID, payload.Email, session.Username, session.Role, server.config.TokenAccessDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error refreshing token: %s", err)
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
			Value:    newToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
			Expires:  payload.ExpiresAt,
		})

	}

	resp := &pb.RefreshTokenResponse{
		AccessToken:          newToken,
		AccessTokenExpiresAt: timestamppb.New(payload.ExpiresAt),
		Session: &pb.Session{
			Id:        session.ID.String(),
			ExpiresAt: timestamppb.New(session.ExpiresAt),
		},
		User: &pb.User{
			Id:       session.UserID.String(),
			Email:    session.Email,
			Username: session.Username,
			Role:     session.Role,
		},
	}
	return resp, nil
}
