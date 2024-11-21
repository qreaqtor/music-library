package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5"
)

func GetConnect(connStr string, connAttempts int, attemptsDelay time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to PostgreSQL: %v", err)
	}

	tiker := time.NewTicker(attemptsDelay)
	defer tiker.Stop()

	for i := 0; i < connAttempts; i++ {
		_, ok := <-tiker.C
		if !ok {
			break
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}
	}

	return nil, fmt.Errorf("can't connect to PostgreSQL")
}
