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

2. Install the Go language (v1.24.0 or higher)

2. Install Fyne command line tool `go install fyne.io/tools/cmd/fyne@latest`.

3. Build `fyne build .` then run `./mayble` and try it out!

## Usage

When the application fist starts, you need to create a database which can be done in the *File* menu bar drop down.

Once open press the **+** button to open the create book entry view. Once there, You can continually add entires to the database, and will be notify of actions that succeed, or failed in the top blank spot of the view. Press `Cancel` to go back to the table view, to see your additions.

The table view dose not show all the columns in the table and can be changed in the *Table* menu bar drop down. You can show, and hide the ID, loaned, and read columns from the entires. These do not filter out entires, and just hides those columns from the table.

## Contributing

I'm not looking for contributors for this project. My aim was to help one person, but I'll consider any issue reports or request if they pop up.


