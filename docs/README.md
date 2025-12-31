
# Mayble Manual

Version: 1.0.0

## Creating Files

When creating a file, be it a `.csv` or `.db` file, the file name extension is **not required**. The appropriate extension will be append to the file name when not found during creation.


## Database File

At first start a database file will be created at `$HOME/Documents/mayble.db`. 
If that directory dose not exists then `$HOME/mayble.db` will be used.

The database file is a sqlite3 database file, and by default uses the `.db` extension.
`.sqlite` and `.sqlite3` are also valid extension.


## Exporting and Importing CSV

The structure of the CSV is:

```
TITLE,AUTHOR,GENRE,RATTING,BORROWER,DATE
```

**RULES**

- The should be no column header in the `.csv` file.

- `TITLE`, `AUTHOR`, `GENRE`, and `RATTING` fields must be filled in.

- If it's on loan then `BORROWER` and `DATE` must be filled in, otherwise keep blank.

- *RATTING:* can only be 0, 1, 2, 3, 4, and 5.

- *DATE* is in a `YYYY-MM-DD` format.


## Text Search

![Search](images/search.png)

Search can only search a specific column from the table, and they must only contain text.
By default the search column is set to Title, and can be changed with the selection box (right of the search box.)

When entering text in the search box, a case insensitive substring match will be preformed.
Resulting in a selection ring, which can be navigated with the arrow buttons, or pressing <ENTER> selecting the next item in the ring.

