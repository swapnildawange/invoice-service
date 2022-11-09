package model

import "time"

const (
	PageSize = 5
)

type Config struct {
	WebPort int
}

// payment status
type PaymentStatus int

const (
	Initiate   PaymentStatus = 1
	InProgress PaymentStatus = 2
	Success    PaymentStatus = 3
	Failed     PaymentStatus = 4
	Retry      PaymentStatus = 5
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var MapRoleToName = map[int]Role{
	1: RoleAdmin,
	2: RoleUser,
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

type Invoice struct {
	Id            string        `json:"invoice_id"`
	UserId        int           `json:"user_id"`
	Paid          float64       `json:"paid"`
	AdminId       int           `json:"admin_id"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type CreateInvoiceRequest struct {
	Id            string        `json:"invoice_id"`
	UserId        int           `json:"user_id"`
	Paid          float64       `json:"paid"`
	AdminId       int           `json:"admin_id"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type GetInvoiceRequest struct {
	Id string `json:"id"` // invoice_id
}

type UpdateInvoiceRequest struct {
	Invoice
}

type CreateUserRequest struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Role      int       `json:"role"`
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
	Token string `json:"token"`
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
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Page      int    `json:"page"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type InvoiceFilter struct {
	Id            string  `json:"id"`
	UserId        int     `json:"user_id"`
	AdminId       int     `json:"admin_id"`
	Paid          float64 `json:"paid"`
	PaymentStatus int     `json:"payment_status"`
	Page          int     `json:"page"`
	SortBy        string  `json:"sort_by"`
	SortOrder     string  `json:"sort_order"`
}

type DeleteUserReq struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type EditUserRequest struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
