package invoice

import (
	"context"
	"fmt"
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

	if createInvoiceReq.AdminId == createInvoiceReq.UserId {
		bl.logger.Log("invoice", "bl", "CreateInvoice", "User id and admin id cant be same")
		return invoice, fmt.Errorf("User id and admin id cant be same")
	}

	invoice, err = bl.repo.CreateInvoice(ctx, createInvoiceReq)
	if err != nil {
		bl.logger.Log("invoice", "bl", "CreateInvoice", "Failed to create invoice", err.Error())
		return invoice, fmt.Errorf("Failed to create invoice", err.Error())
	}

	return invoice, nil
}

func (bl *bl) GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error) {
	var (
		invoice model.Invoice
		err     error
	)
	invoice, err = bl.repo.GetInvoice(ctx, invoiceId)
	if err != nil {
		bl.logger.Log("Failed to get invoice", err.Error())
		return invoice, err
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
	var (
		invoice model.Invoice
		err     error
	)
	err = bl.repo.EditInvoice(ctx, updateInvoiceReq)
	if err != nil {
		bl.logger.Log("Failed to update invoice", err.Error())
		return invoice, err
	}

	invoice, err = bl.repo.GetInvoice(ctx, updateInvoiceReq.Id)
	if err != nil {
		bl.logger.Log("Failed to get updated invoice", err.Error())
		return invoice, err
	}
	bl.logger.Log("Successfully updated invoice")
	return invoice, nil
}

func (bl *bl) DeleteInvoice(ctx context.Context, invoiceId string) error {
	err := bl.repo.DeleteInvoice(ctx, invoiceId)
	if err != nil {
		bl.logger.Log("Fialed to delete invoice")
		return err
	}
	return nil
}
