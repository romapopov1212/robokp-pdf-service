package db

import (
	"database/sql"
	"fmt"
	"github.com/romapopov1212/robokp-pdf-service/internal/config"
)

func NewDatabaseConnection(cfg config.Database) (*sql.DB, string, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open database:%w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, "", fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, dsn, nil
}
