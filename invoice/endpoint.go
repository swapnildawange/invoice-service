package invoice

import (
	"context"
	"fmt"
	"invoice_service/model"
	"invoice_service/security"
	"time"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Endpoints struct {
	CreateInvoice endpoint.Endpoint
	GetInvoice    endpoint.Endpoint
	ListInvoice   endpoint.Endpoint
	UpdateInvoice endpoint.Endpoint
	DeleteInvoice endpoint.Endpoint
}

func NewEndpoints(logger log.Logger, bl BL) Endpoints {
	return Endpoints{
		CreateInvoice: makeCreateInvoice(logger, bl),
		GetInvoice:    makeGetInvoice(logger, bl),
		ListInvoice:   makeListInvoice(logger, bl),
		UpdateInvoice: makeUpdateInvoice(logger, bl),
		DeleteInvoice: makeDeleteInvoiceEndpoint(logger, bl),
	}
}

func makeCreateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "createInvoice",
				"took", time.Since(begin),
			)
		}(time.Now())

		var (
			req              model.CreateInvoiceRequest
			createInvoiceRes model.Invoice
		)
		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid jwt token")
		}

		if JWTClaims.Role == 2 {
			return nil, security.NotAuthorizedErr
		}

		req = request.(model.CreateInvoiceRequest)
		// admin id
		req.AdminId = JWTClaims.Id

		createInvoiceRes, err = bl.CreateInvoice(ctx, req)
		if err != nil {
			return nil, err
		}
		return createInvoiceRes, nil
	}
}

func makeGetInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "getInvoice",
				"took", time.Since(begin),
			)
		}(time.Now())

		var (
			invoice       model.Invoice
			getInvoiceReq model.GetInvoiceRequest
		)

		getInvoiceReq = request.(model.GetInvoiceRequest)
		invoice, err = bl.GetInvoice(ctx, getInvoiceReq.Id)
		if err != nil {
			return nil, err
		}
		return invoice, nil
	}
}

func makeListInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "listInvoice",
				"took", time.Since(begin),
			)
		}(time.Now())

		var (
			invoices      []model.Invoice
			invoiceFilter model.InvoiceFilter
		)

		invoiceFilter = request.(model.InvoiceFilter)
		invoices, err = bl.ListInvoice(ctx, invoiceFilter)
		if err != nil {
			return nil, err
		}
		return invoices, nil
	}
}

func makeUpdateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "updateInvoice",
				"took", time.Since(begin),
			)
		}(time.Now())

		var (
			req     model.UpdateInvoiceRequest
			invoice model.Invoice
		)

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid jwt token")
		}

		if JWTClaims.Role == 2 {
			return nil, security.NotAuthorizedErr
		}

		req = request.(model.UpdateInvoiceRequest)
		invoice, err = bl.UpdateInvoice(ctx, req)
		if err != nil {
			return nil, err
		}
		return invoice, nil
	}
}

func makeDeleteInvoiceEndpoint(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "deleteInvoice",
				"took", time.Since(begin),
			)
		}(time.Now())

		var (
			invoiceId string
		)

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid jwt token")
		}

		if JWTClaims.Role == 2 {
			return nil, security.NotAuthorizedErr
		}

		invoiceId = request.(string)
		err = bl.DeleteInvoice(ctx, invoiceId)
		if err != nil {
			return nil, err
		}
		return "Invoice deleted", nil
	}
}
