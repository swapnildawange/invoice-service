package inithandler

import (
	"os"

	"github.com/go-kit/log"
)

func InitLogger() log.Logger {
	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(w)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "method=", log.DefaultCaller)
	return logger
}
