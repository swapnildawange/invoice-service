package user

import (
	"context"
	"fmt"

	"github.com/invoice-service/security"
	"github.com/invoice-service/svcerror"

	gokitjwt "github.com/go-kit/kit/auth/jwt"
)

func CheckIsAuthorized(ctx context.Context) error {
	JWTClaims, ok := ctx.Value(gokitjwt.JWTClaimsContextKey).(*security.CustomClaims)
	if !ok {
		return fmt.Errorf("invalid jwt token")
	}

	if JWTClaims.Role == 2 {
		return svcerror.ErrNotAuthorized
	}
	return nil
}
