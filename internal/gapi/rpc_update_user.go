package gapi

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "github.com/noueii/gonuxt-starter/internal/db/out"
	"github.com/noueii/gonuxt-starter/internal/pb"
	"github.com/noueii/gonuxt-starter/internal/util"
	"github.com/noueii/gonuxt-starter/internal/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx, []string{util.UserRole, util.AdminRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.Role != util.AdminRole && authPayload.UserID.String() != req.GetId() {

		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's data")
	}

	userUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "can't parse uuid")
	}

	arg := db.UpdateUserByIdParams{
		ID: userUUID,

		Balance: sql.NullInt32{
			Int32: req.GetBalance(),
			Valid: req.Balance != nil,
		},
	}

	if req.GetPassword() != "" {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	if req.GetUsername() != "" {
		arg.Name = sql.NullString{
			String: req.GetUsername(),
			Valid:  true,
		}
	}

	if req.GetRole() != "" {
		if authPayload.Role != util.AdminRole {
			return nil, status.Errorf(codes.PermissionDenied, "insufficient permissions to change user role")
		}

		arg.Role = sql.NullString{
			String: req.GetRole(),
			Valid:  true,
		}
	}

	user, err := server.db.UpdateUserById(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "could not update user")
	}

	resp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if req.Password != nil {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
