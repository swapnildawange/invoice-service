package repository

import (
	"context"
	"database/sql"
	"fmt"
	"invoice_service/model"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (int, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUserFromAuth(ctx context.Context, email string) (userId int, hashedPassword string, err error)
	GetUser(ctx context.Context, userId int) (model.User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetUserFromAuth(ctx context.Context, email string) (userId int, hashedPassword string, err error) {
	var (
		row *sql.Row
	)
	row = repo.db.QueryRowContext(ctx, `select user_id,password from auth where email = $1`, email)
	if row.Err() != nil {
		return userId, hashedPassword, fmt.Errorf("Failed to get User details from auth", row.Err())
	}
	if err = row.Scan(&userId, &hashedPassword); err != nil {
		return userId, hashedPassword, fmt.Errorf("Failed to get user details %s", row.Err().Error())
	}
	return userId, hashedPassword, nil
}

func (repo *repository) GetUser(ctx context.Context, userId int) (model.User, error) {
	var (
		user model.User
		err  error
		row  *sql.Row
	)
	getUserQuery := `select id,first_name,last_name,role,created_at,updated_at from users where id=$1;`
	row = repo.db.QueryRowContext(ctx, getUserQuery, userId)
	if row.Err() != nil {
		return user, fmt.Errorf("Failed to get user", row.Err())
	}

	if err = row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}
	return user, nil
}

func (repo *repository) CreateUser(ctx context.Context, createUserReq model.CreateUserRequest) (int, error) {

	var (
		err    error
		rows   *sql.Rows
		row    *sql.Row
		tx     *sql.Tx
		userId int
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err = repo.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			fmt.Println("ee", err)
			tx.Rollback()
		}
	}()

	if err != nil {
		return userId, fmt.Errorf("failed to begin transaction")
	}

	// check if email is already present
	checkEmailQuery := `select COUNT(*) from auth where email = $1`

	rows, err = tx.QueryContext(ctx, checkEmailQuery, createUserReq.Email)
	if err != nil {
		return userId, err
	}
	defer rows.Close()

	var emailCount int

	for rows.Next() {
		if err := rows.Scan(&emailCount); err != nil {
			return userId, err
		}
	}

	if emailCount != 0 {
		return userId, fmt.Errorf("email is already present")
	}

	insertIntoUsersQuery := `insert into "users"("first_name","last_name","role","created_at","updated_at")
	values 
	($1,$2,$3,$4,$5) RETURNING id`

	row = tx.QueryRowContext(ctx, insertIntoUsersQuery, createUserReq.FirstName, createUserReq.LastName, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt)
	if row.Err() != nil {
		return userId, err
	}

	if err = row.Scan(&userId); err != nil {
		return userId, err
	}

	insertIntoAuthQuery := `insert into auth("user_id","email","password") values($1,$2,$3) ;`
	_, err = tx.ExecContext(ctx, insertIntoAuthQuery, userId, createUserReq.Email, createUserReq.Password)
	if err != nil {
		return userId, err
	}

	tx.Commit()

	return userId, nil
}

func (repo *repository) ListUsers(ctx context.Context) ([]model.User, error) {
	var (
		users = make([]model.User, 0)
		err   error
		rows  *sql.Rows
	)
	listUsersQuery := `select id,first_name,last_name,role,created_at,updated_at from users ;`
	rows, err = repo.db.QueryContext(ctx, listUsersQuery)
	if err != nil {
		return users, fmt.Errorf("Failed to list users", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *repository) queryRowsWithFilter(ctx context.Context, query string, filters model.InvoiceFilter) (string, error) {

	return "", nil
}
