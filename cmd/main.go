package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/invoice-service/invoice"
	"github.com/invoice-service/user"

	"github.com/invoice-service/cmd/inithandler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	// initate viper
	config, err := inithandler.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// initiate logger
	var logger = inithandler.InitLogger()

	// initate db
	db := inithandler.InitDB(logger)
	if db == nil {
		log.Fatal("[debug]", "failed to connect to db")
	}
	logger.Log("[debug]", "The database is connected")

	// initate repository
	var repository = inithandler.InitRepository(db)

	// initiate invoice bl
	var invoiceBL = invoice.NewBL(logger, repository.InvoiceRepo)
	// initiate user bl
	var userBL = user.NewBL(logger, repository.UserRepo)

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
	logger.Log("debug", "Starting the server on ", "port", config.WebPort)
	errs := make(chan error)
	go func() {
		logger.Log("transport", "http", "address", config.WebPort, "msg", "listening")
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", config.WebPort), router)
	}()

	logger.Log("terminated", <-errs)
}
