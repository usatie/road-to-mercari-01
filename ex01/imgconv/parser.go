package imgconv

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// cliArgument is a collection of available options
type cliArgument struct {
	inExt    string
	outExt   string
	rootPath string
	quality  int
	dir      string
	verbose  bool
	encoder  Encoder
}

func (c cliArgument) String() string {
	return fmt.Sprintf("Options:\n  inExt: %s\n  outExt: %s\n  rootPath: %s\n  quality: %d\n  dir: %s\n  verbose: %v", c.inExt, c.outExt, c.rootPath, c.quality, c.dir, c.verbose)
}

const usageOptions = `options:
  -i, --input <extension>    input file format <png, jpeg, gif>
  -o, --output <extension>   output file format <png, jpeg, gif>
  -d, --dir <directory>      destination directory for output files
  -q, --quality <quality>    output encoding quality. higher is better. <1-100>
  -v, --verbose              verbose output`

// parseFlags parse the arguments and writes errors to errStream
func (a *App) parseFlags(args []string) (cliArgument, error) {
	var arg cliArgument
	// Custom flag set to set custom output and usage
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(a.ErrOutput)
	fs.Usage = func() {
		fmt.Fprintf(a.ErrOutput, "Usage: %s [options] image_dir\n", args[0])
		fmt.Fprintln(a.ErrOutput, usageOptions)
	}
	// -i, --input
	fs.StringVar(&arg.inExt, "i", "jpeg", "input file format <png, jpeg, gif>")
	fs.StringVar(&arg.inExt, "input", "jpeg", "input file format <png, jpeg, gif>")
	// -o, --output
	fs.StringVar(&arg.outExt, "o", "png", "output file format <png, jpeg, gif>")
	fs.StringVar(&arg.outExt, "output", "png", "output file format <png, jpeg, gif>")
	// -d, --dir
	fs.StringVar(&arg.dir, "d", "", "destination directory for output files")
	fs.StringVar(&arg.dir, "dir", "", "destination directory for output files")
	// -q, --quality
	fs.IntVar(&arg.quality, "q", 75, "output encoding quality. higher is better. <1-100>")
	fs.IntVar(&arg.quality, "quality", 75, "output encoding quality. higher is better. <1-100>")
	// -v, --verbose
	fs.BoolVar(&arg.verbose, "v", false, "verbose output")
	fs.BoolVar(&arg.verbose, "verbose", false, "verbose output")
	// Parse
	if err := fs.Parse(args[1:]); err != nil {
		return arg, errors.New("Parse Error")
	}
	// rootPath
	arg.rootPath = fs.Arg(0)
	// rootPath must be passed
	if arg.rootPath == "" {
		fmt.Fprintln(a.ErrOutput, "error: invalid argument")
		return arg, errors.New("rootPath is empty")
	}
	// rootPath must exist
	if _, err := os.Stat(arg.rootPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(a.ErrOutput, "error: %s: no such file or directory\n", arg.rootPath)
		} else {
			fmt.Fprintf(a.ErrOutput, "error: %v\n", err)
		}
		return arg, err
	}
	// inExt
	switch arg.inExt {
	case "png", "jpeg", "gif":
		// ok
	default:
		fmt.Fprintln(a.ErrOutput, "Invalid input format.")
		return arg, errors.New("Invalid input format.")
	}
	// outExt
	switch arg.outExt {
	case "png", "jpeg", "gif":
		// ok
	default:
		fmt.Fprintln(a.ErrOutput, "Invalid input format.")
		return arg, errors.New("Invalid input format.")
	}
	// dir must be directory
	if arg.dir != "" {
		info, err := os.Stat(arg.dir)
		if err != nil {
			fmt.Fprintf(a.ErrOutput, "error: %v\n", err)
			return arg, err
		}
		if !info.IsDir() {
			err = errors.New(fmt.Sprintf("%s is not a directory.", arg.dir))
			fmt.Fprintf(a.ErrOutput, "error: %v\n", err)
			return arg, err
		}
	}
	// Encoder
	encoder, err := newEncoder(arg.outExt, arg.quality)
	if err != nil {
		fmt.Fprintf(a.ErrOutput, "error: %v\n", err)
		return arg, err
	}
	arg.encoder = encoder
	return arg, nil
}
