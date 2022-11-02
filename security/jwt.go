package jwt

import (
	"fmt"
	"invoice_service/model"
	"net/http"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/golang-jwt/jwt"
)

type jwtResponse struct {
	Role model.Role
}

type JWT interface {
	GenerateJWT() (string, error)
	VerifyJWT(w http.ResponseWriter, r *http.Request) (jwtResponse, error)
}

func GenerateJWT() (string, error) {
	var (
		token  *jwt.Token
		claims jwt.MapClaims
	)

	token = jwt.New(jwt.SigningMethodEdDSA)
	claims = token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["role"] = "admin"
	claims["email"] = "email"
	claims["id"] = "user_id"

	tokenString, err := token.SignedString("mysecretkey")
	if err != nil {
		return "", nil
	}

	return tokenString, nil

}

// func VerifyJWT(r *http.Request) (jwtResponse, error) {
// 	if r.Header["Token"] == nil {
// 		return jwtResponse{}, fmt.Errorf("No token found")
// 	}
// 	tokenString := r.Header["Token"][0]
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
// 			return nil, fmt.Errorf("There was an error in parsing token")

// 		}
// 		return "mysecretkey", nil
// 	})
// 	if err != nil {
// 		return jwtResponse{}, err
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if claims["role"] == "admin" {
// 			return jwtResponse{
// 				Role: model.RoleAdmin,
// 			}, nil
// 		} else if claims["role"] == "user" {
// 			return jwtResponse{
// 				Role: model.RoleUser,
// 			}, nil
// 		}
// 	}
// 	return jwtResponse{}, nil
// }

// func VerifyJWT(endpointHandler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Header["Token"] == nil {
// 			return
// 		}
// 		tokenString := r.Header["Token"][0]
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
// 				return nil, fmt.Errorf("There was an error in parsing token")

// 			}
// 			return "mysecretkey", nil
// 		})
// 		if err != nil {
// 			return
// 		}
// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			if claims["role"] == "admin" {
// 				r.Header.Set("Role", "admin")
// 				endpointHandler.ServeHTTP(w, r)
// 				return
// 			} else if claims["role"] == "user" {
// 				r.Header.Set("Role", "user")
// 				endpointHandler.ServeHTTP(w, r)
// 				return
// 			}
// 		}
// 	}
// }

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
			return "mysecretkey", nil
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