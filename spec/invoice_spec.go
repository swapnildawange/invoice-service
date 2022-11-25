package spec

import "time"

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
	UserId        int           `json:"user_id" validate:"required"`
	Paid          float64       `json:"paid" validate:"min=0"`
	AdminId       int           `json:"admin_id"`
	PaymentStatus PaymentStatus `json:"payment_status" validate:"required,oneof=1 2 3 4 5"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type GetInvoiceRequest struct {
	Id     string `json:"id"` // invoice_id
	UserId int    `json:"user_id"`
}

type UpdateInvoiceRequest struct {
	Id            string    `json:"invoice_id"`
	Paid          float64   `json:"paid"`
	PaymentStatus int       `json:"payment_status"`
	UpdatedAt     time.Time `json:"updated_at"`
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
