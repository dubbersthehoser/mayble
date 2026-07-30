# Manual

Mayble v2.0.2

- [Github](https://github.com/dubbersthehoser/mayble/)

## Table

**Hide, and show columns** by opening the *Table* menu bar drop down, and pick ID, loaned, and read to hide, or show a set of columns. These do not filter out entires, and just hides those columns from the table.

**Column widths** are adjustable, by mousing in between two column headers, hold left click an move to resize.

**Sort table** by pressing one of the column header labels you desire to sort by. Starting as *descending* order, and then *ascending* when pressed again, and toggled back and forth. When sorted, any none / empty values will always be forced to the bottom of the table when sorted. Keeping all actual value to the top of the sorted table.

**Search** box will select a row cell by its highest scored match. You can change the `search by` setting using the select widget right of the search box. Searching with `All` will search cells from left to right and down, and down a column of cells for the other options. Searches are **case insensitive*.

## Controls

**Create** entries press the ➕ (plus) button to open the create book entry view. Once there, You can continually add entires to the database, and will be notify of actions that succeed, or failed in the top blank spot of the view. Press `Cancel` to go back to the table view, to see your additions.

**Edit** an entry by selecting any cell in that row entry. Then press the ✏️ (pencil) button to open the edit view. This view behaves similar to the create entries view, but closes the view once submitted.

**Delete** an entry by selecting an entry, then pressing the 🗑️ (garbage can) twice, do to the first click being a warning with a cool down.

## Search

- Press ENTER to move down to the next nearest result.
- Press CTRL+ENTER to move back up the search result.

## CSV Format

- The schema header is `Title,Author,Genre,Completed,Rating,Loaned,Borrower`.

- With `Completed` and `Loaned` columns, value must be in a YYYY-MM-DD format, and for `Rating` values can only be from 1-5. The rating 0 and blank are the same.  

- With `Title`, `Author`, and `Genre` column values **MUST** be filled in, and not left blank.

- With `Completed`, and `Rating`, or `Loaned`, and `Borrower`, both will be considered not filled if one is blank. e.g. If `Completed` is empty and `Rating` is filled as 1, both will be consider empty, and `Rating` will be not be present.

