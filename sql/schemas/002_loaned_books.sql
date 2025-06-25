-- +goose Up
CREATE TABLE loaned_books (
	id         INTEGER PRIMARY KEY,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	date       INTEGER NOT NULL,
	name       TEXT NOT NULL,
	book_id    INTEGER NOT NULL,
	FOREIGN KEY(book_id) REFERENCES books(id) ON DELETE CASCADE,
	UNIQUE(book_id)
);

-- +goose Down
DROP TABLE loaned_books;
