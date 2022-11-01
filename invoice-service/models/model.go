package models

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

// email , hash , id 

type Admin struct {
	Id int
}

type User struct {
	Id int
}

type Invoice struct {
	User               User          `json:"user"`
	Id                 string        `json:"invoice_id"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	Paid               float64       `json:"paid"`
	PaymentInitiatedBy Admin         `json:"admin"`
	PaymentStatus      PaymentStatus `json:"status"`
}

type CreateInvoiceRequest struct {
	UserId int     `json:"user_id"`
	Paid   float64 `json:"paid"`
}

type GetInvoiceRequest struct {
	Id string `json:"id"`
}

type UpdateInvoiceRequest struct {
	Invoice Invoice `json:"invoice"`
}
