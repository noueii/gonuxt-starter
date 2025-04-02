package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/noueii/gonuxt-starter/db/out"
)

type Server struct {
	db     *db.Queries
	router *gin.Engine
}

func NewServer(db *db.Queries) *Server {
	server := &Server{db: db}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getAccountById)
	server.router = router

	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
