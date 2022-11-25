package invoice

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
	"github.com/invoice-service/svcerror"
	"github.com/invoice-service/utils"

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

	createInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.CreateInvoice),
		decodeCreateInvoiceReq,
		encodeResponse,
		options...,
	)

	getInvoiceHandler := httptransport.NewServer(
		gokitjwt.NewParser(keys, jwt.SigningMethodHS256, security.GetJWTClaims)(endpoint.GetInvoice),
		decodeGetInvoiceReq,
		encodeResponse,
		options...,
	)

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
	r.Methods(http.MethodGet).Path(GetInvoiceRequestPath).Handler(getInvoiceHandler)
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

	if err == svcerror.ErrInvalidLoginCreds || err == svcerror.ErrNotAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
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

func validateCreateInvoiceRequest(request spec.CreateInvoiceRequest) error {
	err := validate.Struct(request)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New(e.Translate(trans))
		}
	}
	return nil
}

func decodeCreateInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request spec.CreateInvoiceRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	if err := validateCreateInvoiceRequest(request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request spec.GetInvoiceRequest

	query := req.URL.Query()
	invoice_id := query.Get("id")
	request.Id = invoice_id

	return request, nil
}

func decodeListInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var invoiceFilter = spec.InvoiceFilter{
		Paid:      -1,
		SortBy:    "id",
		SortOrder: "ASC",
		Page:      1,
	}

	invoiceId := req.URL.Query().Get("id")
	if invoiceId != "" {
		invoiceFilter.Id = invoiceId
	}

	page := req.URL.Query().Get("page")
	if page != "" {
		page, err := strconv.Atoi(page)
		if err != nil {
			return nil, svcerror.ErrBadRouting
		}
		invoiceFilter.Page = page
	}

	userId := req.URL.Query().Get("user_id")
	if userId != "" {
		userId, err := strconv.Atoi(userId)
		if err != nil {
			return nil, svcerror.ErrBadRouting
		}
		invoiceFilter.UserId = userId
	}

	adminId := req.URL.Query().Get("admin_id")
	if adminId != "" {
		adminId, err := strconv.Atoi(adminId)
		if err != nil {
			return nil, svcerror.ErrBadRouting
		}
		invoiceFilter.AdminId = adminId
	}

	paid := req.URL.Query().Get("paid")
	if paid != "" {
		paid, err := strconv.ParseFloat(paid, 64)
		if err != nil {
			return nil, svcerror.ErrBadRouting
		}
		invoiceFilter.Paid = paid
	}

	paymentStatus := req.URL.Query().Get("payment_status")
	if paymentStatus != "" {
		paymentStatus, err := strconv.Atoi(paymentStatus)
		if err != nil {
			return nil, svcerror.ErrBadRouting
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
		return nil, svcerror.ErrBadRouting
	}

	return invoiceId, nil
}

func decodeUpdateInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request = spec.UpdateInvoiceRequest{
		Paid:          -1,
		PaymentStatus: -1,
	}
	invoiceId, ok := mux.Vars(req)["id"]
	if !ok {
		return nil, svcerror.ErrBadRouting
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, svcerror.ErrBadRouting
	}

	request.Id = invoiceId
	return request, nil
}
