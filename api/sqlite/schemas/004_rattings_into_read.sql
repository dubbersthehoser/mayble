-- +goose Up
INSERT INTO read_books( -- prenounced 'red books'
	created_at,
	updated_at,
	rating,
	date_completed,
	book_id) 
SELECT 
	unixepoch(), 
	unixepoch(), 
	ratting, 
	date('now'), 
	id
FROM books WHERE ratting > 0;

-- +goose Down
DELETE FROM read_books;


