package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/noueii/gonuxt-starter/internal/validator"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	email := req.Email
	if err := validator.ValidateEmail(email); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	password := req.Password
	hashedPassword, err := util.HashPassword(password)

	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := s.db.CreateUser(ctx, db.CreateUserParams{
		Email: email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}

type getUserByIdRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserById(ctx *gin.Context) {
	var req getUserByIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userUUID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.db.GetUserById(ctx, userUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func getUserResponseFromUser(user db.User) *userResponse {
	return &userResponse{
		ID:        user.ID,
		Username:  user.Name,
		CreatedAt: user.CreatedAt,
	}
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,alphanum,min=6"`
	Password string `json:"password" binding:"required,alphanum,min=8"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := server.db.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if err = util.CheckPassword(req.Password, user.HashedPassword.String); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, user.Name, user.Role, server.config.TokenAccessDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, user.Email, user.Name, user.Role, server.config.TokenRefreshDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User: userResponse{
			ID:        user.ID,
			Username:  user.Name,
			CreatedAt: user.CreatedAt,
		},
	}

	ctx.JSON(http.StatusOK, resp)
}
