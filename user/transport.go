package user

import (
	"context"
	"encoding/json"
	"invoice_service/model"
	"invoice_service/security"

	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt/v4"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(_ context.Context, logger log.Logger, r *mux.Router, endpoint Endpoints) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(gokitjwt.HTTPToContext()),
	}

	key := []byte("mysecret")
	keys := func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}

	createUserHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.CreateUser),
		decodeCreateUserRequest,
		encodeResponse,
		options...,
	)

	listUsersHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.ListUsers),
		decodeListUsersReq,
		encodeResponse,
		options...,
	)

	loginHandler := httptransport.NewServer(
		endpoint.LoginHandler,
		decodeLoginReq,
		encodeResponse,
		options...,
	)

	jwtTokenHandler := httptransport.NewServer(
		endpoint.GenerateJWTToken,
		decodeGenerateTokenReq,
		encodeResponse,
		options...,
	)

	r.Methods(http.MethodPost).Path("/create_user").Handler(createUserHandler)
	r.Methods(http.MethodPost).Path("/create_user").Handler(createUserHandler)
	r.Methods(http.MethodGet).Path("/users").Handler(listUsersHandler)
	r.Methods(http.MethodPost).Path("/login").Handler(loginHandler)
	r.Methods(http.MethodPost).Path("/generate_token").Handler(jwtTokenHandler)

	return r
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err == security.InvalidLoginErr || err == security.NotAuthorizedErr {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeCreateUserRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request model.CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeListUsersReq(ctx context.Context, req *http.Request) (interface{}, error) {
	return "", nil
}

func decodeLoginReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request model.LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func decodeGenerateTokenReq(ctx context.Context, req *http.Request) (interface{}, error) {
	return nil, nil
}
