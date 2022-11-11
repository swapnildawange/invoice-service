package user

import (
	"context"
	"encoding/json"
	"fmt"
	"invoice_service/model"
	"invoice_service/security"
	"strconv"

	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"

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

	keys := func(token *jwt.Token) (interface{}, error) {
		key := viper.GetString("JWTSECRET")
		return []byte(key), nil
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

	editUserHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.EditUser),
		decodeEditUserReq,
		encodeResponse,
		options...,
	)

	deleteUserHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.DeleteUser),
		decodeDeleteReq,
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

	r.Methods(http.MethodPost).Path(CreateUserRequestPath).Handler(createUserHandler)
	r.Methods(http.MethodGet).Path("/users").Handler(listUsersHandler)
	r.Methods(http.MethodPatch).Path("/user/{id}").Handler(editUserHandler)
	r.Methods(http.MethodDelete).Path("/user/{id}").Handler(deleteUserHandler)

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
		return nil, fmt.Errorf("failed to decode create user request %v", err)
	}
	return request, nil
}

func decodeListUsersReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		request = model.UserFilter{
			Page:      1,
			SortBy:    "id",
			SortOrder: "ASC",
		}
		err error
	)

	id := req.URL.Query().Get("id")
	if id != "" {
		request.Id, err = strconv.Atoi(id)
		if err != nil || request.Id <= 0 {
			return nil, fmt.Errorf("invalid user id %v", err)
		}

	}

	page := req.URL.Query().Get("page")
	if page != "" {
		request.Page, err = strconv.Atoi(page)
		if err != nil {
			return nil, fmt.Errorf("invalid page value %v", err)
		}
		if request.Page <= 0 {
			request.Page = 1
		}
	}

	firstName := req.URL.Query().Get("first_name")
	if firstName != "" {
		request.FirstName = firstName
	}

	lastName := req.URL.Query().Get("last_name")
	if lastName != "" {
		request.LastName = lastName
	}

	sortBy := req.URL.Query().Get("sort_by")
	if sortBy != "" {
		request.SortBy = sortBy
	}

	sortOrder := req.URL.Query().Get("sort_order")
	if sortOrder != "" {
		request.SortOrder = sortOrder
	}

	return request, nil
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

func decodeDeleteReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		deleteUserReq = model.DeleteUserReq{
			Id: -1,
		}
		err error
	)

	id := req.URL.Query().Get("id")
	if id != "" {
		deleteUserReq.Id, err = strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("Invalid user id")
		}
	}
	deleteUserReq.Email = req.URL.Query().Get("email")

	return deleteUserReq, nil
}

func decodeEditUserReq(ctx context.Context, req *http.Request) (interface{}, error) {

	var request model.EditUserRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}


