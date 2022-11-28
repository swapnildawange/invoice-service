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

//go:generate  mockgen -destination=../mocks/repository.mock.go -package=mocks github.com/invoice-service/user/repository Repository
type Repository interface {
	Create(ctx context.Context, createUserReq spec.CreateUserRequest) (int, error)
	List(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error)
	GetUserFromAuth(ctx context.Context, email string) (userId int, hashedPassword string, err error)
	Get(ctx context.Context, userId int) (spec.User, error)
	Delete(ctx context.Context, deleteUserReq spec.DeleteUserReq) (int, error)
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

	row := repo.db.QueryRowContext(ctx, `SELECT user_id,password FROM auth WHERE email = $1`, email)
	if row.Err() != nil {

		return userId, hashedPassword, err
	}
	if err = row.Scan(&userId, &hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return userId, hashedPassword, svcerror.ErrUserNotFound
		}
		return userId, hashedPassword, row.Err()
	}
	return userId, hashedPassword, nil
}

func (repo *repository) Get(ctx context.Context, userId int) (spec.User, error) {
	var (
		user spec.User
		err  error
		row  *sql.Row
	)
	getUserQuery := `SELECT id,first_name,last_name,role,created_at,updated_at FROM users WHERE id=$1;`
	row = repo.db.QueryRowContext(ctx, getUserQuery, userId)
	if row.Err() != nil {
		return user, err
	}

	err = row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, svcerror.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}

func (repo *repository) Create(ctx context.Context, createUserReq spec.CreateUserRequest) (int, error) {
	var (
		err        error
		row        *sql.Row
		tx         *sql.Tx
		userId     int
		emailCount int
	)
	ctx, cancel := context.WithTimeout(ctx, spec.Timeout*time.Second)
	defer cancel()

	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return userId, fmt.Errorf("failed to begin transaction %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// check if email is already present
	checkEmailQuery := `SELECT COUNT(id) FROM auth WHERE email = $1`

	row = tx.QueryRowContext(ctx, checkEmailQuery, createUserReq.Email)
	if row.Err() != nil {
		return userId, err
	}

	if err := row.Scan(&emailCount); err != nil {
		return userId, err
	}

	if emailCount != 0 {
		return userId, fmt.Errorf("email is already present")
	}

	insertIntoUsersQuery := `insert into users(first_name,last_name,role,created_at,updated_at) values ($1,$2,$3,$4,$5) RETURNING id`

	row = tx.QueryRowContext(ctx, insertIntoUsersQuery, createUserReq.FirstName, createUserReq.LastName, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt)
	if row.Err() != nil {
		return userId, row.Err()
	}

	if err = row.Scan(&userId); err != nil {
		return userId, err
	}

	insertIntoAuthQuery := `insert into auth(user_id,email,password) values($1,$2,$3);`
	_, err = tx.ExecContext(ctx, insertIntoAuthQuery, userId, createUserReq.Email, createUserReq.Password)
	if err != nil {
		return userId, err
	}

	err = tx.Commit()
	if err != nil {
		return userId, fmt.Errorf("failed to commit db transaction %v", err)
	}
	return userId, nil
}

func (repo *repository) List(ctx context.Context, listUserFilter spec.UserFilter) ([]spec.User, error) {
	var (
		users = make([]spec.User, 0)
		err   error
		rows  *sql.Rows
	)
	listUsersQuery := `SELECT id,first_name,last_name,role,created_at,updated_at FROM users `

	listUsersQuery, filterValues := repo.queryUsersWithFilter(ctx, listUsersQuery, listUserFilter)

	rows, err = repo.db.QueryContext(ctx, listUsersQuery, filterValues...)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user spec.User
		if err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			if err == sql.ErrNoRows {
				return users, nil
			}
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

func (repo *repository) Delete(ctx context.Context, deleteUserReq spec.DeleteUserReq) (int, error) {
	var (
		tx     *sql.Tx
		err    error
		userId int
	)

	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return userId, fmt.Errorf("failed to begin transaction %v", err.Error())
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if deleteUserReq.Email != "" {
		userId, _, err = repo.GetUserFromAuth(ctx, deleteUserReq.Email)
		if err != nil {
			return userId, err
		}
		if deleteUserReq.Id != -1 && deleteUserReq.Id != userId {
			return userId, svcerror.ErrUserNotFound
		}
		deleteUserReq.Id = userId
	} else {
		userId = deleteUserReq.Id
	}
	_, err = repo.Get(ctx, userId)
	if err != nil {
		return userId, err
	}
	// delete from users
	usersQuery := `DELETE FROM users where id = $1`
	_, err = tx.ExecContext(ctx, usersQuery, userId)
	if err != nil {
		return userId, fmt.Errorf("failed to execute delete query for users table %v", err)
	}
	// delete from auth
	authQuery := `DELETE FROM auth where id = $1`
	_, err = tx.ExecContext(ctx, authQuery, userId)
	if err != nil {
		return userId, fmt.Errorf("failed to execute  delete query for auth query %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return userId, fmt.Errorf("failed to commit transaction %v", err)
	}

	return userId, nil
}

func (repo *repository) Edit(ctx context.Context, editUserReq spec.EditUserRequest) (spec.User, error) {
	var (
		user   spec.User
		values []interface{}
		tx     *sql.Tx
		err    error
		row    *sql.Row
	)
	tx, err = repo.db.BeginTx(ctx, nil)
	if err != nil {
		return user, fmt.Errorf("failed to begin db transaction %s", err.Error())
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// check if user exsists
	selectQuery := `select id,first_name,last_name,role,created_at,updated_at from users where id = $1`
	row = tx.QueryRowContext(ctx, selectQuery, editUserReq.Id)
	if row.Err() != nil {
		return user, err
	}
	err = row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, svcerror.ErrUserNotFound
		}
		return user, fmt.Errorf("failed to scan  user %s", err.Error())
	}

	updateQuery := `UPDATE users SET`

	if editUserReq.FirstName != "" {
		values = append(values, editUserReq.FirstName)
		updateQuery += ` first_name = $` + strconv.Itoa(len(values)) + ` ,`
		user.FirstName = editUserReq.FirstName
	}
	if editUserReq.LastName != "" {
		values = append(values, editUserReq.LastName)
		updateQuery += ` last_name = $` + strconv.Itoa(len(values)) + ` ,`
		user.LastName = editUserReq.LastName
	}

	values = append(values, time.Now())
	updateQuery += ` updated_at = $` + strconv.Itoa(len(values)) + ` ,`

	values = append(values, editUserReq.Id)
	updateQuery += ` WHERE id = $` + strconv.Itoa(len(values))

	updateQuery = strings.ReplaceAll(updateQuery, ", WHERE", "WHERE ")
	_, err = tx.ExecContext(ctx, updateQuery, values...)
	if err != nil {
		return user, fmt.Errorf("failed to execute update user query %s", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return user, fmt.Errorf("failed to commit db transaction %s", err.Error())
	}

	return user, nil
}
