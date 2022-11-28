package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/invoice-service/spec"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_CreateUser(t *testing.T) {
	var (
		ctx           = context.TODO()
		userId        = 1
		createUserReq = spec.CreateUserRequest{
			Email: "email",
		}
		err = errors.New("some error")
	)
	type args struct {
		ctx           context.Context
		createUserReq spec.CreateUserRequest
	}
	tests := []struct {
		name        string
		args        args
		want        int
		wantErr     bool
		prepareTest func(args, sqlmock.Sqlmock)
	}{
		// Positive
		{
			name: "Positive",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin()
				rows := sqlmock.NewRows([]string{"COUNT(id)"})
				rows = rows.AddRow(0)
				s.ExpectQuery("SELECT COUNT\\(id\\) FROM auth WHERE email = \\$1").WithArgs("email").WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"id"})
				rows = rows.AddRow(userId)
				s.ExpectQuery("insert into users\\(first_name,last_name,role,created_at,updated_at\\) values \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id").WithArgs(createUserReq.FirstName, createUserReq.LastName, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt).WillReturnRows(rows)
				s.ExpectExec("insert into auth\\(user_id,email,password\\) values\\(\\$1,\\$2,\\$3\\);").WithArgs(userId, createUserReq.Email, createUserReq.Password).WillReturnResult(driver.ResultNoRows)
				s.ExpectCommit()
			},
			want:    userId,
			wantErr: false,
		},
		// Negative | Failed to begin transaction
		{
			name: "Negative | Failed to begin transaction",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin().WillReturnError(err)
			},
			want:    0,
			wantErr: true,
		},
		// Negative | Email is already present
		{
			name: "Negative | Email is already present",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin()
				rows := sqlmock.NewRows([]string{"COUNT(id)"})
				rows = rows.AddRow(1)
				s.ExpectQuery("SELECT COUNT\\(id\\) FROM auth WHERE email = \\$1").WithArgs("email").WillReturnRows(rows)
			},
			want:    0,
			wantErr: true,
		},
		// Negative | Failed to execute insert query
		{
			name: "Negative | Failed to execute insert query",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin()
				s.ExpectQuery("SELECT COUNT\\(id\\) FROM auth WHERE email = \\$1").WithArgs("email").WillReturnError(err)
			},
			want:    0,
			wantErr: false,
		},
		// Negative | Failed to insert into users table
		{
			name: "Negative | Failed to insert into users table",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin()
				rows := sqlmock.NewRows([]string{"COUNT(id)"})
				rows = rows.AddRow(0)
				s.ExpectQuery("SELECT COUNT\\(id\\) FROM auth WHERE email = \\$1").WithArgs("email").WillReturnRows(rows)
				s.ExpectQuery("insert into users\\(first_name,last_name,role,created_at,updated_at\\) values \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id").WithArgs(createUserReq.FirstName, createUserReq.LastName, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt).WillReturnError(err)
			},
			want:    0,
			wantErr: true,
		},
		// Negative | Failed to insert into auth
		{
			name: "Negative | Failed to insert into auth",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin()
				rows := sqlmock.NewRows([]string{"COUNT(id)"})
				rows = rows.AddRow(0)
				s.ExpectQuery("SELECT COUNT\\(id\\) FROM auth WHERE email = \\$1").WithArgs("email").WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"id"})
				rows = rows.AddRow(userId)
				s.ExpectQuery("insert into users\\(first_name,last_name,role,created_at,updated_at\\) values \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id").WithArgs(createUserReq.FirstName, createUserReq.LastName, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt).WillReturnRows(rows)
				s.ExpectExec("insert into auth\\(user_id,email,password\\) values\\(\\$1,\\$2,\\$3\\);").WithArgs(userId, createUserReq.Email, createUserReq.Password).WillReturnError(err)
				s.ExpectRollback()
			},
			want:    1,
			wantErr: true,
		},
		// Negative | commit failed
		{
			name: "Negative | commit failed",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			prepareTest: func(a args, s sqlmock.Sqlmock) {
				s.ExpectBegin()
				rows := sqlmock.NewRows([]string{"COUNT(id)"})
				rows = rows.AddRow(0)
				s.ExpectQuery("SELECT COUNT\\(id\\) FROM auth WHERE email = \\$1").WithArgs("email").WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"id"})
				rows = rows.AddRow(userId)
				s.ExpectQuery("insert into users\\(first_name,last_name,role,created_at,updated_at\\) values \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id").WithArgs(createUserReq.FirstName, createUserReq.LastName, createUserReq.Role, createUserReq.CreatedAt, createUserReq.UpdatedAt).WillReturnRows(rows)
				s.ExpectExec("insert into auth\\(user_id,email,password\\) values\\(\\$1,\\$2,\\$3\\);").WithArgs(userId, createUserReq.Email, createUserReq.Password).WillReturnResult(driver.ResultNoRows)
				s.ExpectCommit().WillReturnError(err)
			},
			want:    userId,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			repo := repository{
				db: db,
			}
			tt.prepareTest(tt.args, mock)
			got, err := repo.Create(tt.args.ctx, tt.args.createUserReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repository.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
