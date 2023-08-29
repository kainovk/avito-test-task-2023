package postgres

import (
	"avito-test-task-2023/internal/config"
	"avito-test-task-2023/internal/http-server/handlers/users"
	"avito-test-task-2023/internal/models/segment"
	"avito-test-task-2023/internal/models/user"
	"avito-test-task-2023/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
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
			slug VARCHAR(512) UNIQUE NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS user_segments
		(
			id         BIGSERIAL PRIMARY KEY,
			user_id    BIGINT REFERENCES users (id) ON DELETE CASCADE,
			segment_id BIGINT REFERENCES segments (id) ON DELETE CASCADE,
			delete_at  TIMESTAMP DEFAULT NULL,
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

func (s *Storage) GetUser(id int64) (*user.User, error) {
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

	var usr *user.User
	err = row.Scan(&usr.ID, &usr.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to scan user name: %w", op, err)
	}

	return usr, nil
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

func (s *Storage) SaveSegment(slug string) error {
	const op = "storage.postgres.SaveSegment"

	_, err := s.db.Exec(`INSERT INTO segments(slug) VALUES ($1);`, slug)
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

func (s *Storage) GetSegment(id int64) (*segment.Segment, error) {
	const op = "storage.postgres.GetSegment"

	stmt, err := s.db.Prepare("SELECT id, slug FROM segments WHERE id = $1")
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

	var seg *segment.Segment
	err = row.Scan(&seg.ID, &seg.Slug)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to scan segment name: %w", op, err)
	}

	return seg, nil
}

func (s *Storage) GetSegmentBySlug(slug string) (*segment.Segment, error) {
	const op = "storage.postgres.GetSegmentBySlug"

	stmt, err := s.db.Prepare("SELECT id, slug FROM segments WHERE slug = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRow(slug)
	if row.Err() != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrSegmentNotFound
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	seg := &segment.Segment{}
	err = row.Scan(&seg.ID, &seg.Slug)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to scan segment name: %w", op, err)
	}

	return seg, nil
}

func (s *Storage) GetSegments() ([]*segment.Segment, error) {
	const op = "storage.postgres.GetSegments"

	rows, err := s.db.Query(`SELECT id, slug FROM segments`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var segments []*segment.Segment
	for rows.Next() {
		seg := &segment.Segment{}
		if err := rows.Scan(&seg.ID, &seg.Slug); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		segments = append(segments, seg)
	}

	return segments, nil
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

func (s *Storage) DeleteSegmentBySlug(slug string) error {
	const op = "storage.postgres.DeleteSegmentBySlug"

	res, err := s.db.Exec(`DELETE FROM segments WHERE slug = $1;`, slug)
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

func (s *Storage) AddUserSegmentsBySlugs(userID int64, segmentsToAdd []users.SegmentRequest) error {
	const op = "storage.postgres.AddUserSegmentsBySlugs"

	for _, segmentToAdd := range segmentsToAdd {
		seg, err := s.GetSegmentBySlug(segmentToAdd.Slug)
		if err != nil {
			continue
		}

		_, err = s.db.Exec(`
			INSERT INTO user_segments(user_id, segment_id, delete_at) VALUES ($1, $2, $3);
		`, userID, seg.ID, segmentToAdd.DeleteAt)
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

func (s *Storage) DeleteUserSegmentsBySlugs(userID int64, slugs []string) error {
	const op = "storage.postgres.DeleteUserSegmentsBySlugs"

	for _, slug := range slugs {
		seg, err := s.GetSegmentBySlug(slug)
		if err != nil {
			continue
		}

		_, err = s.db.Exec(`DELETE FROM user_segments WHERE user_id = $1 AND segment_id = $2;`, userID, seg.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) GetUserSegments(userID int64) ([]*segment.Segment, error) {
	const op = "storage.postgres.GetUserSegments"

	rows, err := s.db.Query(`
        SELECT s.id, s.slug
        FROM user_segments AS usr
        JOIN segments AS s ON usr.segment_id = s.id
        WHERE usr.user_id = $1;
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segments []*segment.Segment
	for rows.Next() {
		seg := &segment.Segment{}
		if err := rows.Scan(&seg.ID, &seg.Slug); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		segments = append(segments, seg)
	}

	return segments, nil
}

func (s *Storage) ConfigureUserSegments(userID int64, segAdd []users.SegmentRequest, segDel []string) error {
	// TODO: tx
	const op = "storage.postgres.ConfigureUserSegments"

	err := s.AddUserSegmentsBySlugs(userID, segAdd)
	if err != nil {
		return fmt.Errorf("%s: failed to add segments to user: %w", op, err)
	}

	err = s.DeleteUserSegmentsBySlugs(userID, segDel)
	if err != nil {
		return fmt.Errorf("%s: failed to delete segments from user: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSegmentsTTL() (int64, error) {
	const op = "storage.postgres.DeleteSegmentsTTL"

	currentTime := time.Now()

	res, err := s.db.Exec(`
		DELETE FROM user_segments 
		WHERE delete_at IS NOT NULL AND delete_at < $1
	`, currentTime)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return deleted, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
