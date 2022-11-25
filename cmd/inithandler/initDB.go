package inithandler

import (
	"database/sql"
	"time"

	"github.com/go-kit/log"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
)

// This function will make a connection to the database only once.
func InitDB(logger log.Logger) *sql.DB {
	// connStr := "postgres://postgres:password@127.0.0.1/github.com/invoice-service?sslmode=disable"

	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	return nil, err
	// }

	// if err = db.Ping(); err != nil {
	// 	return nil, err
	// }
	return connectToDB(logger, 3, 5)
}

func openDB(logger log.Logger, dsn string) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		logger.Log("Failed to open Database", err.Error())
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB(logger log.Logger, retries, delay int) *sql.DB {
	for i := 0; i < retries; i++ {
		dsn := viper.GetString("DSN")
		connection, err := openDB(logger, dsn)
		if err != nil {
			logger.Log("[debug]", " Postgres not ready yet", "err", err.Error())
			time.Sleep(time.Duration(delay) * time.Second)
			continue
		}
		logger.Log("[debug]", "Successfully connected postgres DB")
		return connection
	}
	return nil
}