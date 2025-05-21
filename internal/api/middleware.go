package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/noueii/gonuxt-starter/internal/token"
)

const (
	headerAuthorizationKey        = "authorization"
	headerAuthorizationType       = "bearer"
	headerAuthorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(headerAuthorizationKey)
		if len(authHeader) == 0 {
			err := errors.New("missing authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != headerAuthorizationType {
			err := fmt.Errorf("invalid authorization type")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(headerAuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
