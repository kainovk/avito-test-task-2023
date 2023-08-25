package postgres

import (
	"avito-test-task-2023/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(creds config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		creds.Host,
		creds.Port,
		creds.Username,
		creds.Password,
		creds.Database,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}
