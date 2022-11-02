package invoice

import (
	"context"
	"time"

	"invoice_service/invoice/repository"
	"invoice_service/model"

	"github.com/go-kit/kit/log"
)

type BL interface {
	CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error)
	GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error)
	UpdateInvoice(ctx context.Context, updateInvoiceReq model.UpdateInvoiceRequest) (model.Invoice, error)
	DeleteInvoice(ctx context.Context, invoiceId string) error
}

type bl struct {
	logger log.Logger
	repo   repository.Repository
}

func NewBL(logger log.Logger, repo repository.Repository) BL {
	return &bl{
		logger: logger,
		repo:   repo,
	}
}

func (bl *bl) CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error) {
	bl.logger.Log("Creating invoice")

	var (
		invoice model.Invoice
		err     error
	)
	// var (
	// 	createInvoiceRes = model.Invoice{
	// 		CreatedAt:          time.Now(),
	// 		UpdatedAt:          time.Now(),
	// 		Paid:               createInvoiceReq.Paid,
	// 		PaymentInitiatedBy: 1,
	// 		PaymentStatus:      model.PaymentStatus(1),
	// 	}
	// )

	invoice, err = bl.repo.CreateInvoice(ctx, createInvoiceReq)
	if err != nil {
		bl.logger.Log("invoice", "bl", "CreateInvoice", "Failed to create invoice", err.Error())
		return invoice, err
	}

	return invoice, nil
}

func (bl *bl) GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error) {
	var invoice = model.Invoice{
		Id:                 "temp_invoice_id",
		UserId:             1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Paid:               100,
		PaymentInitiatedBy: 2,
		PaymentStatus:      model.PaymentStatus(model.Initiate),
	}
	bl.logger.Log("Successfully get invoice")
	return invoice, nil
}

func (bl *bl) UpdateInvoice(ctx context.Context, updateInvoiceReq model.UpdateInvoiceRequest) (model.Invoice, error) {
	// invoice
	var invoice = model.Invoice{
		UserId:             updateInvoiceReq.Invoice.UserId,
		Id:                 updateInvoiceReq.Invoice.Id,
		CreatedAt:          updateInvoiceReq.Invoice.CreatedAt,
		UpdatedAt:          time.Now(),
		Paid:               updateInvoiceReq.Invoice.Paid,
		PaymentInitiatedBy: updateInvoiceReq.Invoice.PaymentInitiatedBy,
		PaymentStatus:      model.PaymentStatus(model.Initiate),
	}
	bl.logger.Log("Successfully get invoice")
	return invoice, nil
}

func (bl *bl) DeleteInvoice(ctx context.Context, invoiceId string) error {
	return nil
}
