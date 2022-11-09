package inithandler

import (
	"database/sql"

	"github.com/go-kit/log"
)

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
