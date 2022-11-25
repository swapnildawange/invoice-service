package inithandler

import (
	"database/sql"

	userRepository "github.com/invoice-service/user/repository"

	invoiceRepository "github.com/invoice-service/invoice/repository"
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
