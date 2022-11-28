package user

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"
	"github.com/invoice-service/security"
	"github.com/invoice-service/spec"
	"github.com/invoice-service/user/mocks"
)

type blMocks struct {
	logger log.Logger
	repo   *mocks.MockRepository
}

func initContext(id int, role spec.Role) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, gokitjwt.JWTClaimsContextKey, &security.CustomClaims{
		Id:         id,
		Authorized: true,
		Role:       int(role),
	})
	return ctx
}

func GetUserBL(t *testing.T) (bm blMocks, bl BL) {
	ct := gomock.NewController(t)
	var (
		logger = log.NewLogfmtLogger(os.Stderr)
		repo   = mocks.NewMockRepository(ct)
	)
	bm = blMocks{
		logger: logger,
		repo:   repo,
	}

	bl = NewBL(bm.logger, repo)
	return
}

func Test_makeCreateUser(t *testing.T) {
	var (
		ctx           = initContext(1, spec.RoleAdmin)
		createUserReq = spec.CreateUserRequest{
			Email:     defaultEmail,
			FirstName: defaultFirstName,
			LastName:  defaultLastName,
			Role:      defaultRole,
		}
		createUserRes = spec.User{
			Email:     defaultEmail,
			FirstName: defaultFirstName,
			LastName:  defaultLastName,
			Role:      defaultRole,
		}
		err = errors.New("some error")
	)

	type args struct {
		ctx     context.Context
		request spec.CreateUserRequest
	}
	tests := []struct {
		name        string
		args        args
		prepareTest func(*mocks.MockBL)
		want        interface{}
		wantErr     bool
	}{
		// Positive
		{
			name: "Positive",
			args: args{
				ctx:     ctx,
				request: createUserReq,
			},
			prepareTest: func(bm *mocks.MockBL) {
				bm.EXPECT().CreateUser(ctx, createUserReq).Return(createUserRes, nil)
			},
			want:    createUserRes,
			wantErr: false,
		},
		// Negative | user not authorized
		{
			name: "Negative | user not authorized",
			args: args{
				ctx:     initContext(1, spec.RoleUser),
				request: createUserReq,
			},
			prepareTest: func(bm *mocks.MockBL) {
			},
			want:    nil,
			wantErr: true,
		},
		// Negative | Failed to create user
		{
			name: "Negative | Failed to create user",
			args: args{
				ctx:     ctx,
				request: createUserReq,
			},
			prepareTest: func(bm *mocks.MockBL) {
				bm.EXPECT().CreateUser(ctx, createUserReq).Return(spec.User{}, err)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ct     = gomock.NewController(t)
				logger = log.NewLogfmtLogger(os.Stderr)
				bm     *mocks.MockBL
			)
			bm = mocks.NewMockBL(ct)
			tt.prepareTest(bm)

			fun := makeCreateUser(logger, bm)
			res, err := fun(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeCreateUser() error = %v", err)
				return
			}
			if tt.wantErr && res != nil {
				t.Errorf("makeCreateUser() = res should be nil in case of error")
			}
			if !tt.wantErr && res == nil {
				t.Errorf("makeCreateUser() = res should not be nil in case of success")
			}
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("makeCreateUser() = %v, want %v", res, tt.want)
			}
		})
	}
}
