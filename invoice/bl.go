package invoice

import (
	"context"
	"fmt"
	"time"

	"invoice_service/invoice/repository"
	"invoice_service/spec"

	"github.com/go-kit/log"
	"github.com/google/uuid"
)

type BL interface {
	CreateInvoice(ctx context.Context, createInvoiceReq spec.CreateInvoiceRequest) (spec.Invoice, error)
	GetInvoice(ctx context.Context, invoiceId string) (spec.Invoice, error)
	ListInvoice(ctx context.Context, invoiceFilter spec.InvoiceFilter) ([]spec.Invoice, error)
	UpdateInvoice(ctx context.Context, updateInvoiceReq spec.UpdateInvoiceRequest) (spec.Invoice, error)
	DeleteInvoice(ctx context.Context, invoiceId string) error
}

type bl struct {
	logger log.Logger
	repo   repository.Repository
}

func NewBL(logger log.Logger, repo repository.Repository) BL {
	return bl{
		logger: logger,
		repo:   repo,
	}
}

func (bl bl) CreateInvoice(ctx context.Context, createInvoiceReq spec.CreateInvoiceRequest) (spec.Invoice, error) {
	var (
		invoice spec.Invoice
		err     error
	)
	createInvoiceReq.Id = uuid.NewString()
	createInvoiceReq.CreatedAt = time.Now()
	createInvoiceReq.UpdatedAt = time.Now()

	if createInvoiceReq.AdminId == createInvoiceReq.UserId {
		bl.logger.Log("[debug]", "User id and admin id cant be same")
		return invoice, fmt.Errorf("user id and admin id cant be same")
	}

	invoice, err = bl.repo.Create(ctx, createInvoiceReq)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to create invoice ", "err", err.Error())
		return invoice, fmt.Errorf("failed to create invoice %v", err.Error())
	}

	return invoice, nil
}

func (bl bl) GetInvoice(ctx context.Context, invoiceId string) (spec.Invoice, error) {
	var (
		invoice spec.Invoice
		err     error
	)
	invoice, err = bl.repo.Get(ctx, invoiceId)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to get invoice ", "err", err.Error())
		return invoice, err
	}
	bl.logger.Log("[debug]", "Successfully get invoice")
	return invoice, nil
}

func (bl bl) ListInvoice(ctx context.Context, invoiceFilter spec.InvoiceFilter) ([]spec.Invoice, error) {
	var (
		invoices []spec.Invoice
		err      error
	)
	invoices, err = bl.repo.List(ctx, invoiceFilter)
	if err != nil {
		return invoices, err
	}
	if len(invoices) == 0 {
		bl.logger.Log("[debug]", "No invoice found")
	}
	bl.logger.Log("[debug]", "Successfully listed invoice")
	return invoices, nil
}

func (bl bl) UpdateInvoice(ctx context.Context, updateInvoiceReq spec.UpdateInvoiceRequest) (spec.Invoice, error) {
	var (
		invoice spec.Invoice
		err     error
	)
	err = bl.repo.Edit(ctx, updateInvoiceReq)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to update invoice", "err", err.Error())
		return invoice, err
	}

	invoice, err = bl.repo.Get(ctx, updateInvoiceReq.Id)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to get updated invoice", "err", err.Error())
		return invoice, err
	}
	bl.logger.Log("Successfully updated invoice")
	return invoice, nil
}

func (bl bl) DeleteInvoice(ctx context.Context, invoiceId string) error {
	err := bl.repo.Delete(ctx, invoiceId)
	if err != nil {
		bl.logger.Log("[debug]", "Failed to delete invoice", "err", err.Error())
		return err
	}
	return nil
}
