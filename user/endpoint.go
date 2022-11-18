package user

import (
	"context"
	"invoice_service/security"
	"invoice_service/spec"
	"invoice_service/svcerror"
	"time"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
)

type Endpoints struct {
	CreateUser       endpoint.Endpoint
	ListUsers        endpoint.Endpoint
	GenerateJWTToken endpoint.Endpoint
	LoginHandler     endpoint.Endpoint
	DeleteUser       endpoint.Endpoint
	EditUser         endpoint.Endpoint
}

func NewEndpoints(logger log.Logger, bl BL) Endpoints {
	return Endpoints{
		CreateUser:       makeCreateUser(logger, bl),
		ListUsers:        makeListUsers(logger, bl),
		GenerateJWTToken: makeGenerateJWT(logger, bl),
		LoginHandler:     makeLoginHandler(logger, bl),
		DeleteUser:       makeDeleteUser(logger, bl),
		EditUser:         makeEditUser(logger, bl),
	}
}

func makeLoginHandler(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		loginReq, ok := request.(spec.LoginRequest)
		if !ok {
			logger.Log("[debug]", svcerror.ErrInvalidRequest)
			return nil, svcerror.ErrInvalidRequest
		}
		user, token, err := bl.Login(ctx, loginReq)
		if err != nil {
			return user, svcerror.ErrLoginFailed
		}

		return spec.LoginResponse{
			User:  user,
			Token: token,
		}, nil
	}
}

func makeGenerateJWT(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token, err := security.GenerateJWT(viper.GetString("JWTSECRET"), 1, 1)
		if err != nil {
			return "", svcerror.ErrFailedToGenerateJWT
		}
		return token, nil
	}
}

func makeCreateUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "createUser",
				"took", time.Since(begin),
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
			return nil, svcerror.ErrNotAuthorized
		}

		req = request.(spec.CreateUserRequest)
		if err != nil {
			logger.Log("[debug]", "invalid request for create user", "err", err)
			return nil, svcerror.ErrInvalidRequest
		}
		user, err = bl.CreateUser(ctx, req)
		if err != nil {
			logger.Log("[debug]", err)
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
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			return nil, svcerror.ErrNotAuthorized
		}

		req := request.(spec.UserFilter)

		response, err = bl.ListUsers(ctx, req)
		if err != nil {
			logger.Log("[debug]", err)
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
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			return nil, svcerror.ErrNotAuthorized
		}

		editUserReq := request.(spec.EditUserRequest)
		user, err := bl.EditUser(ctx, editUserReq)
		if err != nil {
			logger.Log("[debug]", err)
			return nil, svcerror.ErrFailedToUpdateUser
		}
		return user, nil
	}
}

func makeDeleteUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "deleteUser",
				"took", time.Since(begin),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			return nil, svcerror.ErrNotAuthorized
		}

		deleteUserReq := request.(spec.DeleteUserReq)
		response, err = bl.DeleteUser(ctx, deleteUserReq)
		if err != nil {
			logger.Log("[debug]", err)
			return nil, svcerror.ErrFailedToDeleteUser
		}
		return response, nil
	}
}
