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

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
