-- +goose Up
CREATE TABLE books (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	title      TEXT NOT NULL,
	author     TEXT NOT NULL,
	genre      TEXT NOT NULL,
	ratting    INTEGER NOT NULL
);

-- +goose Down
DROP TABLE books;
