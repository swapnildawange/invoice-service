package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
	"github.com/invoice-service/svcerror"

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
		DeleteInvoice: makeDeleteInvoice(logger, bl),
	}
}

func makeCreateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req              spec.CreateInvoiceRequest
			createInvoiceRes spec.Invoice
		)
		defer func(begin time.Time) {
			logger.Log(
				"method", "createInvoice",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%+v", request),
				"reponse", fmt.Sprintf("%+v", createInvoiceRes),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		req = request.(spec.CreateInvoiceRequest)
		// admin id
		req.AdminId = JWTClaims.Id

		createInvoiceRes, err = bl.CreateInvoice(ctx, req)
		if err != nil {
			logger.Log("[debug]", "failed to create invoice", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return nil, svcerror.ErrFailedToCreateInvoice
		}
		return createInvoiceRes, nil
	}
}

func makeGetInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			invoice       spec.Invoice
			getInvoiceReq spec.GetInvoiceRequest
		)
		defer func(begin time.Time) {
			logger.Log(
				"method", "getInvoice",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%+v", request),
				"reponse", fmt.Sprintf("%+v", invoice),
			)
		}(time.Now())

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}
		getInvoiceReq = request.(spec.GetInvoiceRequest)
		if JWTClaims.Role == int(spec.RoleAdmin) {
			// admin can get any invoice
			invoice, err = bl.GetInvoice(ctx, getInvoiceReq)

		} else if JWTClaims.Role == int(spec.RoleUser) {
			// but user can get only his invoice
			getInvoiceReq.UserId = JWTClaims.Id
			invoice, err = bl.GetInvoice(ctx, getInvoiceReq)
		}

		if err != nil {
			logger.Log("[debug]", "failed to get invoice", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return nil, svcerror.ErrFailedToGetInvoice
		}
		return invoice, nil
	}
}

func makeListInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			invoices      []spec.Invoice
			invoiceFilter spec.InvoiceFilter
		)
		defer func(begin time.Time) {
			logger.Log(
				"method", "listInvoice",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%+v", request),
				"reponse", fmt.Sprintf("%+v", invoices),
			)
		}(time.Now())

		invoiceFilter = request.(spec.InvoiceFilter)

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleAdmin) {
			// admin can get any invoice
			invoices, err = bl.ListInvoice(ctx, invoiceFilter)

		} else if JWTClaims.Role == int(spec.RoleUser) {
			// but user can get only his invoice
			invoiceFilter.UserId = JWTClaims.Id
			invoices, err = bl.ListInvoice(ctx, invoiceFilter)
		}
		if err != nil {
			logger.Log("[debug]", "failed to list invoices", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
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
				"request", fmt.Sprintf("%+v", request),
				"reponse", fmt.Sprintf("%+v", response),
			)
		}(time.Now())

		var (
			req     spec.UpdateInvoiceRequest
			invoice spec.Invoice
		)

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken

		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		req = request.(spec.UpdateInvoiceRequest)
		invoice, err = bl.UpdateInvoice(ctx, req)
		if err != nil {
			logger.Log("[debug]", "failed to update invoice", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return nil, svcerror.ErrFailedToUpdateInvoice
		}
		return invoice, nil
	}
}

func makeDeleteInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		defer func(begin time.Time) {
			logger.Log(
				"method", "deleteInvoice",
				"took", time.Since(begin),
				"request", fmt.Sprintf("%+v", request),
				"reponse", fmt.Sprintf("%+v", response),
			)
		}(time.Now())

		var invoiceId string

		JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
		if !ok {
			logger.Log("[debug]", "Invalid JWT token")
			return nil, svcerror.ErrInvalidToken
		}

		if JWTClaims.Role == int(spec.RoleUser) {
			logger.Log("[debug]", "User is not authorized")
			return nil, svcerror.ErrNotAuthorized
		}

		invoiceId = request.(string)
		invoiceId, err = bl.DeleteInvoice(ctx, invoiceId)
		if err != nil {
			logger.Log("[debug]", "failed to update invoice", "err", err)
			_, ok := err.(*svcerror.CustomErrString)
			if ok {
				return nil, err
			}
			return nil, svcerror.ErrFailedToDeleteInvoice
		}
		return invoiceId, nil
	}
}
