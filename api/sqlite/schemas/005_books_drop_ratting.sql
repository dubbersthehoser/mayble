
-- +goose Up
-- Remove rattings column from books table
CREATE TABLE new_books(
	id         INTEGER PRIMARY KEY,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	title      TEXT    NOT NULL,
	author     TEXT    NOT NULL,
	genre      TEXT    NOT NULL
);

INSERT INTO new_books(
	id,
	created_at,
	updated_at,
	title,
	author,
	genre
)
SELECT 
	id,
	created_at,
	updated_at,
	title,
	author,
	genre
FROM 
	books;

DROP TABLE books;

ALTER TABLE new_books RENAME TO books;


-- +goose Down
CREATE TABLE old_books(
	id         INTEGER PRIMARY KEY,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	title      TEXT    NOT NULL,
	author     TEXT    NOT NULL,
	genre      TEXT    NOT NULL,
	ratting    INTEGER NOT NULL CHECK(ratting >= 0 AND ratting <= 5)
);

INSERT INTO old_books(
	id,
	created_at,
	updated_at,
	title,
	author,
	genre,
	ratting
) 
SELECT
	id,
	created_at,
	updated_at,
	title,
	author,
	genre,
	0
FROM
	books;

-- update ratting values
UPDATE old_books
SET ratting = (
	SELECT read_books.rating 
	FROM read_books
	WHERE read_books.book_id = old_books.id
)
WHERE old_books.id IN (
	SELECT read_books.book_id
	FROM read_books
);

DROP TABLE books;

ALTER TABLE old_books RENAME TO books;
