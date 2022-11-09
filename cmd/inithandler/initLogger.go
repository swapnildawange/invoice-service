package inithandler

import (
	"os"

	"github.com/go-kit/log"
)

func InitLogger() log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "method=", log.DefaultCaller)
	return logger
}
