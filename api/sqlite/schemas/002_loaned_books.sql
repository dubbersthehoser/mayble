-- +goose Up
CREATE TABLE loaned_books (
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	date       TEXT NOT NULL,
	name       TEXT NOT NULL,
	book_id    INTEGER NOT NULL,
	FOREIGN KEY(book_id) REFERENCES books(id) 
		ON DELETE CASCADE
		ON UPDATE NO ACTION,
	UNIQUE(book_id)
);

-- +goose Down
DROP TABLE loaned_books;
