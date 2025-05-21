package gapi

import (
	"context"
	"net/http"

	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
)

func (server *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if w, ok := GetResponseWriter(ctx); ok {
		secure := false
		sameSite := http.SameSiteLaxMode

		if server.config.Environment == util.Production {

			sameSite = http.SameSiteLaxMode
			secure = true
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/v1/token/refresh",
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
		})

	}

	return &pb.LogoutResponse{}, nil

}
