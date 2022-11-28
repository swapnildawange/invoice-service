package user

import (
	"context"
	"fmt"
	"time"

	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
	"github.com/invoice-service/svcerror"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	CreateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	ListUsers    endpoint.Endpoint
	LoginHandler endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
	EditUser     endpoint.Endpoint
}

func NewEndpoints(logger log.Logger, bl BL) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUser(logger, bl),
		GetUser:    makeGetUser(logger, bl),
		ListUsers:  makeListUsers(logger, bl),
		DeleteUser: makeDeleteUser(logger, bl),
		EditUser:   makeEditUser(logger, bl),
	}
}

func makeCreateUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "createUser",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%#v", request),
				"response", fmt.Sprintf("%#v", response),
			)
		}(time.Now())

		var (
			req  spec.CreateUserRequest
			user spec.User
		)

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}
		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		req = request.(spec.CreateUserRequest)
		user, err = bl.CreateUser(ctx, req)
		if err != nil {
			logger.Log("[debug]", "failed to create user", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return nil, svcerror.ErrFailedToCreateUser
		}
		return user, nil
	}
}

func makeListUsers(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "listUsers",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%#v", request),
				"response", fmt.Sprintf("%#v", response),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			response, err = bl.GetUser(ctx, JWTClaims.Id)
			if err != nil {
				_, ok := err.(*svcerror.CustomErrString)
				if ok {
					return nil, err
				}
				return response, svcerror.ErrFailedToListUsers
			}
			return
		}

		req := request.(spec.UserFilter)
		response, err = bl.ListUsers(ctx, req)
		if err != nil {
			logger.Log("[debug]", "failed to list users", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return response, svcerror.ErrFailedToListUsers
		}
		return
	}
}

func makeEditUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "editUser",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%#v", request),
				"response", fmt.Sprintf("%#v", response),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		editUserReq := request.(spec.EditUserRequest)
		user, err := bl.EditUser(ctx, editUserReq)
		if err != nil {
			logger.Log("[debug]", "failed to edit user", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return response, svcerror.ErrFailedToUpdateUser
		}
		return user, nil
	}
}

func makeDeleteUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var userId int
		defer func(begin time.Time) {
			logger.Log(
				"method", "deleteUser",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%#v", request),
				"response", fmt.Sprintf("%#v", userId),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		deleteUserReq := request.(spec.DeleteUserReq)
		userId, err = bl.DeleteUser(ctx, deleteUserReq)
		if err != nil {
			logger.Log("[debug]", "failed to delete user", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return userId, svcerror.ErrFailedToDeleteUser
		}
		return userId, nil
	}
}

func makeGetUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var user spec.User
		defer func(begin time.Time) {
			logger.Log(
				"method", "deleteUser",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%#v", request),
				"response", fmt.Sprintf("%#v", response),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		userId := request.(int)
		user, err = bl.GetUser(ctx, userId)
		logger.Log("[debug]", "Failed to get user", "err", err)
		if err != nil {
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return response, svcerror.ErrFailedToListUsers
		}
		return user, nil
	}
}
