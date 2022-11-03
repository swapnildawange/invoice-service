package invoice

import (
	"context"
	"time"

	"invoice_service/invoice/repository"
	"invoice_service/model"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type BL interface {
	CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error)
	GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error)
	ListInvoice(ctx context.Context, userId int) ([]model.Invoice, error)
	UpdateInvoice(ctx context.Context, updateInvoiceReq model.UpdateInvoiceRequest) (model.Invoice, error)
	DeleteInvoice(ctx context.Context, invoiceId string) error
	CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error)
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
	createInvoiceReq.Id = uuid.NewString()
	createInvoiceReq.CreatedAt = time.Now()
	createInvoiceReq.UpdatedAt = time.Now()
	createInvoiceReq.AdminId = 1
	createInvoiceReq.UserId = 1
	createInvoiceReq.PaymentStatus = 3

	invoice, err = bl.repo.CreateInvoice(ctx, createInvoiceReq)
	if err != nil {
		bl.logger.Log("invoice", "bl", "CreateInvoice", "Failed to create invoice", err.Error())
		return invoice, err
	}

	return invoice, nil
}

func (bl *bl) GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error) {
	var invoice = model.Invoice{
		Id:            "temp_invoice_id",
		UserId:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Paid:          100,
		AdminId:       2,
		PaymentStatus: model.PaymentStatus(model.Initiate),
	}
	bl.logger.Log("Successfully get invoice")
	return invoice, nil
}

func (bl *bl) ListInvoice(ctx context.Context, userId int) ([]model.Invoice, error) {
	var (
		invoices []model.Invoice
		err      error
	)
	invoices, err = bl.repo.ListInvoice(ctx, userId)
	if err != nil {
		return invoices, err
	}
	if len(invoices) == 0 {
		bl.logger.Log("No invoice found")
	}
	bl.logger.Log("Successfully listed invoice")
	return invoices, nil
}

func (bl *bl) UpdateInvoice(ctx context.Context, updateInvoiceReq model.UpdateInvoiceRequest) (model.Invoice, error) {
	// invoice
	var invoice = model.Invoice{
		UserId:        updateInvoiceReq.Invoice.UserId,
		Id:            updateInvoiceReq.Invoice.Id,
		CreatedAt:     updateInvoiceReq.Invoice.CreatedAt,
		UpdatedAt:     time.Now(),
		Paid:          updateInvoiceReq.Invoice.Paid,
		AdminId:       updateInvoiceReq.Invoice.AdminId,
		PaymentStatus: model.PaymentStatus(model.Initiate),
	}
	bl.logger.Log("Successfully get invoice")
	return invoice, nil
}

func (bl *bl) DeleteInvoice(ctx context.Context, invoiceId string) error {
	return nil
}

func (bl *bl) CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error) {

	createUserReq.CreatedAt = time.Now()
	createUserReq.UpdatedAt = time.Now()
	user, err := bl.repo.CreateUser(ctx, createUserReq)
	if err != nil {
		bl.logger.Log(err)
		return user, err
	}
	return user, nil
}
