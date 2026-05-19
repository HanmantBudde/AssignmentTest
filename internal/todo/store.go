package todo

import (
	"context"
	"crypto/rand"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"errors"
	"time"
)

//go:embed schema.sql
var schemaSQL string

var ErrNotFound = errors.New("todo not found")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ApplySchema(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schemaSQL)
	return err
}

func (s *Store) Create(ctx context.Context, text string, dueDate time.Time) (Todo, error) {
	t := Todo{
		ID:      newID(),
		Text:    text,
		DueDate: dueDate,
	}
	const insertQ = `INSERT INTO todos (id, text, due_date, completed) VALUES ($1, $2, $3, $4)`
	if _, err := s.db.ExecContext(ctx, insertQ, t.ID, t.Text, t.DueDate, t.Completed); err != nil {
		return Todo{}, err
	}
	const selectQ = `SELECT created_at, updated_at FROM todos WHERE id = $1`
	if err := s.db.QueryRowContext(ctx, selectQ, t.ID).Scan(&t.CreatedAt, &t.UpdatedAt); err != nil {
		return Todo{}, err
	}
	return t, nil
}

func (s *Store) Get(ctx context.Context, id string) (Todo, error) {
	const q = `
		SELECT id, text, due_date, completed, created_at, updated_at
		FROM todos
		WHERE id = $1`
	var t Todo
	err := s.db.QueryRowContext(ctx, q, id).
		Scan(&t.ID, &t.Text, &t.DueDate, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Todo{}, ErrNotFound
	}
	if err != nil {
		return Todo{}, err
	}
	return t, nil
}

// Update applies non-nil fields to the todo with the given id. NULL parameters
// are preserved via COALESCE so unspecified fields stay untouched. updated_at
// is always bumped to now().
func (s *Store) Update(ctx context.Context, id string, text *string, dueDate *time.Time, completed *bool) (Todo, error) {
	const updateQ = `
		UPDATE todos
		SET text       = COALESCE($2, text),
		    due_date   = COALESCE($3, due_date),
		    completed  = COALESCE($4, completed),
		    updated_at = now()
		WHERE id = $1`
	res, err := s.db.ExecContext(ctx, updateQ, id, text, dueDate, completed)
	if err != nil {
		return Todo{}, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return Todo{}, err
	}
	if n == 0 {
		return Todo{}, ErrNotFound
	}
	return s.Get(ctx, id)
}

func (s *Store) Delete(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
