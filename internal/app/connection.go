package app

import (
	"database/sql"
	"fmt"

	"github.com/qreaqtor/music-library/internal/config"

	_ "github.com/lib/pq"
)

func getPostgresConn(cfg config.PostgresConfig) (*sql.DB, error) {
	sslMode := "disable"
	if cfg.SSL {
		sslMode = "enable"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DB,
		sslMode,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to PostgreSQL: %v", err)
	}

	return conn, nil
}
