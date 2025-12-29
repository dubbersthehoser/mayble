
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

### Opening Database File

At first start a database file will be created at `$HOME/Documents/mayble.db`. 
If that directory dose not exists then `$HOME/mayble.db` will be used.

### Creating Files

When creating a file, be it exporting `.csv` or creating database file, with in the file dialog, the file name extension is **not required**. The appropriate extension will be append to the name, if is not found.


## Text Search

Search searches specific fields from entries in the table.
By default the search field is set to title.
You can change search field with selection box (6) right of the search box.

When entering text in to search entry, an case insensitive substring match will be preformed.
Resulting in a item ring, which can be navigated with the arrow buttons. (4)

Pressing ENTER in the text box will go to the next item in the search selected item in ring.
