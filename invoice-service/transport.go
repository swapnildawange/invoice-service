package invoice

import (
	"context"
	"encoding/json"
	"invoicing/invoice-service/models"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(logger log.Logger, bl BL) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	createInvoiceHandler := httptransport.NewServer(
		makeCreateInvoice(logger, bl),
		decodeCreateInvoiceReq,
		encodeResponse,
		options...,
	)

	getInvoiceHandler := httptransport.NewServer(makeGetInvoice(logger, bl), decodeGetInvoiceReq, encodeResponse, options...)

	r := mux.NewRouter()
	r.Methods("POST").Path("/create_invoice").Handler(createInvoiceHandler)
	r.Methods("POST").Path("/").Handler(getInvoiceHandler)
	return r
}

func decodeCreateInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request models.CreateInvoiceRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeGetInvoiceReq(ctx context.Context, req *http.Request) (interface{}, error) {
	var request models.GetInvoiceRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
