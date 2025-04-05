package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/token"
	"github.com/noueii/gonuxt-starter/util"
)

type Server struct {
	config     *util.Config
	db         *db.Queries
	router     *gin.Engine
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
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(tokenMaker))

	authRoutes.GET("/users/:id", server.getUserById)
	server.router = router

	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
