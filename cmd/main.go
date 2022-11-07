package main

import (
	"context"
	"database/sql"
	"fmt"
	"invoice_service/invoice"
	invoiceRepository "invoice_service/invoice/repository"
	"invoice_service/security"
	"invoice_service/user"
	userRepository "invoice_service/user/repository"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	webPort  = ":8080"
	host     = "localhost"
	port     = 5432
	psqlUser = "postgres"
	password = "password"
	dbname   = "users"
)

func main() {
	s, _ := security.HashPassword("secret")
	fmt.Println(s)
	ctx := context.Background()

	// initiate logger
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "method=", log.DefaultCaller)

	// initate db
	db, err := InitDB(logger)
	if err != nil {
		panic("can't connect to Postgres")
	}
	logger.Log("The database is connected")

	// initate repository
	var userRepo = userRepository.NewRepository(db)
	var inviceRepo = invoiceRepository.NewRepository(db)
	// initiate invoice bl
	var invoiceBL = invoice.NewBL(logger, inviceRepo)
	// initiate user bl
	var userBL = user.NewBL(logger, userRepo)

	// initate invoice endpoints
	var invoiceEndpoints = invoice.NewEndpoints(logger, invoiceBL)
	// initiate user endpoints
	var userEndpoints = user.NewEndpoints(logger, userBL)
	// initiate router
	var router = mux.NewRouter()
	// initate handlers
	router = invoice.NewHTTPHandler(ctx, logger, router, invoiceEndpoints)

	router = user.NewHTTPHandler(ctx, logger, router, userEndpoints)
	// start the server
	logger.Log("Starting the server on port", webPort)
	err = http.ListenAndServe(webPort, router)
	if err != nil {
		logger.Log(err)
	}

}

// This function will make a connection to the database only once.
func InitDB(logger log.Logger) (*sql.DB, error) {
	var err error
	connStr := "postgres://postgres:password@localhost/invoicing?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// func openDB(logger log.Logger, dsn string) (*sql.DB, error) {
// 	var (
// 		db  *sql.DB
// 		err error
// 	)
// 	db, err = sql.Open("pgx", dsn)
// 	if err != nil {
// 		logger.Log("Failed to open Database", err.Error())
// 		return db, err
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

// func connectToDB(logger log.Logger, retries, delay int) *sql.DB {
// 	dsn := os.Getenv("DSN")
// 	for i := 0; i < retries; i++ {
// 		connection, err := openDB(logger, dsn)
// 		if err != nil {
// 			logger.Log("Postgres not ready yet", err.Error())
// 			time.Sleep(time.Duration(delay) * time.Second)
// 		} else {
// 			logger.Log("Successfully connected postgres DB")
// 			return connection
// 		}
// 	}
// 	return nil
// }
