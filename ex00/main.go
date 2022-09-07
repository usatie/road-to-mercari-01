package cat

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	app := &App{InStream: os.Stdin, OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(app.Run(os.Args))
}

// Exit Status
const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 1
	ExitCodeConvertError   = 1
)

// App consists of output/error streams.
type App struct {
	InStream             io.Reader
	OutStream, ErrStream io.Writer
	args                 Args
}

type Args struct {
	filenames []string
}

// Run is the entry point to the cli. Parses the arguments slice and routes to the proper flag/args combination
func (a *App) Run(args []string) int {
	if err := a.parseFlags(args); err != nil {
		return ExitCodeParseFlagError
	}

	if len(a.args.filenames) == 0 {
		return a.catInStream()
	} else {
		return a.catFiles()
	}
}

func (a *App) catInStream() int {
	status := ExitCodeOK
	if err := cat(a.InStream, a.OutStream); err != nil {
		fmt.Fprintln(a.ErrStream, err)
		status = ExitCodeConvertError
	}
	return status
}

func (a *App) catFiles() int {
	status := ExitCodeOK
	for _, file := range a.args.filenames {
		if f, err := os.Open(file); err != nil {
			fmt.Fprintln(a.ErrStream, err)
			status = ExitCodeConvertError
		} else if err := cat(f, a.OutStream); err != nil {
			fmt.Fprintln(a.ErrStream, err)
			status = ExitCodeConvertError
		}
	}
	return status
}

// Concat input strings
const bufSize = 1024

func cat(r io.Reader, w io.Writer) error {
	buf := make([]byte, bufSize)
	for {
		if n, err := r.Read(buf); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		} else if n > 0 {
			if n, err := w.Write(buf[:n]); err != nil {
				return err
			} else if n < bufSize {
				return nil
			}
		}
	}
	return nil
}

// Parse flags/args
func (a *App) parseFlags(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(a.ErrStream)
	fs.Usage = func() {
		fmt.Fprintf(a.ErrStream, "usage: %s [file ...]\n", args[0])
	}
	if err := fs.Parse(args[1:]); err != nil {
		return errors.New("Parse Error")
	}
	a.args.filenames = fs.Args()
	return nil
}
