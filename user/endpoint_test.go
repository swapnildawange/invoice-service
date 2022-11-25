package user

// import (
// 	"context"
// 	"os"
// 	"testing"

// 	"github.com/invoice-service/security"
// 	"github.com/invoice-service/user/mocks"

// 	"github.com/go-kit/kit/auth/jwt"
// 	"github.com/go-kit/log"
// 	"github.com/golang/mock/gomock"
// )

// type blMocks struct {
// 	logger log.Logger
// 	repo   *mocks.MockRepository
// }

// func initContext() context.Context {
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, jwt.JWTClaimsContextKey, security.CustomClaims{
// 		Id:         1,
// 		Authorized: true,
// 		Role:       1,
// 	})
// 	return ctx
// }
// func GetUserBL(t *testing.T) (bm blMocks, bl BL) {
// 	ct := gomock.NewController(t)
// 	var (
// 		logger = log.NewLogfmtLogger(os.Stderr)
// 		repo   = mocks.NewMockRepository(ct)
// 	)
// 	bm = blMocks{
// 		logger: logger,
// 		repo:   repo,
// 	}

// 	bl = NewBL(bm.logger, repo)
// 	return
// }

// // func Test_makeCreateUser(t *testing.T) {

// // 	var (
// // 		ctx           = initContext()
// // 		createUserReq = spec.CreateUserRequest{
// // 			Email:     defaultEmail,
// // 			FirstName: defaultFirstName,
// // 			LastName:  defaultLastName,
// // 			Role:      defaultRole,
// // 		}
// // 		createUserRes = spec.User{
// // 			Email:     defaultEmail,
// // 			FirstName: defaultFirstName,
// // 			LastName:  defaultLastName,
// // 			Role:      defaultRole,
// // 		}
// // 	)

// // 	type args struct {
// // 		ctx     context.Context
// // 		request spec.CreateUserRequest
// // 	}
// // 	tests := []struct {
// // 		name        string
// // 		args        args
// // 		prepareTest func(*mocks.MockBL)
// // 		want        spec.User
// // 		wantErr     bool
// // 	}{
// // 		// Positive
// // 		{
// // 			name: "Positive",
// // 			args: args{
// // 				ctx:     ctx,
// // 				request: createUserReq,
// // 			},
// // 			prepareTest: func(bm *mocks.MockBL) {
// // 				bm.EXPECT().CreateUser(ctx, createUserReq).Return(createUserRes, nil)
// // 			},
// // 		},
// // 	}
// // 	for _, tt := range tests {
// // 		t.Run(tt.name, func(t *testing.T) {
// // 			var (
// // 				ct     = gomock.NewController(t)
// // 				logger = log.NewLogfmtLogger(os.Stderr)
// // 				bm     *mocks.MockBL
// // 				bl     BL
// // 			)
// // 			bm = mocks.NewMockBL(ct)
// // 			_, bl = GetUserBL(t)
// // 			tt.prepareTest(bm)

// // 			fun := makeCreateUser(logger, bl)
// // 			res, err := fun(tt.args.ctx, tt.args.request)
// // 			if err != nil != tt.wantErr {
// // 				t.Errorf("makeCreateUser() error %v", err)
// // 			}
// // 			if tt.wantErr && err == nil {
// // 				t.Errorf("makeCreateUser() got %v want %v", err, tt.wantErr)
// // 			}
// // 			if !tt.wantErr && res == nil {
// // 				t.Errorf("makeCreateUser() = res should not be nil in case of error")
// // 			}
// // 			if !reflect.DeepEqual(res, tt.want) {
// // 				t.Errorf("makeCreateUser() = %v, want %v", res, tt.want)
// // 			}
// // 		})
// // 	}
// // }
