package model

import "time"

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

type User struct {
	Id         int       `json:"id"`
	Email      string    `json:"email"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Password   string    `json:"password"`
	Role       Role      `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
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
	Invoice Invoice `json:"invoice"`
}

type CreateUserRequest struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
