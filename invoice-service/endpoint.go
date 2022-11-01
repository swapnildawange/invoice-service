package invoice

import (
	"context"
	"invoicing/invoice-service/models"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// type Endpoints struct {
// 	CreateInvoice endpoint.Endpoint
// }

// func MakeEndpoints(logger log.Logger, bl BL) Endpoints {
// 	return Endpoints{
// 		CreateInvoice: makeCreateInvoice(logger, bl),
// 	}
// }

func makeCreateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req              models.CreateInvoiceRequest
			createInvoiceRes models.Invoice
		)

		req = request.(models.CreateInvoiceRequest)
		createInvoiceRes, err = bl.CreateInvoice(ctx, req)
		if err != nil {
			return nil, err
		}
		return createInvoiceRes, nil
	}
}

func makeGetInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req     models.GetInvoiceRequest
			invoice models.Invoice
		)
		req = request.(models.GetInvoiceRequest)
		invoice, err = bl.GetInvoice(ctx, req)
		if err != nil {
			return nil, err
		}
		return invoice, nil
	}
}

func makeUpdateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req     models.UpdateInvoiceRequest
			invoice models.Invoice
		)

		req = request.(models.UpdateInvoiceRequest)
		invoice, err = bl.UpdateInvoice(ctx, req)
		if err != nil {
			return nil, err
		}
		return invoice, nil
	}
}



transport  endpoint bl dl 