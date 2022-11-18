package invoice

import (
	"context"
	"invoice_service/security"
	"invoice_service/spec"
	"invoice_service/svcerror"
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
			req              spec.CreateInvoiceRequest
			createInvoiceRes spec.Invoice
		)
		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			return nil, svcerror.ErrNotAuthorized
		}

		req = request.(spec.CreateInvoiceRequest)
		// admin id
		req.AdminId = JWTClaims.Id

		createInvoiceRes, err = bl.CreateInvoice(ctx, req)
		if err != nil {
			return nil, svcerror.ErrFailedToCreateInvoice
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
			invoice       spec.Invoice
			getInvoiceReq spec.GetInvoiceRequest
		)

		getInvoiceReq = request.(spec.GetInvoiceRequest)
		invoice, err = bl.GetInvoice(ctx, getInvoiceReq.Id)
		if err != nil {
			return nil, svcerror.ErrFailedToGetInvoice
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
			invoices      []spec.Invoice
			invoiceFilter spec.InvoiceFilter
		)

		invoiceFilter = request.(spec.InvoiceFilter)
		invoices, err = bl.ListInvoice(ctx, invoiceFilter)
		if err != nil {
			return nil, svcerror.ErrFailedToListInvoice
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
			req     spec.UpdateInvoiceRequest
			invoice spec.Invoice
		)

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, svcerror.ErrInvalidToken

		}

		if JWTClaims.Role == int(spec.RoleUser) {
			return nil, svcerror.ErrNotAuthorized
		}

		req = request.(spec.UpdateInvoiceRequest)
		invoice, err = bl.UpdateInvoice(ctx, req)
		if err != nil {
			return nil, svcerror.ErrFailedToUpdateInvoice
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

		var invoiceId string

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			return nil, svcerror.ErrNotAuthorized
		}

		invoiceId = request.(string)
		err = bl.DeleteInvoice(ctx, invoiceId)
		if err != nil {
			return nil, svcerror.ErrFailedToDeleteInvoice
		}
		return "Invoice deleted", nil
	}
}
