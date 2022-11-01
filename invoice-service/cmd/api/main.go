package main

import (
	"invoicing/invoice-service"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

const webPort = ":8080"

func main() {
	// initiate logger
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "method=", log.DefaultCaller)

	// initiate bl
	bl := invoice.NewBL(logger)

	router := invoice.NewHTTPServer(logger, bl)
	// start the server

	logger.Log("Starting the server on port", webPort)
	err := http.ListenAndServe(webPort, router)
	if err != nil {
		logger.Log(err)
	}

}
