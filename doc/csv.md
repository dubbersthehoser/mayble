# CSV Format

Mayble Version 2.0.2

## Schema

The default schema header is `Title,Author,Genre,Completed,Rating,Loaned,Borrower`, but the order dose not matter when the header schema is present with correct labels. When the imported file has no schema header the default schema will be used. 

It is **recommend** to use default schema header, and to use the default header order, when modifying, or creating a CSV file.

**Values**

For `Completed` and `Loaned` value must be in a YYYY-MM-DD format, and for `Rating` values can only be from 0-5. 

With `Title`, `Author`, `Genre` values **MUST** be filled in, and not left blank. When one of the two, either `Completed`, and `Rating`, or `Loaned`, and `Borrower`, are both will be considered empty. e.g. If `Completed` is empty and `Rating` is filled as 1, both will be consider empty, and `Rating` will be imported as empty / zero.

