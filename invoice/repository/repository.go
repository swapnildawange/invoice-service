package repository

import (
	"context"
	"database/sql"
	"fmt"
	"invoice_service/model"
	"time"
)

type Repository interface {
	CreateInvoice(ctx context.Context, createInvoiceReq model.CreateInvoiceRequest) (model.Invoice, error)
	ListInvoice(ctx context.Context, userId int) ([]model.Invoice, error)
	CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
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

func (repo *repository) ListInvoice(ctx context.Context, userId int) ([]model.Invoice, error) {
	var (
		invoices = make([]model.Invoice, 0)
	)

	list := `select * from invoice where user_id = $1 ;`
	rows, err := repo.db.QueryContext(ctx, list, userId)
	if err != nil {
		return invoices, err
	}

	defer rows.Close()

	for rows.Next() {
		var invoice model.Invoice
		if err := rows.Scan(&invoice.Id, &invoice.UserId, &invoice.Paid, &invoice.PaymentStatus, &invoice.AdminId, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
			return invoices, fmt.Errorf("Failed to scan ", err.Error())
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil

}

func (repo *repository) CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (model.User, error) {

	var (
		user      model.User
		err       error
		rows      *sql.Rows
		row       *sql.Row
		numOfRows int64
		tx        *sql.Tx
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
		return user, fmt.Errorf("failed to begin transaction")
	}

	// check if email is already present
	checkEmailQuery := `select COUNT(*) from users where email = $1`

	rows, err = tx.QueryContext(ctx, checkEmailQuery, createUserReq.Email)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	var emailCount int

	for rows.Next() {
		if err := rows.Scan(&emailCount); err != nil {
			return user, err
		}
	}

	if emailCount != 0 {
		return user, fmt.Errorf("email is already present")
	}

	insertQuery := `insert into "users"("email","first_name","last_name","password","role","created_at","updated_at")
	values 
	($1,$2,$3,$4,$5,$6,$7) `
	result, err := tx.ExecContext(ctx, insertQuery, createUserReq.Email, createUserReq.FirstName, createUserReq.LastName, createUserReq.Password, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt)
	if err != nil {
		return user, err
	}

	numOfRows, err = result.RowsAffected()
	if err != nil {
		return user, err
	}
	if numOfRows != 1 {
		return user, fmt.Errorf("expected to affect 1 row, affected %d", numOfRows)
	}

	// get newly created user from database
	getUserQuery := `select id,email,first_name,last_name,role,created_at,updated_at from users where email = $1`
	row = tx.QueryRowContext(ctx, getUserQuery, createUserReq.Email)

	if err = row.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, fmt.Errorf("Failed to get user details %s", row.Err().Error())
	}
	tx.Commit()

	return user, nil
}

func (repo *repository) ListUsers(ctx context.Context) ([]model.User, error) {
	var (
		users = make([]model.User, 0)
		err   error
		rows  *sql.Rows
	)
	listUsersQuery := `select id,first_name,last_name,email,created_at,updated_at from users ;`
	rows, err = repo.db.QueryContext(ctx, listUsersQuery)
	if err != nil {
		return users, fmt.Errorf("Failed to list users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
