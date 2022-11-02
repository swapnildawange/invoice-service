package invoice

import (
	"context"
	"encoding/json"
	"fmt"
	"invoice_service/model"

	"net/http"

	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(logger log.Logger, endpoint Endpoints) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	createInvoiceHandler := httptransport.NewServer(
		endpoint.CreateInvoice,
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

	updateInvoiceHandler := httptransport.NewServer(
		endpoint.UpdateInvoice,
		decodeUpdateInvoiceReq,
		encodeResponse,
		options...,
	)

	deleteInvoiceHandler := httptransport.NewServer(
		endpoint.DeleteInvoice,
		decodeDeleteInvoiceReq,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter()
	r.Methods(http.MethodPost).Path("/create_invoice").Handler(createInvoiceHandler)
	r.Methods(http.MethodGet).Path("/invoice").Handler(getInvoiceHandler)
	r.Methods(http.MethodPatch).Path("/update_invoice/{id}").Handler(updateInvoiceHandler)
	r.Methods(http.MethodDelete).Path("/invoice/{id}").Handler(deleteInvoiceHandler)

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	if err != nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusInternalServerError)
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
