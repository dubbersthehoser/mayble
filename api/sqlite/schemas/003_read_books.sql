
-- +goose UP
CREATE TABLE read_books (
	-- Books that have been fulling read.
	created_at      INTEGER NOT NULL,
	updated_at      INTEGER NOT NULL,
	rating          INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
	date_completed  TEXT    NOT NULL,
	book_id         INTEGER NOT NULL,
	FOREIGN KEY(book_id) REFERENCES books(id)
		ON DELETE CASCADE
		ON UPDATE NO ACTION,
	UNIQUE(book_id)
);

-- +goose Down
DROP TABLE read_books;
