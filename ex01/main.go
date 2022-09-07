/*
convert converts images from some format to another.
It converts all the images in the path including ones in the sub directories.

By default, it converts JPG images to PNG images.

Usage:

	./convert image_dir [options]

The flags are:

	-i string
		input file extension
		<png, jpeg, jpg> (default "jpg")
	-o string
		output file extension
		<png, jpeg, jpg> (default "png")
*/
package main

import (
	"os"

	"convert/imgconv"
)

// Convert images in a directory from some format to another.
func main() {
	app := &imgconv.App{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(app.Run(os.Args))
}
