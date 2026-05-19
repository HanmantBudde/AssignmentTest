CREATE TABLE IF NOT EXISTS todos (
    id         TEXT PRIMARY KEY,
    text       TEXT NOT NULL,
    due_date   TIMESTAMPTZ NOT NULL,
    completed  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Backfill timestamp columns on pre-existing installations.
ALTER TABLE todos ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE todos ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS todos_due_date_idx ON todos (due_date);
