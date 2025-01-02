# bldiff

A basic file diff tool written in Go!

## Usage

 - Build the binary with `go build .`
 - Pass two filenames to the binary to see similarities and differences.

`bldiff <file1> <file2>`

## About

This tool is much simpler than the built in `diff` command that comes preinstalled on Unix-like systems. The goal was to create something that was a bit less cryptic to read and that could be used more easily at a glace. It does not have any fancy features beyond basic line comparison and colored output.
