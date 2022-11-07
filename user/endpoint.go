package user

import (
	"context"
	"fmt"
	"invoice_service/model"
	"invoice_service/security"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Endpoints struct {
	CreateUser       endpoint.Endpoint
	ListUsers        endpoint.Endpoint
	GenerateJWTToken endpoint.Endpoint
	LoginHandler     endpoint.Endpoint
}

func NewEndpoints(logger log.Logger, bl BL) Endpoints {
	return Endpoints{
		CreateUser:       makeCreateUser(logger, bl),
		ListUsers:        makeListUsers(logger, bl),
		GenerateJWTToken: makeGenerateJWT(logger, bl),
		LoginHandler:     makeLoginHandler(logger, bl),
	}
}

func makeLoginHandler(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		loginReq, ok := request.(model.LoginRequest)

		if !ok {
			return nil, fmt.Errorf("Invalid login request")
		}
		user, token, err := bl.Login(ctx, loginReq)
		if err != nil {
			return user, err
		}

		return model.LoginResponse{
			User:  user,
			Token: token,
		}, nil
	}
}

func makeGenerateJWT(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token, err := security.GenerateJWT("mysecret", 1, 1)
		if err != nil {
			return "", err
		}
		return token, nil
	}
}

func makeCreateUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req  model.CreateUserRequest
			user model.User
		)
		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid jwt token")
		}
		if JWTClaims.Role == 2 {
			return nil, security.NotAuthorizedErr
		}

		req = request.(model.CreateUserRequest)
		if err != nil {
			return nil, fmt.Errorf("invalid request for create user")
		}
		user, err = bl.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func makeListUsers(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid jwt token")
		}

		if JWTClaims.Role == 2 {
			return nil, security.NotAuthorizedErr
		}

		response, nil := bl.ListUsers(ctx)
		if err != nil {
			logger.Log("endpoint", "makeListUsers", "Failed to list users", err.Error())
			return
		}
		return

	}
}