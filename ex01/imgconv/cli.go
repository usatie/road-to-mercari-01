// Package cli provides functionalities to convert images from some format to another.
package imgconv

import (
	"bufio"
	"errors"
	"fmt"
	"image"

	//_ "image/gif"
	_ "image/jpeg"
	//_ "image/png"
	"io"
	"os"
	"path/filepath"
)

const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 1
	ExitCodeConvertError   = 1
)

// App consists of output/error streams.
type App struct {
	Input             io.Reader
	Output, ErrOutput io.Writer
	verbose           bool
}

func (a *App) verbosePrintf(format string, args ...any) {
	if a.verbose {
		fmt.Fprintf(a.Output, format, args...)
	}
}

func (a *App) convertFile(dstPath, srcPath, inputFormat string, encoder Encoder) error {
	a.verbosePrintf("Convert a file: %s ------> %s\n", srcPath, dstPath)
	// Open file
	fin, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer fin.Close()

	// Decode file as image
	img, err := decodeImage(fin, inputFormat)
	if err != nil {
		return fmt.Errorf("%s is not a valid file", srcPath)
	}

	/*
		Do not create file before decodeImage.
		Otherwise, empty file is created even when decodeImage fails.
	*/
	// Create output file
	fout, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer fout.Close()

	// Encode image to file
	err = encoder.Encode(bufio.NewWriter(fout), img)
	return err
}

func (a *App) convert(inputFormat string, encoder Encoder) error {
	// Decode stdin as image
	img, err := decodeImage(a.Input, inputFormat)
	if err != nil {
		return err
	}

	// Encode image to stdout
	err = encoder.Encode(bufio.NewWriter(a.Output), img)
	if err != nil {
		return err
	}
	return nil
}

func decodeImage(r io.Reader, expectedFormat string) (image.Image, error) {
	img, format, err := image.Decode(bufio.NewReader(r))
	if err == image.ErrFormat || "."+format != expectedFormat {
		return img, fmt.Errorf("Invalid Format")
	}
	return img, err
}

// Run is the entry point to the cli. Parses the arguments slice and routes to the proper flag/args combination
func (a *App) Run(args []string) int {
	arg, err := a.parseFlags(args)
	if err != nil {
		return ExitCodeParseFlagError
	}
	inExt := "." + arg.inExt
	outExt := "." + arg.outExt
	// verbose output (print args)
	a.verbosePrintf("%v\n", arg)

	paths := []string{}
	err = filepath.WalkDir(arg.rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("WalkDir: %w", err)
		}
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	var cnt int
	for _, path := range paths {
		dstPath := getOutputPath(path, arg.dir, outExt)
		err := a.convertFile(dstPath, path, inExt, arg.encoder)
		if err != nil {
			fmt.Fprintf(a.ErrOutput, "error: %v\n", err)
			return ExitCodeConvertError
		}
		cnt++
	}

	// verbose output (completion)
	a.verbosePrintf("\n\nconverted %d files\n", cnt)
	return ExitCodeOK
}

// removeExt removes the extension from path.
// "image/hoge.jpg" -> "image/hoge"
func removeExt(path string) string {
	ext := filepath.Ext(path)
	baseLen := len(path) - len(ext)
	baseName := path[:baseLen]
	return baseName
}

/*
getOutputPath returns the output file name which does not overwrite existing files..

getOutputPath("image/bear.jpg", "output", ".png")
// "output/bear.png"
getOutputPath("image/bear.jpg", "output", ".png")
// "output/bear (2).png"
getOutputPath("image/bear.jpg", "output", ".png")
// "output/bear (3).png"
*/
func getOutputPath(path, outputDir, outExt string) string {
	// "image/bear.jpg"
	outputPath := path
	if outputDir != "" {
		// "output/bear.jpg"
		outputPath = filepath.Join(outputDir, filepath.Base(path))
	}
	baseName := removeExt(outputPath)
	// "output/bear.png"
	outputPath = baseName + outExt
	for n := 2; ; n++ {
		if _, err := os.Stat(outputPath); err == nil {
			// Already exists
			// "output/bear (2).png"
			outputPath = fmt.Sprintf("%s (%d)%s", baseName, n, outExt)
			continue
		} else if errors.Is(err, os.ErrNotExist) {
			// Does not exists (new file name)
			break
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			panic("File may or may not exist. os.Stat error")
		}
	}
	return outputPath
}
