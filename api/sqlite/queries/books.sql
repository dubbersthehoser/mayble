-- name: CreateBook :one
INSERT INTO books(created_at, updated_at, id, title, author, genre)
VALUES (
	unixepoch(),
	unixepoch(),
	NULL,
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
	genre  = ?

WHERE id = ?
RETURNING *;

-- name: GetAllBooks :many
SELECT id, title, author, genre FROM books;

-- name: GetBookByID :one
SELECT id, title, author, genre FROM books
WHERE id = ?;
