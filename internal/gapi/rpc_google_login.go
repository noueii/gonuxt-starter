package gapi

import (
	"context"
	"fmt"

	"github.com/noueii/gonuxt-starter/internal/pb"
	"golang.org/x/oauth2"
)

func (server *Server) GoogleLogin(context context.Context, req *pb.GoogleLoginRequest) (*pb.GoogleLoginResponse, error) {
	fmt.Println("I'm here")
	url := server.config.OAuth.Google.AuthCodeURL("some-state", oauth2.AccessTypeOffline)
	return &pb.GoogleLoginResponse{
		RedirectUrl: url,
	}, nil

}
