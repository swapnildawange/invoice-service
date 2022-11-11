package user

import (
	"fmt"
	"invoice_service/model"
	"invoice_service/security"
	"reflect"
	"time"
)

type UserMatcher struct {
	user model.CreateUserRequest
}

func (e UserMatcher) Matches(x model.CreateUserRequest) bool {

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

// func TestBL_CreateUser(t *testing.T) {
// 	var (
// 		ctx           = context.Background()
// 		userId        = defaultUserId
// 		createUserReq = model.CreateUserRequest{
// 			Email:     defaultEmail,
// 			FirstName: defaultFirstName,
// 			LastName:  defaultLastName,
// 			Role:      defaultRole,
// 		}
// 		createUserRes = model.User{
// 			Id:        userId,
// 			Email:     defaultEmail,
// 			FirstName: defaultFirstName,
// 			LastName:  defaultLastName,
// 			Role:      defaultRole,
// 		}
// 	)
// 	type args struct {
// 		ctx           context.Context
// 		createUserReq model.CreateUserRequest
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		prepareTest func(args, blMocks)
// 		want        model.User
// 		wantErr     bool
// 	}{
// 		{
// 			name: "Positive",
// 			args: args{
// 				ctx:           ctx,
// 				createUserReq: createUserReq,
// 			},
// 			want:    createUserRes,
// 			wantErr: false,
// 			prepareTest: func(a args, bm blMocks) {
// 				bm.repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(userId, nil)
// 				bm.repo.EXPECT().GetUser(ctx, gomock.Any()).Return(createUserRes, nil)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			bm, bl := GetUserBL(t)
// 			tt.prepareTest(tt.args, bm)
// 			got, err := bl.CreateUser(tt.args.ctx, tt.args.createUserReq)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("bl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("bl.CreateUser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
