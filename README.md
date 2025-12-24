# Mayble

A simple desktop book management application for a family member.

My sister has a collection of physical books, and wants to keep track of what book has been read, and to whom she lent them out to.

She asked me to make this application to help her out.

![Screenshot](Screenshot.png)

## The Requirements

The [requirements](requirements.txt) my sister gave me. 

Technical requirements were to have builds for MacOS and ChromeOS. 

The MacOS build is created using [OSX-KVM](https://github.com/kholia/OSX-KVM) for a proper build environment (I hit road blocks using cgo for Linux to Mac cross-compilation).

For ChromeOS a `.deb` build is used with ChromeOS's Linux development environment.

Windows is not a priority for my sister.

## Current Features

- Sort by:
  + Title
  + Author
  + Genre
  + Ratting
  + Borrower
  + and Date

- Search by:
  + Title
  + Author
  + Genre
  + and Borrower

- The ability to Undo and Redo changes.

- To import and export by CSV.

# Build

