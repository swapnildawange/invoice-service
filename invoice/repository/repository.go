package repository

import (
	"context"
	"database/sql"
	"fmt"
	"invoice_service/model"

	"github.com/go-kit/kit/log"
)

type Repository interface {
	CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error)
}

type repository struct {
	log log.Logger
	db  *sql.DB
}

func NewRepository(logger log.Logger, db *sql.DB) Repository {
	return &repository{
		log: logger,
		db:  db,
	}
}

func (repo *repository) CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error) {
	fmt.Println("here")
	var (
		invoice model.Invoice
	)
	query := `insert into "public"."invoice"("id","user_id","admin_id","paid","payment_status","created_at","updated_at")
	values 
	(E'invoice_id1',E'2',E'1',E'1000.46',E'1',E'2022-03-14 00:00:00',E'2022-03-14 00:00:00')`
	_, err := repo.db.ExecContext(ctx, query)
	if err != nil {
		fmt.Println("err", err)
		repo.log.Log(err)
		return invoice, err
	}
	return invoice, nil
}
