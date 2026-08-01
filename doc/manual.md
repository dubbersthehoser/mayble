# Manual

Mayble v2.0.2

- Github: [https://github.com/dubbersthehoser/mayble/](https://github.com/dubbersthehoser/mayble/)

## Overview

There are three major visual elements: Menu Bar, Status Line, and View.

- Menu Bar is at the top, and contains *File*, *Table*, and *Help*.
    + *File* create or open a database, or import / export with CSV from this drop down.
    + *Table* set settings for hiding a set of columns from the table view.
    + *Help* to open this manual.

- Status Line is on the top of every view, (and left of the table controls), and will give a status of an action taken with in the application.

- Body is the current open view for user interactions. Being: Table, Create, and Edit Entry views

*note: There are other views that pertain to errors, and status of the application.*

## Table View

The table view shows all the book entries in the database.

- *Column widths:* Adjust, by moving the mouse in between two column headers, hold left click and move left to right to resize. Sizes will be saved for next session.

- *Sorting* Press one of the column header buttons you desire to sort by. Starting from *descending* order, and then *ascending*, when pressed a second time. When sorted, none, or empty values will always be forced to the bottom of the table. Keeping all actual value to the top of the table.

### Controls

There are four controls and three of them are active when a entry is selected from the table. Being `Edit` ✏️, `Delete` 🗑️, and Cancel '❌'. 

- `Create` (➕) opens the create entry view.
- `Edit` (✏️) opens the edit view for editing the selected entry. 
- `Delete` (🗑️) removes the selected entry from the database, but it must be pressed twice.
- `Cancel` (❌) deselects the current selected entry.

### Search

Once text has been entered into the search box, will create a grouping of nearest matches and select highest scored cell with that search. You can cycle through them by:

- Press ENTER to move down to the next search result.

- Press CTRL+ENTER to move back up the search result.

These hotkeys only work when you focus is in the search entry.

The searches are **case insensitive**.

### Search By

You can change the style of search, and narrow down to a single column, or to every cell in the table.

## Create Entry View

With in this view you can create as many book entries as one needs. Press `Cancel` to go back to table view. Any errors, or successes will be shown on the Status Line. To Open this view check **Controls** section under **Table View**.

## Edit Entry View

This view is similar to the create view, but it will go back to the table view once the entry is updated.

## Import / Export

### CSV Format Rules

- The header row is for the CSV file is '`Title,Author,Genre,Completed,Rating,Loaned,Borrower`'.

- `Title`, `Author`, and `Genre` column values **MUST** be filled in, and not left blank.

- `Completed` and `Loaned` columns must be in a YYYY-MM-DD format. 
 
- `Rating` values can only be from 1 to 5, and the rating 0 and blank are treated as the same, otherwise it will be an error.  

- `Completed`, and `Rating`, or `Loaned`, and `Borrower`  will be considered not filled if one is blank out of the two. e.g. If `Completed` is empty and `Rating` is filled as 1, both will be consider empty, and `Rating` will be not be present, as well with `Loaned`, and `Borrower`.

