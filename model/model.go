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

// email , hash , id
type Admin struct {
	Id int
}

type User struct {
	Id int
}

// when the admin craete the invoice it already  have the admin details
// while getting the invoice user already have the user details
// so no need to add the user details and admin details in invoice only
// respective ids need to add into database

type Invoice struct {
	Id                 string        `json:"invoice_id"`
	UserId             int           `json:"user_id"`
	Paid               float64       `json:"paid"`
	PaymentInitiatedBy int           `json:"admin_id"`
	PaymentStatus      PaymentStatus `json:"status"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

type CreateInvoiceRequest struct {
	UserId int     `json:"user_id"`
	Paid   float64 `json:"paid"`
}

type GetInvoiceRequest struct {
	Id string `json:"id"` // invoice_id
}

type UpdateInvoiceRequest struct {
	Invoice Invoice `json:"invoice"`
}
