
# Mayble Manual

Version: 1.0.0

## Overview

1. Database management menu. Creating, and opening database files, and importing, and exporting via CSV.

1. Undo and redo changes.

1. Create entry, edit, or delete selected book entry.

1. Next, and previous search matched item.

1. Text search entry.

1. Search field selection.

1. Table header, and table ordering.

![Overview](images/overview.png)


## Database Management Menu

![Dialog Window](images/database-management-menu.png)

### Creating Files

When creating a file, be it exporting `.csv` or creating database file, with in the file dialog, the file name extension is **not required**. The appropriate extension will be append to the name, if is not found in given file name.

### Database File

The database file is a sqlite3 database file.

At first start a database file will be created at `$HOME/Documents/mayble.db`. 
If that directory dose not exists then `$HOME/mayble.db` will be used.

### Exporting / Importing CSV

The structure of the CSV is:

```
TITLE,AUTHOR,GENRE,RATTING,BORROWER,DATE
TITLE,AUTHOR,GENRE,RATTING,,
```

**RULES**

- The should be no column header in the `.csv` file.

- `TITLE`, `AUTHOR`, `GENRE`, and `RATTING` fields must be filled in.

- If it's on loan then `BORROWER` and `DATE` must be filled in, otherwise keep blank.

- *RATTING:* can only be 0, 1, 2, 3, 4, and 5.

- *DATE* is in a `YYYY-MM-DD` format.


## Text Search

![Search](images/search.png)

Search only searches specific column from table, and they have to only contain text.
By default the search column is set to title, and can be with the selection box (right of the search box.)

When entering text in to search entry, an case insensitive substring match will be preformed.
Resulting in a selection ring, which can be navigated with the arrow buttons or inserting <ENTER> in the text box will
selecting the next selection in the ring.
