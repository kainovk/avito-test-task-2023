package postgres

import (
	"avito-test-task-2023/internal/config"
	"avito-test-task-2023/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

	err = initSchema(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	return &Storage{db}, nil
}

func initSchema(db *sql.DB) error {
	op := "storage.postgres.initSchema"

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users
		(
			id   BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS segments
		(
			id   BIGSERIAL PRIMARY KEY,
			name VARCHAR(512) UNIQUE NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS user_segments
		(
			id         BIGSERIAL PRIMARY KEY,
			user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE,
			segment_id BIGINT REFERENCES segments (id) ON DELETE CASCADE,
			UNIQUE (user_id, segment_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveUser(name string) error {
	const op = "storage.postgres.SaveUser"

	_, err := s.db.Exec(`INSERT INTO users(name) VALUES ($1);`, name)
	if err != nil {
		// handle unique constraint error
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUser(id int64) (*sql.Row, error) {
	const op = "storage.postgres.GetUser"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return row, nil
}

func (s *Storage) DeleteUser(userId string) error {
	const op = "storage.postgres.DeleteUser"

	res, err := s.db.Exec(`DELETE FROM users WHERE id = $1;`, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrUserNotExists)
	}

	return nil
}

func (s *Storage) SaveSegment(name string) error {
	const op = "storage.postgres.SaveSegment"

	_, err := s.db.Exec(`INSERT INTO segments(name) VALUES ($1);`, name)
	if err != nil {
		// handle unique constraint error
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetSegment(id int64) (*sql.Row, error) {
	const op = "storage.postgres.GetSegment"

	stmt, err := s.db.Prepare("SELECT * FROM segments WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrSegmentNotFound
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return row, nil
}

func (s *Storage) DeleteSegment(segmentId string) error {
	const op = "storage.postgres.DeleteSegment"

	res, err := s.db.Exec(`DELETE FROM segments WHERE id = $1;`, segmentId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrSegmentNotExists)
	}

	return nil
}

func (s *Storage) AddUserSegments(userID int64, segmentIDs []int64) error {
	const op = "storage.postgres.AddUserSegments"

	for _, segmentID := range segmentIDs {
		_, err := s.db.Exec(`INSERT INTO user_segments(user_id, segment_id) VALUES ($1, $2);`, userID, segmentID)
		if err != nil {
			// handle unique constraint error
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrUserAlreadyHaveSegment)
			}

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) GetUserSegments(userID int64) ([]string, error) {
	const op = "storage.postgres.GetUserSegments"

	rows, err := s.db.Query(`
        SELECT s.name
        FROM user_segments AS usr
        JOIN segments AS s ON usr.segment_id = s.id
        WHERE usr.user_id = $1;
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var segments []string
	for rows.Next() {
		var segmentName string
		if err := rows.Scan(&segmentName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		segments = append(segments, segmentName)
	}

	return segments, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
