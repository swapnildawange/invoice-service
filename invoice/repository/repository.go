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
	ListInvoice(ctx context.Context, userId int) ([]model.Invoice, error)
	CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error)
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
	var (
		invoice model.Invoice
		err     error
	)
	insertQuery := `insert into "invoice"("id","user_id","admin_id","paid","payment_status","created_at","updated_at")
	values 
	($1,$2,$3,$4,$5,$6,$7)`
	_, err = repo.db.ExecContext(ctx, insertQuery, createInvoiceReq.Id, createInvoiceReq.UserId, createInvoiceReq.AdminId, createInvoiceReq.Paid, createInvoiceReq.PaymentStatus, createInvoiceReq.CreatedAt, createInvoiceReq.UpdatedAt)
	if err != nil {
		repo.log.Log(err)
		return invoice, err
	}

	invoice.Id = createInvoiceReq.Id
	invoice.AdminId = createInvoiceReq.AdminId
	invoice.UserId = createInvoiceReq.UserId
	invoice.Paid = createInvoiceReq.Paid
	invoice.PaymentStatus = createInvoiceReq.PaymentStatus
	invoice.CreatedAt = createInvoiceReq.CreatedAt
	invoice.UpdatedAt = createInvoiceReq.UpdatedAt

	return invoice, nil
}

func (repo *repository) ListInvoice(ctx context.Context, userId int) ([]model.Invoice, error) {
	var (
		invoices = make([]model.Invoice, 0)
	)

	list := `select * from invoice where user_id = $1 ;`
	rows, err := repo.db.QueryContext(ctx, list, userId)
	if err != nil {
		repo.log.Log(err, userId)
		return invoices, err
	}

	defer rows.Close()

	for rows.Next() {
		var invoice model.Invoice
		if err := rows.Scan(&invoice.Id, &invoice.UserId, &invoice.Paid, &invoice.PaymentStatus, &invoice.AdminId, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			repo.log.Log("Failed to scan")
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil

}

func (repo *repository) CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error) {

	var user model.User
	insertQuery := `insert into "users"("email","first_name","last_name","password","role","created_at","updated_at")
	values 
	($1,$2,$3,$4,$5,$6,$7) `
	result, err := repo.db.ExecContext(ctx, insertQuery, createUserReq.Email, createUserReq.FirstName, createUserReq.LastName, createUserReq.Password, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt)
	if err != nil {
		repo.log.Log(err)
		return user, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		repo.log.Log(err)
		return user, err
	}
	if rows != 1 {
		repo.log.Log(fmt.Sprintf("expected to affect 1 row, affected %d", rows))
		return user, err
	}
	return user, nil
}
