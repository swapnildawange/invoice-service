package inithandler

import (
	"database/sql"
	invoiceRepository "invoice_service/invoice/repository"
	userRepository "invoice_service/user/repository"
)

type Repository struct {
	UserRepo    userRepository.Repository
	InvoiceRepo invoiceRepository.Repository
}

func InitRepository(db *sql.DB) Repository {

	// initate repository
	var userRepo = userRepository.NewRepository(db)
	var inviceRepo = invoiceRepository.NewRepository(db)
	return Repository{
		UserRepo:    userRepo,
		InvoiceRepo: inviceRepo,
	}
}
