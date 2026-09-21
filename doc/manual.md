# Manual

Mayble v2.0.2

- Github: [https://github.com/dubbersthehoser/mayble/](https://github.com/dubbersthehoser/mayble/)

## Overview

This manual contains information on hidden features and non-obvious details about this program.

## Search

Once text has been entered into the search box, will create a grouping of nearest matches and select highest scored cell with that search. You can cycle through them by:

- Pressing ENTER to move down to the next search results, or

- Pressing CTRL+ENTER to move back up the search results.

These hotkeys only work when you focus is in the search entry. 
And searches are **case insensitive**.

## CSV Importing Format Rules

- The header row is for the CSV file is '`Title,Author,Genre,Completed,Rating,Loaned,Borrower`'.

- `Title`, `Author`, and `Genre` column values **MUST** be filled in, and not left blank.

- `Completed` and `Loaned` columns must be in a YYYY-MM-DD format. 
 
- `Rating` values can only be from 1 to 5, and the rating 0 and blank are treated as the same, otherwise it will be an error.  

- `Completed`, and `Rating`, or `Loaned`, and `Borrower`  will be considered not filled if one is blank out of the two. e.g. If `Completed` is empty and `Rating` is filled as 1, both will be consider empty, and `Rating` will be not be present, as well with `Loaned`, and `Borrower`.

