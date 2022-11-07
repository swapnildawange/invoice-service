package invoice

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

	createInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.CreateInvoice),
		decodeCreateInvoiceReq,
		encodeResponse,
		options...,
	)

	getInvoiceHandler := httptransport.NewServer(
		endpoint.GetInvoice,
		decodeGetInvoiceReq,
		encodeResponse,
		options...,
	)

	// listInvoiceHandler := httptransport.NewServer(
	// 	endpoint.ListInvoice,
	// 	decodeListInvoiceReq,
	// 	encodeResponse,
	// 	options...,
	// )

	updateInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.UpdateInvoice),
		decodeUpdateInvoiceReq,
		encodeResponse,
		options...,
	)

	deleteInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.DeleteInvoice),
		decodeDeleteInvoiceReq,
		encodeResponse,
		options...,
	)

	

	// admin routes
	r.Methods(http.MethodPost).Path("/create_invoice").Handler(createInvoiceHandler)
	r.Methods(http.MethodPatch).Path("/update_invoice/{id}").Handler(updateInvoiceHandler)
	r.Methods(http.MethodDelete).Path("/invoice/{id}").Handler(deleteInvoiceHandler)

	// user routes

	// common routes
	r.Methods(http.MethodGet).Path("/invoice/{id}").Handler(getInvoiceHandler)
	// r.Methods(http.MethodGet).Path("/invoice").Handler(listInvoiceHandler) //need to change


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


func decodeCreateInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request model.CreateInvoiceRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeGetInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request model.GetInvoiceRequest

	query := req.URL.Query()
	invoice_id := query.Get("id")
	request.Id = invoice_id

	return request, nil
}

func decodeListInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	id, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, fmt.Errorf("user id not found in url path")
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func decodeDeleteInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {

	invoiceId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, fmt.Errorf("Invoice id not found in url path")
	}

	return invoiceId, nil
}

func decodeUpdateInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request model.UpdateInvoiceRequest
	invoiceId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, fmt.Errorf("Invoice id not found in url path")
	}

	if err := json.NewDecoder(req.Body).Decode(&request.Invoice); err != nil {
		return nil, err
	}

	request.Invoice.Id = invoiceId
	return request, nil
}

func decodeCreateUserRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request model.CreateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

