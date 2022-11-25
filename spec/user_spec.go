package spec

import (
	"time"
)

type Config struct {
	WebPort int
}

var MapRoleToName = map[Role]string{
	1: "admin",
	2: "user",
}

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// when the admin craete the invoice it already  have the admin details
// while getting the invoice user already have the user details
// so no need to add the user details and admin details in invoice only
// respective ids need to add into database

type CreateUserRequest struct {
	Id        int       `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name"  validate:"required"`
	Password  string    `json:"password"  validate:"required,passwd"`
	Role      int       `json:"role"  validate:"required,oneof=1 2"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// request
type AuthRequest struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// response
type AuthResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  int    `json:"role"`
	Mesg  string `json:"mesg,omitempty"`
	Err   error  `json:"err,omitempty"`
}

type SignUpRequest struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type SignUpResponse struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ListUsersRequest struct {
	Page      int    `json:"page"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserFilter struct {
	Id        int    `json:"id" `
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Page      int    `json:"page" validate:"gt=0"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order" validate:"oneof=ASC DESC"`
}

type DeleteUserReq struct {
	Id    int    `json:"id"`
	Email string `json:"email" validate:"email"`
}

type EditUserRequest struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
