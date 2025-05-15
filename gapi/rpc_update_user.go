package gapi

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/noueii/gonuxt-starter/db/out"
	"github.com/noueii/gonuxt-starter/pb"
	"github.com/noueii/gonuxt-starter/util"
	"github.com/noueii/gonuxt-starter/validator"
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

	if authPayload.Role != util.AdminRole && authPayload.Username != req.Username {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's data")
	}

	arg := db.UpdateUserByNameParams{
		Name: req.GetUsername(),
		Balance: sql.NullInt32{
			Int32: req.GetBalance(),
			Valid: req.Balance != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to has password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

	}

	user, err := server.db.UpdateUserByName(ctx, arg)
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
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

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
