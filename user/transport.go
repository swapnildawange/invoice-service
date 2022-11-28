package user

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
	"github.com/invoice-service/svcerror"
	"github.com/invoice-service/utils"

	"net/http"

	"github.com/go-kit/log"
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
		key := viper.GetString("ACCESS_TOKEN_SECRET")
		return []byte(key), nil
	}

	createUserHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.CreateUser),
		decodeCreateUserRequest,
		encodeResponse,
		options...,
	)

	getUserHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.GetUser),
		decodeGetUserRequest,
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
		decodeDeleteUserReq,
		encodeResponse,
		options...,
	)

	r.Methods(http.MethodPost).Path(spec.CreateUserRequestPath).Handler(createUserHandler)
	r.Methods(http.MethodGet).Path(spec.GetUserRequestPath).Handler(getUserHandler)
	r.Methods(http.MethodGet).Path(spec.ListUsersRequestPath).Handler(listUsersHandler)
	r.Methods(http.MethodPatch).Path(spec.EditUserRequestPath).Handler(editUserHandler)
	r.Methods(http.MethodDelete).Path(spec.DeleteUserRequestPath).Handler(deleteUserHandler)

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

	if err == svcerror.ErrInvalidLoginCreds || err == svcerror.ErrNotAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		_, ok := err.(*svcerror.CustomErrString)
		if ok {
			if err == svcerror.ErrBadRouting {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate, _ = utils.InitValidator()

	trans = utils.InitTranslator()
}

func validateCreateUserRequest(request spec.CreateUserRequest) error {
	err := validate.Struct(request)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New(e.Translate(trans))
		}
	}
	return nil
}

func decodeCreateUserRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request spec.CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, svcerror.ErrInvalidRequest
	}
	if err := validateCreateUserRequest(request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeListUsersReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		request = spec.UserFilter{
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
			return nil, svcerror.ErrBadRouting
		}
	}

	page := req.URL.Query().Get("page")
	if page != "" {
		request.Page, err = strconv.Atoi(page)
		if err != nil {
			return nil, svcerror.ErrBadRouting

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

func decodeGetUserRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		err error
		id  int
	)

	userId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, svcerror.ErrBadRouting
	} else if userId != "" {
		id, err = strconv.Atoi(userId)
		if err != nil {
			return nil, svcerror.ErrBadRouting
		}
	}
	return id, nil
}

func decodeDeleteUserReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		deleteUserReq = spec.DeleteUserReq{
			Id: -1,
		}
		err error
	)
	id := req.URL.Query().Get("id")
	if id != "" {
		deleteUserReq.Id, err = strconv.Atoi(id)
		if err != nil || deleteUserReq.Id <= 0 {
			return nil, svcerror.ErrBadRouting
		}
	}
	deleteUserReq.Email = req.URL.Query().Get("email")
	if deleteUserReq.Id == -1 && deleteUserReq.Email == "" {
		return nil, svcerror.ErrBadRouting
	}
	return deleteUserReq, nil
}

func decodeEditUserReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		request spec.EditUserRequest
		err     error
	)
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	userId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, svcerror.ErrBadRouting
	} else if userId != "" {
		request.Id, err = strconv.Atoi(userId)
		if err != nil {
			return nil, svcerror.ErrBadRouting
		}
	}

	return request, nil
}
