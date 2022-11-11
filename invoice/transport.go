package invoice

import (
	"context"
	"encoding/json"
	"fmt"
	"invoice_service/model"
	"invoice_service/security"
	"strconv"

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
		key := viper.GetString("JWTSECRET")
		return []byte(key), nil
	}

	createInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.CreateInvoice),
		decodeCreateInvoiceReq,
		encodeResponse,
		options...,
	)

	// getInvoiceHandler := httptransport.NewServer(
	// 	endpoint.GetInvoice,
	// 	decodeGetInvoiceReq,
	// 	encodeResponse,
	// 	options...,
	// )

	listInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.ListInvoice),
		decodeListInvoiceReq,
		encodeResponse,
		options...,
	)

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

	r.Methods(http.MethodPost).Path(CreateInvoiceRequestPath).Handler(createInvoiceHandler)
	r.Methods(http.MethodPatch).Path(EditInvoiceRequestPath).Handler(updateInvoiceHandler)
	r.Methods(http.MethodDelete).Path(DeleteInvoiceRequestPath).Handler(deleteInvoiceHandler)
	r.Methods(http.MethodGet).Path(ListInvoiceRequestPath).Handler(listInvoiceHandler)

	return r
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
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

// func decodeGetInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
// 	var request model.GetInvoiceRequest

// 	query := req.URL.Query()
// 	invoice_id := query.Get("id")
// 	request.Id = invoice_id

// 	return request, nil
// }

func decodeListInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var invoiceFilter = model.InvoiceFilter{
		Paid:      -1,
		SortBy:    "id",
		SortOrder: "ASC",
		Page:      1,
	}

	invoiceId := req.URL.Query().Get("id")
	if invoiceId != "" {
		invoiceFilter.Id = invoiceId
	}

	userId := req.URL.Query().Get("user_id")
	if userId != "" {
		userId, err := strconv.Atoi(userId)
		if err != nil {
			return nil, fmt.Errorf("invalid user id provided in query params")
		}
		invoiceFilter.UserId = userId
	}

	adminId := req.URL.Query().Get("admin_id")
	if adminId != "" {
		adminId, err := strconv.Atoi(adminId)
		if err != nil {
			return nil, fmt.Errorf("invalid admin id provided in query params")
		}
		invoiceFilter.AdminId = adminId
	}

	paid := req.URL.Query().Get("paid")
	if paid != "" {
		paid, err := strconv.ParseFloat(paid, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid paid amount provided in query params")
		}
		invoiceFilter.Paid = paid
	}

	paymentStatus := req.URL.Query().Get("payment_status")
	if paymentStatus != "" {
		paymentStatus, err := strconv.Atoi(paymentStatus)
		if err != nil {
			return nil, fmt.Errorf("invalid paymentStatus amount provided in query params")
		}
		invoiceFilter.PaymentStatus = paymentStatus
	}
	sortBy := req.URL.Query().Get("sort_by")
	if sortBy != "" {
		invoiceFilter.SortBy = sortBy
	}

	sortOrder := req.URL.Query().Get("sort_order")
	if sortOrder != "" {
		invoiceFilter.SortOrder = sortOrder
	}

	return invoiceFilter, nil
}

func decodeDeleteInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {

	invoiceId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, fmt.Errorf("invoice id not found in url path")
	}

	return invoiceId, nil
}

func decodeUpdateInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request = model.UpdateInvoiceRequest{
		Paid:          -1,
		PaymentStatus: -1,
	}
	invoiceId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, fmt.Errorf("invoice id not found in url path")
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.Id = invoiceId
	return request, nil
}
