
-- name: CreateLoan :one
INSERT INTO loaned_books( created_at, updated_at, name, date, book_id)
VALUES (
	unixepoch(),
	unixepoch(),
	?,
	?,
	?
)
RETURNING *;

-- name: DeleteLoan :exec
DELETE FROM loaned_books WHERE id = ?;

-- name: UpdateLoan :one
UPDATE loaned_books
SET
	updated_at = unixepoch(),
	name = ?,
	date = ?
WHERE id = ?
RETURNING *;

-- name: GetAllLoans :many
SELECT id, name, date, book_id FROM loaned_books;

