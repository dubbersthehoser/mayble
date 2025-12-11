-- +goose Up
CREATE TABLE books (
	id         INTEGER PRIMARY KEY,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	title      TEXT NOT NULL,
	author     TEXT NOT NULL,
	genre      TEXT NOT NULL,
	ratting    INTEGER NOT NULL CHECK(ratting >= 0 AND ratting <= 5)
);

-- +goose Down
DROP TABLE books;
