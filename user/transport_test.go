package user

// import (
// 	"context"
// 	"io"
// 	"net/http"
// 	"reflect"
// 	"strings"
// 	"testing"

// 	"github.com/invoice-service/spec"
// 	"github.com/invoice-service/utils"
// )

// var (
// 	defaultEmail     = "test@email.com"
// 	defaultFirstName = "testFirstName"
// 	defaultLastName  = "testLastName"
// 	defaultPassword  = "testPassword"
// 	defaultRole      = 2
// 	defaultUserId    = 100
// )

// func Test_decodeCreateUserRequest(t *testing.T) {
// 	var platformJWT = utils.InitPlatformJWT(t, defaultUserId, 1)
// 	validRequestBody := io.NopCloser(strings.NewReader(`
// 					{
// 						"email":"test@email.com",
//     					"first_name":"testFirstName",
//     					"last_name":"testLastName",
//     					"password":"testPassword",
//     					"role":2
// 					}`))

// 	mockValidRequest, _ := http.NewRequest("POST", CreateUserRequestPath, validRequestBody)
// 	mockValidRequest.Header.Set("Authorization", platformJWT)

// 	inValidRequestBody := io.NopCloser(strings.NewReader(`
// 					{
// 						"email":"test@email.com,
//     					"first_name":"testFirstName",
//     					"last_name":"testLastName",
//     					"password":"testPassword",
//     					"role":2
// 					}`))

// 	mockInvalidRequest, _ := http.NewRequest("POST", CreateUserRequestPath, inValidRequestBody)
// 	mockInvalidRequest.Header.Set("Authorization", platformJWT)

// 	type args struct {
// 		ctx context.Context
// 		req *http.Request
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    interface{}
// 		wantErr bool
// 	}{
// 		// Positive
// 		{
// 			name: "Positive",
// 			args: args{
// 				ctx: context.Background(),
// 				req: mockValidRequest,
// 			},
// 			want: spec.CreateUserRequest{
// 				Email:     defaultEmail,
// 				FirstName: defaultFirstName,
// 				LastName:  defaultLastName,
// 				Password:  defaultPassword,
// 				Role:      2,
// 			},
// 			wantErr: false,
// 		},
// 		// Negative | Invalid request body
// 		{
// 			name: "Negative | Invalid request body",
// 			args: args{
// 				ctx: context.Background(),
// 				req: mockValidRequest,
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := decodeCreateUserRequest(tt.args.ctx, tt.args.req)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("decodeCreateUserRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("decodeCreateUserRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
