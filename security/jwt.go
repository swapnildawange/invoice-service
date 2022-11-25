package security

import (
	"fmt"
	"net/http"
	"time"

	"github.com/invoice-service/spec"
	"github.com/invoice-service/svcerror"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
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

// func GenerateJWT(key string, userId, role int) (string, error) {
// 	var (
// 		token  *jwt.Token
// 		claims jwt.MapClaims
// 	)

// 	token = jwt.New(jwt.SigningMethodHS256)
// 	claims = token.Claims.(jwt.MapClaims)

// 	claims["exp"] = time.Now().Add(time.Minute * 120).Unix()
// 	claims["authorized"] = true
// 	claims["role"] = role
// 	claims["id"] = userId

// 	tokenString, err := token.SignedString([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }

func GenerateJWT(id, role int) (TokenDetails, error) {
	var (
		token        *jwt.Token
		atClaims     = jwt.MapClaims{}
		rtClaims     = jwt.MapClaims{}
		tokenDetails TokenDetails
		key          []byte
		err          error
	)
	// create access token
	tokenDetails.AccessTokenExpiry = time.Now().Add(time.Minute * 10).Unix()

	atClaims["exp"] = tokenDetails.AccessTokenExpiry
	atClaims["authorized"] = true
	atClaims["role"] = role
	atClaims["id"] = id

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	key = []byte(viper.GetString("ACCESS_TOKEN_SECRET"))
	tokenDetails.AccessToken, err = token.SignedString(key)
	if err != nil {
		return tokenDetails, svcerror.ErrFailedToGenerateAccessToken
	}

	// create refresh token
	tokenDetails.RefreshTokenExpiry = time.Now().Add(24 * time.Hour).Unix()
	rtClaims["exp"] = tokenDetails.RefreshTokenExpiry
	rtClaims["id"] = id

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	key = []byte(viper.GetString("REFRESH_TOKEN_SECRET"))
	tokenDetails.RefreshToken, err = token.SignedString(key)
	if err != nil {
		return tokenDetails, svcerror.ErrFailedToGenerateRefreshToken
	}
	return tokenDetails, nil
}

func VerifyJWT(endpointHandler *httptransport.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			return
		}
		tokenString := r.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("There was an error in parsing token")
			}
			return viper.Get("JWTSECRET"), nil
		})
		if err != nil {
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				r.Header.Set("Role", "admin")
				endpointHandler.ServeHTTP(w, r)
				return
			} else if claims["role"] == "user" {
				r.Header.Set("Role", "user")
				endpointHandler.ServeHTTP(w, r)
				return
			}
		}
	}
}
