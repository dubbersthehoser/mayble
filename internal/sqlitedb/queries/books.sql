-- name: CreateBook :one
INSERT INTO books(created_at, updated_at, title, author, genre, ratting)
VALUES (
	unixepoch(),
	unixepoch(),
	?,
	?,
	?,
	?
)
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = ?;

-- name: UpdateBook :one
UPDATE books
SET 
	updated_at = unixepoch(),
	title  = ?,
	author = ?,
	genre  = ?,
	ratting = ?

WHERE id = ?
RETURNING *;

-- name: GetAllBooks :many
SELECT id, title, author, genre, ratting FROM books;

-- name: GetBookByID :one
SELECT id, title, author, genre, ratting FROM books
WHERE id = ?;
