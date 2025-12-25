# Mayble

A desktop application book management for a family member.

My sister has a collection of physical books, and wants to keep track of what book has been read, and to whom she lent them out to.

She asked me to make this application to help her out.

## The Requirements

The [requirements](requirements.txt) my sister gave me. 

Build targets are MacOS and ChromeOS. 

The MacOS build is created using [OSX-KVM](https://github.com/kholia/OSX-KVM) for a proper build environment (I hit road blocks using cgo for Linux to Mac cross-compilation).

For ChromeOS a `.deb` build is used with ChromeOS's Linux development environment.

Windows is not a priority for my sister.

## Current Features

![Screenshot](Screenshot.png)

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


# Building

## MacOS Setup

1. Install Xcode Command Line Tools.

``` sh
xcode-select --install
```
1. Install [Go](https://go.dev/dl)

NOTE: Can't run build under [OSX-VM](https://github.com/kholia/OSX-KVM) without GPU pass-through. OpenGL will crash  the app under a virtual Graphics.


## Quick Run

``` sh
go run .
```

## Packaging

Dependencies

``` sh
go install github.com/fyne-io/fyne-cross@latest
```

**MacOS**

With in a MacOS environment.

```
fyne-cross darwin -arch=ARCH  # ARCH = amd64, or arm64
```

Once finished the `.app` file will be shown in output of `fyne-cross`. Example:

``` sh
[✓] Package: "./fyne-cross/dist/darwin-*/Mayble.app"
```

**Debian**

NOTE: This package is not intended for general Debian deployment and only for [personal use](https://wiki.debian.org/MakeAPrivatePackage).

With in a Debian environment.

Dependencies

``` sh
sudo apt update
sudo apt install build-essential devscripts debhelper dh-make fakeroot
```

Package

``` sh
fyne-cross linux -arch=ARCH      # ARCH = amd64, or arm64
./package-deb.sh ARCH            # Create Debian packages of selected ARCH
```









