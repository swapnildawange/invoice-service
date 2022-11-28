package user

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
)

type UserMatcher struct {
	user spec.CreateUserRequest
}

func (e UserMatcher) Matches(x spec.CreateUserRequest) bool {

	e.user.CreatedAt = time.Now()
	e.user.UpdatedAt = time.Now()
	hashedPassword, err := security.HashPassword(x.Password)
	if err != nil {
		return false
	}
	e.user.Password = hashedPassword
	return reflect.DeepEqual(e.user, x)
}

func (e UserMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.user)
}

func TestBL_CreateUser(t *testing.T) {
	var (
		ctx           = context.TODO()
		userId        = defaultUserId
		createUserReq = spec.CreateUserRequest{
			Email:     defaultEmail,
			FirstName: defaultFirstName,
			LastName:  defaultLastName,
			Role:      defaultRole,
		}
		createUserRes = spec.User{
			Id:        userId,
			Email:     defaultEmail,
			FirstName: defaultFirstName,
			LastName:  defaultLastName,
			Role:      defaultRole,
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
		prepareTest func(args, blMocks)
		want        spec.User
		wantErr     bool
	}{
		// positive
		{
			name: "Positive",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			want:    createUserRes,
			wantErr: false,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Create(ctx, gomock.Any()).Return(userId, nil)
				bm.repo.EXPECT().Get(ctx, gomock.Any()).Return(createUserRes, nil)
			},
		},
		// Negative | Failed to Create user in repo
		{
			name: "Negative | Failed to Create user in repo",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			want:    spec.User{},
			wantErr: true,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Create(ctx, gomock.Any()).Return(userId, err)
			},
		},
		// Negative | Failed to get updated user
		{
			name: "Negative | Failed to get updated user",
			args: args{
				ctx:           ctx,
				createUserReq: createUserReq,
			},
			want:    spec.User{},
			wantErr: true,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Create(ctx, gomock.Any()).Return(userId, nil)
				bm.repo.EXPECT().Get(ctx, gomock.Any()).Return(spec.User{}, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm, bl := GetUserBL(t)
			tt.prepareTest(tt.args, bm)
			got, err := bl.CreateUser(tt.args.ctx, tt.args.createUserReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("bl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bl.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBL_GetUser(t *testing.T) {

	var (
		ctx  = context.TODO()
		user = spec.User{}
		err  = errors.New("some error")
	)
	type args struct {
		ctx    context.Context
		userId int
	}
	tests := []struct {
		name        string
		args        args
		want        spec.User
		prepareTest func(args, blMocks)
		wantErr     bool
	}{
		// Positive
		{
			name: "Positive",
			args: args{
				ctx:    ctx,
				userId: defaultUserId,
			},
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Get(ctx, a.userId).Return(user, nil)
			},
			wantErr: false,
			want:    user,
		},
		// Negative
		{
			name: "Negative",
			args: args{
				ctx:    ctx,
				userId: defaultUserId,
			},
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Get(ctx, a.userId).Return(user, err)
			},
			wantErr: true,
			want:    user,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm, bl := GetUserBL(t)
			tt.prepareTest(tt.args, bm)
			got, err := bl.GetUser(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("bl.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bl.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBL_ListUsers(t *testing.T) {
	type args struct {
		ctx            context.Context
		listUserFilter spec.UserFilter
	}
	var (
		resp = []spec.User{
			{
				Id:        defaultUserId,
				Email:     defaultEmail,
				FirstName: defaultFirstName,
				LastName:  defaultLastName,
				Role:      defaultRole,
				CreatedAt: defaultTime,
			},
		}
		err       = errors.New("some error")
		ctx       = context.TODO()
		inputArgs = args{
			ctx: ctx,
			listUserFilter: spec.UserFilter{
				Id: 1,
			},
		}
	)
	tests := []struct {
		name        string
		prepareTest func(args, blMocks)
		args        args
		want        []spec.User
		wantErr     bool
	}{
		// Positive
		{
			name:    "Positive",
			args:    inputArgs,
			want:    resp,
			wantErr: false,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().List(a.ctx, a.listUserFilter).Return(resp, nil)
			},
		},
		// Negative | List user failed in repo
		{
			name:    "Negative | List user failed in repo",
			args:    inputArgs,
			want:    resp,
			wantErr: true,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().List(a.ctx, a.listUserFilter).Return(resp, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm, bl := GetUserBL(t)
			tt.prepareTest(tt.args, bm)
			got, err := bl.ListUsers(tt.args.ctx, tt.args.listUserFilter)
			if (err != nil) != tt.wantErr {
				t.Errorf("bl.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bl.ListUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBL_DeleteUser(t *testing.T) {

	type args struct {
		ctx           context.Context
		deleteUserReq spec.DeleteUserReq
	}
	var (
		inputArgs = args{
			ctx: context.TODO(),
			deleteUserReq: spec.DeleteUserReq{
				Id: defaultUserId,
			},
		}
		err = errors.New("some error")
	)
	tests := []struct {
		name        string
		args        args
		want        int
		wantErr     bool
		prepareTest func(args, blMocks)
	}{
		// Positive
		{
			name:    "Positive",
			args:    inputArgs,
			want:    defaultUserId,
			wantErr: false,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Delete(a.ctx, a.deleteUserReq).Return(defaultUserId, nil)
			},
		},
		// Negative | DeleteUser failed in repo
		{
			name:    "Negative | DeleteUser failed in repo",
			args:    inputArgs,
			want:    defaultUserId,
			wantErr: true,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Delete(a.ctx, a.deleteUserReq).Return(defaultUserId, err)
			},
		},
	}
	for _, tt := range tests {
		bm, bl := GetUserBL(t)
		tt.prepareTest(tt.args, bm)
		t.Run(tt.name, func(t *testing.T) {
			got, err := bl.DeleteUser(tt.args.ctx, tt.args.deleteUserReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("bl.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bl.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBL_EditUser(t *testing.T) {
	type args struct {
		ctx         context.Context
		editUserReq spec.EditUserRequest
	}
	var (
		inputArgs = args{
			ctx: context.TODO(),
			editUserReq: spec.EditUserRequest{
				Id:        defaultUserId,
				Email:     defaultEmail,
				FirstName: defaultFirstName,
				LastName:  defaultLastName,
			},
		}
		resp = spec.User{
			Id:        defaultUserId,
			Email:     defaultEmail,
			FirstName: defaultFirstName,
			LastName:  defaultLastName,
		}
		err = errors.New("some error")
	)
	tests := []struct {
		name        string
		args        args
		want        spec.User
		wantErr     bool
		prepareTest func(args, blMocks)
	}{
		// Positive
		{
			name:    "Positive",
			args:    inputArgs,
			want:    resp,
			wantErr: false,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Edit(a.ctx, a.editUserReq).Return(resp, nil)
			},
		},
		// Negative | EditUser failed in repo
		{
			name:    "Negative | EditUser failed in repo",
			args:    inputArgs,
			want:    resp,
			wantErr: true,
			prepareTest: func(a args, bm blMocks) {
				bm.repo.EXPECT().Edit(a.ctx, a.editUserReq).Return(resp, err)
			},
		},
	}
	for _, tt := range tests {
		bm, bl := GetUserBL(t)
		tt.prepareTest(tt.args, bm)
		t.Run(tt.name, func(t *testing.T) {
			got, err := bl.EditUser(tt.args.ctx, tt.args.editUserReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("bl.EditUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bl.EditUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
