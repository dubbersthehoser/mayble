
-- name: CreateRead :one
INSERT INTO read_books(created_at, updated_at, rating, date_completed, book_id)
VALUES (
	unixepoch(),
	unixepoch(),
	?,
	?,
	?
)
RETURNING *;

-- name: DeleteRead :exec
DELETE FROM read_books WHERE book_id = ?;

-- name: UpdateRead :one
UPDATE read_books
SET
	updated_at = unixepoch(),
	rating = ?,
	date_completed = ?
WHERE book_id = ?
RETURNING *;

-- name: GetReadByBookID :one
SELECT rating, date_completed, book_id FROM read_books
WHERE book_id = ?;

