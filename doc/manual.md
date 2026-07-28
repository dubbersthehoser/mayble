# Manual

Mayble Version 2.0.2

## CSV Format

- The schema header is `Title,Author,Genre,Completed,Rating,Loaned,Borrower`.

- With `Completed` and `Loaned` columns, value must be in a YYYY-MM-DD format, and for `Rating` values can only be from 1-5. The rating 0 and blank are the same.  

- With `Title`, `Author`, and `Genre` column values **MUST** be filled in, and not left blank.

- With `Completed`, and `Rating`, or `Loaned`, and `Borrower`, both will be considered not filled if one is blank. e.g. If `Completed` is empty and `Rating` is filled as 1, both will be consider empty, and `Rating` will be not be present.

