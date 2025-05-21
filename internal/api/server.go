package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/token"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
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

func (server *Server) Shutdown() error {
	return server.Shutdown()
}

func RunGinServer(ctx context.Context, waitGroup *errgroup.Group, config *util.Config, queries *db.Queries) {
	httpServer, err := NewServer(config, queries)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create http server:")
	}

	waitGroup.Go(func() error {
		fmt.Println("Starting server")
		err = httpServer.Start(config.HTTPAddr)

		if err != nil {
			log.Fatal().Err(err).Msg("cannot start HTTP server")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shutting down HTTP server")

		err = httpServer.Shutdown()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			log.Error().Err(err).Msg("could not shutdown HTTP server")
			return err
		}

		return nil
	})

}
