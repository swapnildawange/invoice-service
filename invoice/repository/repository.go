package repository

import (
	"context"
	"database/sql"
	"fmt"
	"invoice_service/model"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error)
	ListInvoice(ctx context.Context, invoiceFilter model.InvoiceFilter) ([]model.Invoice, error)
	GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error)
	EditInvoice(ctx context.Context, updateInvoiceReq model.UpdateInvoiceRequest) error
	DeleteInvoice(ctx context.Context, invoiceId string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error) {
	var (
		invoice model.Invoice
		err     error
		tx      *sql.Tx
	)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err = repo.db.BeginTx(ctx, nil)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return invoice, fmt.Errorf("Failed to begin transaction", err.Error())
	}

	insertQuery := `insert into "invoice"("id","user_id","admin_id","paid","payment_status","created_at","updated_at")
	values 
	($1,$2,$3,$4,$5,$6,$7) `
	_, err = tx.ExecContext(ctx, insertQuery, createInvoiceReq.Id, createInvoiceReq.UserId, createInvoiceReq.AdminId, createInvoiceReq.Paid, createInvoiceReq.PaymentStatus, createInvoiceReq.CreatedAt, createInvoiceReq.UpdatedAt)
	if err != nil {
		return invoice, err
	}

	err = tx.Commit()
	if err != nil {
		return invoice, fmt.Errorf("Failed to commit transaction", err.Error())
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

func (repo *repository) ListInvoice(ctx context.Context, invoiceFilter model.InvoiceFilter) ([]model.Invoice, error) {
	var (
		invoices = make([]model.Invoice, 0)
	)

	listInvoiceQUery := `select id,user_id,admin_id,paid,payment_status,created_at,updated_at from invoice `

	listInvoiceQUery, filterValues := repo.queryInvoiceWithFilter(listInvoiceQUery, invoiceFilter)

	rows, err := repo.db.QueryContext(ctx, listInvoiceQUery, filterValues...)
	if err != nil {
		return invoices, err
	}

	defer rows.Close()

	for rows.Next() {
		var invoice model.Invoice
		if err := rows.Scan(&invoice.Id, &invoice.UserId, &invoice.AdminId, &invoice.Paid, &invoice.PaymentStatus, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return invoices, fmt.Errorf("Failed to scan ", err.Error())
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil

}
func (repo *repository) queryInvoiceWithFilter(query string, filter model.InvoiceFilter) (string, []interface{}) {

	var filterValues []interface{}

	if filter.Id != "" {
		filterValues = append(filterValues, filter.Id)
		query += ` AND id = $` + strconv.Itoa(len(filterValues))
	}
	if filter.UserId != 0 {
		filterValues = append(filterValues, filter.UserId)
		query += ` AND user_id = $` + strconv.Itoa(len(filterValues))
	}
	if filter.AdminId != 0 {
		filterValues = append(filterValues, filter.AdminId)
		query += ` AND admin_id = $` + strconv.Itoa(len(filterValues))
	}
	if filter.Paid >= 0 {
		filterValues = append(filterValues, filter.Paid)
		query += ` AND paid = $` + strconv.Itoa(len(filterValues))
	}
	if filter.PaymentStatus > 0 {
		filterValues = append(filterValues, filter.PaymentStatus)
		query += ` AND payment_status = $` + strconv.Itoa(len(filterValues))
	}

	query += fmt.Sprintf(` ORDER BY %s %s`, filter.SortBy, filter.SortOrder)

	filterValues = append(filterValues, model.PageSize)
	query += ` LIMIT $` + strconv.Itoa(len(filterValues))

	filterValues = append(filterValues, (filter.Page-1)*model.PageSize)
	query += ` OFFSET $` + strconv.Itoa(len(filterValues))

	if len(filterValues) >= 1 {
		query = strings.Replace(query, "AND", "WHERE", 1)
	}

	return query, filterValues

}

func (repo *repository) GetInvoice(ctx context.Context, invoiceId string) (model.Invoice, error) {
	var (
		invoice model.Invoice
		row     *sql.Row
		err     error
	)
	row = repo.db.QueryRowContext(ctx, `select * from invoice where id=$1`, invoiceId)
	if row.Err() != nil {
		return invoice, row.Err()
	}
	if err = row.Scan(&invoice.Id, &invoice.UserId, &invoice.AdminId, &invoice.Paid, &invoice.PaymentStatus, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
		return invoice, err
	}
	return invoice, nil
}

func (repo *repository) EditInvoice(ctx context.Context, updateInvoiceReq model.UpdateInvoiceRequest) error {
	var err error
	_, err = repo.db.ExecContext(ctx, `update invoice set paid = $1 ,payment_status =$2 ,updated_at=$3 where id=$4`, updateInvoiceReq.Paid, updateInvoiceReq.PaymentStatus, time.Now(), updateInvoiceReq.Id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repository) DeleteInvoice(ctx context.Context, invoiceId string) error {
	var (
		tx  *sql.Tx
		err error
	)

	tx, err = repo.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return fmt.Errorf("Failed to begin transaction", err.Error())
	}

	_, err = tx.ExecContext(ctx, "delete from invoice where id = $1", invoiceId)
	if err != nil {
		return fmt.Errorf("Failed to delete invoice", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction", err.Error())
	}
	return nil
}
