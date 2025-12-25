# Mayble

A desktop application book management for a family member.

My sister has a collection of physical books, and wants to keep track of what book has been read, and to whom she lent them out to. And so she asked me to make this application to help her out.

The [requirements](requirements.txt) my sister gave me, and the build targets are for MacOS and ChromeOS. Windows is not a priority for her. I'm also going with to use a Debian package for ChromeOS, since she has install Linux apps on her Chrome Book before.


## Features.

![Screenshot](Screenshot.png)

- Sorting by field.

- Field based Text Searching.

- Undo and Redo changes.

- Import and export with CSV.


# Building

## Requirements.

fyne-cross and [Go](https://go.dev/dl).

``` sh
go install github.com/fyne-io/fyne-cross@latest
```


## MacOS

The MacOS build is created using [OSX-KVM](https://github.com/kholia/OSX-KVM) for a proper build environment. (I hit road blocks using cgo and fyne-cross for Linux to Mac cross-compilation).

NOTE: Can't run app under [OSX-VM](https://github.com/kholia/OSX-KVM) without GPU pass-through. OpenGL will crash  the app under a virtual Graphics.

Install Xcode Command Line Tools.

``` sh
xcode-select --install
```

Build

``` sh
fyne-cross darwin -arch=ARCH  # ARCH = amd64, or arm64
```

Once finished the `.app` file will be shown in output of `fyne-cross`. Example:

``` sh
[✓] Package: "./fyne-cross/dist/darwin-*/Mayble.app"
```

## Debian / ChromeOS

The ChromeOS build is a `.deb` package for the Linux Development Environment. The code writing assumes a desktop environment, and will crash if the `fyne/driver` is a phone environment.

NOTE: This package is not intended for general Debian deployment and only for [personal use](https://wiki.debian.org/MakeAPrivatePackage).


Requirements.

``` sh
sudo apt update
sudo apt install build-essential devscripts debhelper dh-make fakeroot
```

Packaging.

``` sh
fyne-cross linux -arch=ARCH      # ARCH = amd64, or arm64
./package-deb.sh ARCH            # Create Debian packages of selected ARCH
```

The `package-deb.sh` will take the Linux tarball created by `fyne-cross` and creates a Debian package with it.

The resulting Debian package will be in `./build/deb/mayble-X.X.X.deb`


