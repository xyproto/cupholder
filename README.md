# cupholder [![Build Status](https://travis-ci.com/xyproto/cupholder.svg?branch=master)](https://travis-ci.com/xyproto/cupholder) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/cupholder)](https://goreportcard.com/report/github.com/xyproto/cupholder) [![License](https://img.shields.io/badge/License-GPL2-brightgreen)](https://raw.githubusercontent.com/xyproto/cupholder/master/LICENSE)

Eject the CD tray, on Linux, using only Go.

If you have a server room with many servers with CD ROMs, this can be used for finding a specific server.

It can also be used for providing an additional tea or coffee cup holder to your desktop setup.

## Usage

    ./cupholder

Or

    ./cupholder /dev/cdrom

All arguments are treated as device filenames.

Use `-s` for silent operation and set the `NO_COLOR` environment variable to disable colors.

## Screenshot

![cupholder in action](img/screenshot.png)

## Web server

`cupholder` includes a HTTP server that can be launched with the `-l` or `--listen` flag, that lets you eject the CD remotely, just by visiting the HTTP page at port 3280 (`3280 = 0x0CD0`).

![webpage](img/httpserver.png)

## Requirements

* Go >= 1.10

## General info

* Version: 1.0.0
* License: GPL2
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;


