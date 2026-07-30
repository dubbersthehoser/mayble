# Mayble

A simple desktop book management application.

![Mable Screen Shot](./images/Screenshot.png)

## Motivation

I was asked by a family member of mine, whom has a collection of physical books, and needed a way to keep track of them without needing pen and paper, and so I create this CRUD application for them. The application's aim is to keep track of what book has been read, and what has been loaned out, and to whom.

## Quick Start

There are pre-build packages in the [release page](https://github.com/dubbersthehoser/mayble/releases/latest).

*note:* For the **Linux** package, once extracted you can run the binary directly from the `/usr/share/bin/` directory.

### But my platform is not in the releases?

Then build from source.

1. Clone down repo: `git clone github.com/dubbersthehoser/mayble`

1. Install dependencies.
   - **Linux (Debian):** 
     ```
     sudo apt-get install golang gcc libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev
     ```

   - **MacOS:** As long you have `xcode` installed. 
     ```
     xcode-select --install
     ```
   - **Other:** Check [fyne quick start](https://docs.fyne.io/started/quick/).

1. Install the Go language (v1.24.0 or higher), and a C Compiler for graphics libraries.

1. Install Fyne command line tool `go install fyne.io/tools/cmd/fyne@latest`.

1. Build `fyne build .` then run `./mayble` and try it out!

## Usage

The application starts without a database. You'll need to create a database which is done in the menu *File* drop down.

**Create** entries press the ➕ (plus) button to open the create book entry view. Once there, You can continually add entires to the database, and will be notify of actions that succeed, or failed in the top blank spot of the view. Press `Cancel` to go back to the table view, to see your additions.

**Hide, and show columns** by opening the *Table* menu bar drop down, and pick ID, loaned, and read to hide, or show a set of columns. These do not filter out entires, and just hides those columns from the table.

**Column widths** are adjustable, by mousing in between two column headers, hold left click an move to resize.

**Sort table** by pressing one of the column header labels you desire to sort by. Starting as *descending* order, and then *ascending* when pressed again, and toggled back and forth. When sorted, any none / empty values will always be forced to the bottom of the table when sorted. Keeping all actual value to the top of the sorted table.

**Edit** an entry by selecting any cell in that row entry. Then press the ✏️ (pencil) button to open the edit view. This view behaves similar to the create entries view, but closes the view once submitted.

**Delete** an entry by selecting an entry, then pressing the 🗑️ (garbage can) twice, do to the first click being a warning with a cool down.

**Search** box will select a row cell by its highest scored match. You can change the `search by` setting using the select widget right of the search box. Searching with `All` will search cells from left to right and down, and down a column of cells for the other options. Searches are **case insensitive*.

## Contributing

I'm not looking for contributors for this project. My aim was to help one person, but I'll consider any issue reports or request if they pop up.

