package security

import (
	"net/http"

	"github.com/invoice-service/spec"

	"github.com/golang-jwt/jwt/v4"
)

type jwtResponse struct {
	Role spec.Role
}

type JWT interface {
	GenerateJWT() (TokenDetails, error)
	VerifyJWT(w http.ResponseWriter, r *http.Request) (jwtResponse, error)
}

type CustomClaims struct {
	Id         int  `json:"id"`
	Authorized bool `json:"authorized"`
	Role       int  `json:"role"`
	jwt.StandardClaims
}

type AceessTokenClaims struct {
}

type RefreshTokenClaims struct {
}

type TokenDetails struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string
	AceessTokenClaims  AceessTokenClaims
	RefreshTokenClaims RefreshTokenClaims
	AccessTokenExpiry  int64
	RefreshTokenExpiry int64
}

func GetJWTClaims() jwt.Claims {
	return &CustomClaims{}
}
