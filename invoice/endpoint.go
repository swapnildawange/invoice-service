package invoice

import (
	"context"
	"fmt"
	"invoice_service/model"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Endpoints struct {
	CreateInvoice endpoint.Endpoint
	GetInvoice    endpoint.Endpoint
	ListInvoice   endpoint.Endpoint
	UpdateInvoice endpoint.Endpoint
	DeleteInvoice endpoint.Endpoint
	CreateUser    endpoint.Endpoint
	ListUsers     endpoint.Endpoint
}

func NewEndpoints(logger log.Logger, bl BL) Endpoints {
	return Endpoints{
		CreateInvoice: makeCreateInvoice(logger, bl),
		GetInvoice:    makeGetInvoice(logger, bl),
		ListInvoice:   makeListInvoice(logger, bl),
		UpdateInvoice: makeUpdateInvoice(logger, bl),
		DeleteInvoice: makeDeleteEndpoint(logger, bl),
		CreateUser:    makeCreateUser(logger, bl),
		ListUsers:     makeListUsers(logger, bl),
	}
}

func makeCreateUser(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req  model.CreateUserRequest
			user model.User
		)
		req = request.(model.CreateUserRequest)
		if err != nil {
			return nil, fmt.Errorf("invalid request for create user")
		}
		user, err = bl.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func makeCreateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req              model.CreateInvoiceRequest
			createInvoiceRes model.Invoice
		)

		req = request.(model.CreateInvoiceRequest)
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
			invoice   model.Invoice
			invoiceId string
		)

		invoiceId = request.(string)
		invoice, err = bl.GetInvoice(ctx, invoiceId)
		if err != nil {
			return nil, err
		}
		return invoice, nil
	}
}

func makeListInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			invoices []model.Invoice
			userId   int
		)

		userId = request.(int)
		invoices, err = bl.ListInvoice(ctx, userId)
		if err != nil {
			return nil, err
		}
		return invoices, nil
	}
}

func makeUpdateInvoice(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			req     model.UpdateInvoiceRequest
			invoice model.Invoice
		)

		req = request.(model.UpdateInvoiceRequest)
		invoice, err = bl.UpdateInvoice(ctx, req)
		if err != nil {
			return nil, err
		}
		return invoice, nil
	}
}

func makeDeleteEndpoint(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			invoiceId string
		)

		invoiceId = request.(string)
		err = bl.DeleteInvoice(ctx, invoiceId)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("Invoice deleted"), nil
	}
}

func makeListUsers(logger log.Logger, bl BL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		response, nil := bl.ListUsers(ctx)
		if err != nil {
			logger.Log("endpoint", "makeListUsers", "Failed to list users", err.Error())
			return
		}
		return

	}
}
