package user

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/invoice-service/spec"
)

var (
	defaultEmail     = "test@email.com"
	defaultFirstName = "testFirstName"
	defaultLastName  = "testLastName"
	defaultPassword  = "testPassword"
	defaultRole      = 2
	defaultUserId    = 100
	defaultTime      = time.Now()
)

func Test_decodeCreateUserRequest(t *testing.T) {
	var platformJWT = ""
	// utils.InitPlatformJWT(t, defaultUserId, 1)
	validRequestBody := io.NopCloser(strings.NewReader(`
					{
						"email":"test@email.com",
    					"first_name":"testFirstName",
    					"last_name":"testLastName",
    					"password":"testPassword",
    					"role":2
					}`))

	mockValidRequest, _ := http.NewRequest("POST", spec.CreateUserRequestPath, validRequestBody)
	mockValidRequest.Header.Set("Authorization", platformJWT)

	inValidRequestBody := io.NopCloser(strings.NewReader(`
					{
						"email":"test@email.com,
    					"first_name":"testFirstName",
    					"last_name":"testLastName",
    					"password":"testPassword",
    					"role":2
					}`))

	mockInvalidRequest, _ := http.NewRequest("POST", spec.CreateUserRequestPath, inValidRequestBody)
	mockInvalidRequest.Header.Set("Authorization", platformJWT)

	mockValidRequest.Header.Set("Authorization", platformJWT)

	inValidRequestBody1 := io.NopCloser(strings.NewReader(`
					{
						"email":"invalidEmail",
    					"first_name":"testFirstName",
    					"last_name":"testLastName",
    					"password":"testPassword",
    					"role":2
					}`))

	mockInvalidRequest1, _ := http.NewRequest("POST", spec.CreateUserRequestPath, inValidRequestBody1)
	mockInvalidRequest1.Header.Set("Authorization", platformJWT)

	type args struct {
		ctx context.Context
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// Positive
		{
			name: "Positive",
			args: args{
				ctx: context.Background(),
				req: mockValidRequest,
			},
			want: spec.CreateUserRequest{
				Email:     defaultEmail,
				FirstName: defaultFirstName,
				LastName:  defaultLastName,
				Password:  defaultPassword,
				Role:      2,
			},
			wantErr: false,
		},
		// Negative | Invalid request body
		{
			name: "Negative | Invalid request body",
			args: args{
				ctx: context.Background(),
				req: mockValidRequest,
			},
			want:    nil,
			wantErr: true,
		},
		// Negative | invalid email
		{
			name: "Negative | invalid email",
			args: args{
				ctx: context.Background(),
				req: mockInvalidRequest1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeCreateUserRequest(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCreateUserRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeCreateUserRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
