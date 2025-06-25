-- name: CreateBook :one
INSERT INTO books(created_at, updated_at, title, author, genre, ratting)
VALUES (
	?,
	?,
	?,
	?,
	?,
	?
)
RETURNING *;

-- name: DestroyBook :exec
DELETE FROM books WHERE id = ?;

-- name: GetAllBooks :many
SELECT id, title, author, genre, ratting FROM books;

-- name: UpdateBookRatting :exec
UPDATE books
SET updated_at = ?, ratting = ?
WHERE id = ?;

-- name: UpdateBookAuthor :exec 
UPDATE books 
SET updated_at = ?, author = ?
WHERE id = ?;
