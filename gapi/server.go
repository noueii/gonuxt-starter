package gapi

import (
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/pb"
	"github.com/noueii/gonuxt-starter/token"
	"github.com/noueii/gonuxt-starter/util"
)

type Server struct {
	pb.UnimplementedGoNuxtServer
	config     *util.Config
	db         *db.Queries
	tokenMaker token.Maker
}

func NewServer(config *util.Config, db *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		db:         db,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
