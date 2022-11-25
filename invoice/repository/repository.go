package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/invoice-service/spec"
	"github.com/invoice-service/svcerror"
)

type Repository interface {
	Create(ctx context.Context, createInvoiceReq spec.CreateInvoiceRequest) (spec.Invoice, error)
	List(ctx context.Context, invoiceFilter spec.InvoiceFilter) ([]spec.Invoice, error)
	Get(ctx context.Context, invoiceId string, userId int) (spec.Invoice, error)
	Edit(ctx context.Context, updateInvoiceReq spec.UpdateInvoiceRequest) error
	Delete(ctx context.Context, invoiceId string) (spec.Invoice, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) Create(ctx context.Context, createInvoiceReq spec.CreateInvoiceRequest) (spec.Invoice, error) {
	var (
		invoice spec.Invoice
		err     error
		tx      *sql.Tx
		row     *sql.Row
		userId  int
	)
	ctx, cancel := context.WithTimeout(ctx, spec.Timeout*time.Second)
	defer cancel()

	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return invoice, fmt.Errorf("failed to begin transaction %v", err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row = tx.QueryRowContext(ctx, `SELECT id FROM users where id =$1`, createInvoiceReq.UserId)
	if err = row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return invoice, svcerror.ErrUserNotFound
		}
		return invoice, err
	}

	insertQuery := `insert into "invoice"("id","user_id","admin_id","paid","payment_status","created_at","updated_at")
	values 
	($1,$2,$3,$4,$5,$6,$7) `
	_, err = tx.ExecContext(ctx, insertQuery, createInvoiceReq.Id, createInvoiceReq.UserId, createInvoiceReq.AdminId, createInvoiceReq.Paid, createInvoiceReq.PaymentStatus, createInvoiceReq.CreatedAt, createInvoiceReq.UpdatedAt)
	if err != nil {
		return invoice, fmt.Errorf("failed to execute insert invoice %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return invoice, fmt.Errorf("failed to commit transaction %v", err.Error())
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

func (repo *repository) List(ctx context.Context, invoiceFilter spec.InvoiceFilter) ([]spec.Invoice, error) {	
	var (
		invoices = make([]spec.Invoice, 0)
	)

	listInvoiceQUery := `select id,user_id,admin_id,paid,payment_status,created_at,updated_at from invoice `

	listInvoiceQUery, filterValues := repo.queryInvoiceWithFilter(listInvoiceQUery, invoiceFilter)
	rows, err := repo.db.QueryContext(ctx, listInvoiceQUery, filterValues...)
	if err != nil {
		return invoices, fmt.Errorf("failed to execute list invoice %v", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var invoice spec.Invoice
		if err := rows.Scan(&invoice.Id, &invoice.UserId, &invoice.AdminId, &invoice.Paid, &invoice.PaymentStatus, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			if err == sql.ErrNoRows {
				return invoices, nil
			}
			return invoices, fmt.Errorf("failed to scan %v", err.Error())
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil

}
func (repo *repository) queryInvoiceWithFilter(query string, filter spec.InvoiceFilter) (string, []interface{}) {

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

	filterValues = append(filterValues, spec.PageSize)
	query += ` LIMIT $` + strconv.Itoa(len(filterValues))

	filterValues = append(filterValues, (filter.Page-1)*spec.PageSize)
	query += ` OFFSET $` + strconv.Itoa(len(filterValues))

	if len(filterValues) >= 1 {
		query = strings.Replace(query, "AND", "WHERE", 1)
	}

	return query, filterValues
}

func (repo *repository) Get(ctx context.Context, invoiceId string, userId int) (spec.Invoice, error) {
	var (
		invoice   spec.Invoice
		row       *sql.Row
		err       error
		arguments []interface{}
	)
	getInvoiceQuery := `SELECT FROM invoice WHERE id =$1 ` + invoiceId
	arguments = append(arguments, invoiceId)
	if userId > 0 {
		getInvoiceQuery += `AND user_id = $2`
		arguments = append(arguments, userId)
	}
	row = repo.db.QueryRowContext(ctx, getInvoiceQuery, arguments...)

	if row.Err() != nil {
		return invoice, fmt.Errorf("failed to execute get invoice query %v", row.Err())
	}
	if err = row.Scan(&invoice.Id, &invoice.UserId, &invoice.AdminId, &invoice.Paid, &invoice.PaymentStatus, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return invoice, svcerror.ErrInvoiceNotFound
		}
		return invoice, fmt.Errorf("failed to scan row %v", err.Error())
	}
	return invoice, nil
}

func (repo *repository) Edit(ctx context.Context, updateInvoiceReq spec.UpdateInvoiceRequest) error {
	var (
		err    error
		values []interface{}
	)

	updateQuery := `UPDATE invoice SET `

	if updateInvoiceReq.Paid != -1 {
		values = append(values, updateInvoiceReq.Paid)
		updateQuery += ` paid = $` + strconv.Itoa(len(values)) + ` , `
	}

	if updateInvoiceReq.PaymentStatus != -1 {
		values = append(values, updateInvoiceReq.PaymentStatus)
		updateQuery += ` payment_status = $` + strconv.Itoa(len(values)) + ` , `
	}

	values = append(values, time.Now())
	updateQuery += ` updated_at = $` + strconv.Itoa(len(values)) + ` , `

	values = append(values, updateInvoiceReq.Id)
	updateQuery += `WHERE id = $` + strconv.Itoa(len(values))

	updateQuery = strings.Replace(updateQuery, ", WHERE", "WHERE", 1)

	_, err = repo.db.ExecContext(ctx, updateQuery, values...)
	if err != nil {
		return fmt.Errorf("failed to execute update invoice query %v", err.Error())
	}
	return nil
}

func (repo *repository) Delete(ctx context.Context, invoiceId string) (spec.Invoice, error) {
	var (
		tx      *sql.Tx
		err     error
		invoice spec.Invoice
	)

	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return invoice, fmt.Errorf("failed to begin transaction %v", err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	invoice, err = repo.Get(ctx, invoiceId, 0)
	if err != nil {
		return invoice, err
	}
	_, err = tx.ExecContext(ctx, "DELETE from invoice where id = $1", invoiceId)
	if err != nil {
		return invoice, fmt.Errorf("failed to delete invoice %v", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return invoice, fmt.Errorf("failed to commit transaction %v", err.Error())
	}
	return invoice, nil
}
