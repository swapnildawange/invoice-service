package repository

import (
	"context"
	"invoice_service/spec"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_CreateUser(t *testing.T) {
	var (
		ctx           = context.TODO()
		createUserReq = spec.CreateUserRequest{
			Email: "email",
		}
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
				s.ExpectQuery("select COUNT\\(id\\) from auth where email = \\$1").WithArgs("email").WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
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
