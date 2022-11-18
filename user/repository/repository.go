package repository

import (
	"context"
	"database/sql"
	"fmt"
	"invoice_service/spec"
	"strconv"
	"strings"
	"time"
)

//go:generate  mockgen -destination=../mocks/repository.mock.go -package=mocks invoice_service/user/repository Repository
type Repository interface {
	Create(ctx context.Context, createUserReq spec.CreateUserRequest) (int, error)
	List(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error)
	GetUserFromAuth(ctx context.Context, email string) (userId int, hashedPassword string, err error)
	Get(ctx context.Context, userId int) (spec.User, error)
	Delete(ctx context.Context, deleteUserReq spec.DeleteUserReq) error
	Edit(ctx context.Context, editUserReq spec.EditUserRequest) (spec.User, error)
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

	row := repo.db.QueryRowContext(ctx, `select user_id,password from auth where email = $1`, email)
	if row.Err() != nil {
		return userId, hashedPassword, fmt.Errorf("failed to get User details from auth %s", row.Err().Error())
	}
	if err = row.Scan(&userId, &hashedPassword); err != nil {
		return userId, hashedPassword, fmt.Errorf("failed to scan user details %s", row.Err().Error())
	}
	return userId, hashedPassword, nil
}

func (repo *repository) Get(ctx context.Context, userId int) (spec.User, error) {
	var (
		user spec.User
		err  error
		row  *sql.Row
	)
	getUserQuery := `select id,first_name,last_name,role,created_at,updated_at from users where id=$1;`
	row = repo.db.QueryRowContext(ctx, getUserQuery, userId)
	if row.Err() != nil {
		return user, fmt.Errorf("failed to get user %s ", row.Err())
	}

	if err = row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return user, err
	}
	return user, nil
}

func (repo *repository) Create(ctx context.Context, createUserReq spec.CreateUserRequest) (int, error) {
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
			tx.Rollback()
		}
	}()

	if err != nil {
		return userId, fmt.Errorf("failed to begin transaction %w", err)
	}

	// check if email is already present
	checkEmailQuery := `select COUNT(id) from auth where email = $1`

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
		return userId, fmt.Errorf("email is already present %w", err)
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

	err = tx.Commit()
	if err != nil {
		return userId, fmt.Errorf("failed to commit db transaction %w", err)
	}
	return userId, nil
}

func (repo *repository) List(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error) {
	var (
		users = make([]spec.User, 0)
		err   error
		rows  *sql.Rows
	)
	listUsersQuery := `select id,first_name,last_name,role,created_at,updated_at from users `

	listUsersQuery, filterValues := repo.queryUsersWithFilter(ctx, listUsersQuery, listUserFilter)

	rows, err = repo.db.QueryContext(ctx, listUsersQuery, filterValues...)
	if err != nil {
		return users, fmt.Errorf("failed to list users %v", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var user spec.User
		if err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *repository) queryUsersWithFilter(ctx context.Context, query string, filter spec.UserFilter) (string, []interface{}) {
	var (
		filterValues []interface{}
		count        int
	)

	if filter.FirstName != "" {
		query += " AND first_name LIKE '%" + filter.FirstName + "%'"
		count += 1
	}

	if filter.LastName != "" {
		query += " AND last_name LIKE '%" + filter.LastName + "%'"
		count += 1
	}

	if filter.Id != 0 {
		filterValues = append(filterValues, filter.Id)
		count += 1
		query += ` AND id = $` + strconv.Itoa(len(filterValues))
	}

	query += fmt.Sprintf(` ORDER BY %s %s`, filter.SortBy, filter.SortOrder)

	filterValues = append(filterValues, spec.PageSize)
	query += ` LIMIT $` + strconv.Itoa(len(filterValues))

	filterValues = append(filterValues, (filter.Page-1)*spec.PageSize)
	query += ` OFFSET $` + strconv.Itoa(len(filterValues))

	if count >= 1 {
		query = strings.Replace(query, "AND", "WHERE", 1)
	}

	return query, filterValues
}

func (repo *repository) Delete(ctx context.Context, deleteUserReq spec.DeleteUserReq) error {
	var (
		tx     *sql.Tx
		err    error
		userId int
	)

	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction %v", err.Error())
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if deleteUserReq.Email != "" {
		userId, _, err = repo.GetUserFromAuth(ctx, deleteUserReq.Email)
		if err != nil {
			return fmt.Errorf("failed to get user from  auth table %v", err.Error())
		}
		if deleteUserReq.Id != -1 && deleteUserReq.Id != userId {
			return fmt.Errorf("user not found ")
		}
		deleteUserReq.Id = userId
	} else {
		userId = deleteUserReq.Id
	}
	// delete from users
	usersQuery := `DELETE FROM users where id = $1`
	_, err = tx.ExecContext(ctx, usersQuery, userId)
	if err != nil {
		return fmt.Errorf("failed to execute delete query for users table %w", err)
	}
	// delete from auth
	authQuery := `DELETE FROM auth where id = $1`
	_, err = tx.ExecContext(ctx, authQuery, userId)
	if err != nil {
		return fmt.Errorf("failed to execute  delete query for auth query %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction %w", err)
	}

	return nil
}

func (repo *repository) Edit(ctx context.Context, editUserReq spec.EditUserRequest) (spec.User, error) {
	var (
		newUser spec.User
		values  []interface{}
		tx      *sql.Tx
		err     error
		row     *sql.Row
	)
	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return newUser, fmt.Errorf("failed to begin db transaction %s", err.Error())
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateQuery := `UPDATE users SET`

	if editUserReq.FirstName != "" {
		values = append(values, editUserReq.FirstName)
		updateQuery += ` first_name = $` + strconv.Itoa(len(values)) + ` ,`
	}
	if editUserReq.LastName != "" {
		values = append(values, editUserReq.LastName)
		updateQuery += ` last_name = $` + strconv.Itoa(len(values)) + ` ,`
	}

	values = append(values, time.Now())
	updateQuery += ` updated_at = $` + strconv.Itoa(len(values)) + ` ,`

	values = append(values, editUserReq.Id)
	updateQuery += ` WHERE id = $` + strconv.Itoa(len(values))

	updateQuery = strings.ReplaceAll(updateQuery, ", WHERE", "WHERE ")

	_, err = tx.ExecContext(ctx, updateQuery, values...)
	if err != nil {
		return newUser, fmt.Errorf("failed to execute update user query %s", err.Error())
	}

	selectQuery := `select id,first_name,last_name,role,created_at,updated_at from users where id = $1`
	row = tx.QueryRowContext(ctx, selectQuery, editUserReq.Id)
	if row.Err() != nil {
		return newUser, fmt.Errorf("failed to get updated user %s", err.Error())
	}
	if err = row.Scan(&newUser.Id, &newUser.FirstName, &newUser.LastName, &newUser.Role, &newUser.CreatedAt, &newUser.UpdatedAt); err != nil {
		return newUser, fmt.Errorf("failed to scan updated user %s", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return newUser, fmt.Errorf("failed to commit db transaction %s", err.Error())
	}

	return newUser, nil
}
