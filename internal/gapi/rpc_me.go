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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Me(ctx context.Context, req *pb.MeRequest) (*pb.MeResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		grpcMeta := server.getMetadata(ctx)
		token := grpcMeta.RefreshToken
		fmt.Println(grpcMeta)
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

	resp := &pb.MeResponse{
		Id:        session.ID.String(),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
		User: &pb.User{
			Id:        session.UserID.String(),
			Email:     session.Email,
			Username:  session.Username,
			Role:      session.Role,
			CreatedAt: timestamppb.New(session.CreatedAt),
		},
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
			Name:     "session",
			Value:    resp.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
			Expires:  payload.ExpiresAt,
		})

	}

	return resp, nil

}
