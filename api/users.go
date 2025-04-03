package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/util"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required,alphanum,min=1"`
	Password string `json:"password" binding:"required,min=8"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	name := req.Name
	password := req.Password
	hashedPassword, err := util.HashPassword(password)

	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := s.db.CreateUser(ctx, db.CreateUserParams{
		Name:           name,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}

type getAccountByIdRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountById(ctx *gin.Context) {
	var req getAccountByIdRequest

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
