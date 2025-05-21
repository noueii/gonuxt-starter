package gapi

import (
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Name,
		CreatedAt: timestamppb.New(user.CreatedAt),
		Role:      user.Role,
	}
}
