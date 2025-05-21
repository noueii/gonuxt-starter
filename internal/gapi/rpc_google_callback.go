package gapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GoogleCallback(ctx context.Context, req *pb.GoogleCallbackRequest) (*pb.GoogleCallbackResponse, error) {
	token, err := server.config.OAuth.Google.Exchange(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "exchange failed: %s", err.Error())
	}

	userInfo, err := getUserInfo(ctx, &server.config.OAuth.Google, token)

	fmt.Println(userInfo)

	fmt.Println("Passed?")

	if err != nil {
		if w, ok := GetResponseWriter(ctx); ok {
			http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !userInfo.VerifiedEmail {
		if w, ok := GetResponseWriter(ctx); ok {
			http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

		}
		return nil, status.Error(codes.PermissionDenied, "Email must be verified")
	}

	user, err := server.db.GetUserByEmail(ctx, userInfo.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			args := db.CreateUserParams{
				Email:         userInfo.Email,
				EmailVerified: true,
				Name:          userInfo.Name,
			}

			user, err = server.db.CreateUser(ctx, args)
			if err != nil {
				if w, ok := GetResponseWriter(ctx); ok {
					http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

				}
				return nil, status.Errorf(codes.Internal, "Error creating user: %s", err)
			}
		} else {
			if w, ok := GetResponseWriter(ctx); ok {
				http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

			}
			return nil, status.Error(codes.Internal, err.Error())

		}

	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, user.Name, user.Role, server.config.TokenAccessDuration)

	if err != nil {
		if w, ok := GetResponseWriter(ctx); ok {
			http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

		}
		return nil, status.Errorf(codes.Internal, "Error creating access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, user.Name, user.Role, server.config.TokenRefreshDuration)

	if err != nil {
		if w, ok := GetResponseWriter(ctx); ok {
			http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

		}
		return nil, status.Errorf(codes.Internal, "Error creating refresh token: %s", err)
	}

	metadata := server.getMetadata(ctx)

	_, err = server.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		if w, ok := GetResponseWriter(ctx); ok {
			http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

		}
		return nil, status.Errorf(codes.Internal, "Error creating refresh token: %s", err)
	}

	if w, ok := GetResponseWriter(ctx); ok {
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

		http.Redirect(w, &http.Request{}, "http://localhost:3000/auth/redirect", http.StatusTemporaryRedirect)

		fmt.Println("Actually setting cookies")

	}

	fmt.Println("Finished setting cookies")
	return &pb.GoogleCallbackResponse{
		JwtToken: accessToken,
	}, nil

}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func getUserInfo(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (*GoogleUser, error) {
	client := config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var user GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
