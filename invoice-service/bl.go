package invoice

import (
	"context"
	"time"

	"invoicing/invoice-service/models"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type BL interface {
	CreateInvoice(ctx context.Context, createInvoiceReq models.CreateInvoiceRequest) (models.Invoice, error)
	GetInvoice(ctx context.Context, getInoiceReq models.GetInvoiceRequest) (models.Invoice, error)
}

type bl struct {
	logger log.Logger
}

func NewBL(logger log.Logger) *bl {
	return &bl{
		logger: logger,
	}
}

func (bl *bl) CreateInvoice(ctx context.Context, createInvoiceReq models.CreateInvoiceRequest) (models.Invoice, error) {
	bl.logger.Log("Creating invoice")

	var (
		admin            models.Admin
		createInvoiceRes = models.Invoice{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Paid:      createInvoiceReq.Paid,
			User: models.User{
				Id: createInvoiceReq.UserId,
			},
			Id:                 uuid.NewString(),
			PaymentInitiatedBy: admin,
			PaymentStatus:      models.PaymentStatus(1),
		}
	)

	return createInvoiceRes, nil
}

func (bl *bl) GetInvoice(ctx context.Context, getInvoiceReq models.GetInvoiceRequest) (models.Invoice, error) {
	var invoice = models.Invoice{
		User: models.User{
			Id: 1,
		},
		Id:        "temp_invoice_id",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Paid:      100,
		PaymentInitiatedBy: models.Admin{
			Id: 2,
		},
		PaymentStatus: models.PaymentStatus(models.Initiate),
	}
	bl.logger.Log("Successfully get invoice")
	return invoice, nil
}

func (bl *bl) UpdateInvoice(ctx context.Context, updateInvoiceReq models.UpdateInvoiceRequest) (models.Invoice, error) {
	// invoice
}


