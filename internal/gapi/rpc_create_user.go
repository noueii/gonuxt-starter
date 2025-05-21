package gapi

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/noueii/gonuxt-starter/internal/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)

	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	email := req.GetEmail()
	sections := strings.Split(email, "@")
	if len(sections) == 0 {
		return nil, status.Errorf(codes.Internal, "bad email: %s", err)
	}
	username := sections[0]

	arg := db.CreateUserParams{
		Email: req.GetEmail(),
		Name:  username,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  hashedPassword != "",
		},
	}

	user, err := server.db.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "could not create user")
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return resp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
