package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

/*
Tests can be done with the following script

```
$ vim test.sh
$ chmod +x test.sh
$ ./test.sh
```

`test.sh`
```
#!/bin/bash
echo "hello" > hello
echo "world" > world
echo "secret" > secret
chmod 200 secret
rm -f nosuchfile
diff -U 3 <(cat hello) <(./ft_cat hello)
diff -U 3 <(cat hello world) <(./ft_cat hello world)
diff -U 3 <(cat <hello) <(./ft_cat <hello)
diff -U 3 <(echo "ola" | cat) <(echo "ola" | ./ft_cat)
diff -U 3 <(echo "ola" | cat hello) <(echo "ola" | ./ft_cat hello)
diff -U 3 <(cat nosuchfile 2>&1) <(./ft_cat nosuchfile 2>&1)
diff -U 3 <(cat secret 2>&1) <(./ft_cat secret 2>&1)
echo "OK:D"
```
*/
func main() {
	app := &App{Name: "cat", InStream: os.Stdin, OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(app.Run(os.Args))
}

// Exit Status
const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 1
	ExitCodeConcatError    = 1
)

// App consists of output/error streams.
type App struct {
	Name                 string
	InStream             io.Reader
	OutStream, ErrStream io.Writer
	bufOut, bufErr       *bufio.Writer
	args                 Args
}

// Container for command line args/flags
type Args struct {
	filenames []string
}

// Run is the entry point to the cli. Parses the arguments slice and routes to the proper flag/args combination
func (a *App) Run(args []string) int {
	defer a.flush()
	a.bufOut = bufio.NewWriter(a.OutStream)
	a.bufErr = bufio.NewWriter(a.ErrStream)
	if err := a.parseFlags(args); err != nil {
		return ExitCodeParseFlagError
	}

	if len(a.args.filenames) == 0 {
		return a.concatInputStream()
	} else {
		return a.concatFiles()
	}
}

// Parse flags/args
func (a *App) parseFlags(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(a.ErrStream)
	fs.Usage = func() {
		a.bufErr.WriteString("usage: " + a.Name + " [file ...]\n")
	}
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	a.args.filenames = fs.Args()
	return nil
}

// Concat input stream
func (a *App) concatInputStream() int {
	status := ExitCodeOK
	if _, err := a.bufOut.ReadFrom(a.InStream); err != nil {
		a.printError(err)
		status = ExitCodeConcatError
	}
	return status
}

// Concat files
func (a *App) concatFiles() int {
	status := ExitCodeOK
	for _, file := range a.args.filenames {
		if f, err := os.Open(file); err != nil {
			a.printError(err)
			status = ExitCodeConcatError
		} else if _, err := a.bufOut.ReadFrom(f); err != nil {
			a.printError(err)
			status = ExitCodeConcatError
		}
		a.flush()
	}
	return status
}

// Print error message in cat command style
// 1. Replace 'open' with 'cat'
// 2. Capitalize error description after filename
//
// (examples)
// err.Error()       : "open hoge: no such file"
// cat error message : "cat hoge: No such file"
//
// err.Error()       : "open hoge: permission denied"
// cat error message : "cat hoge: Permission denied"
func (a *App) printError(err error) {
	s := err.Error()
	z := []string{a.Name + ":"}
	z = append(z, strings.SplitN(s, " ", 3)[1:]...)
	z[2] = capitalize(z[2])
	s = strings.Join(z, " ") + "\n"
	a.bufErr.WriteString(s)
}

// Capitalize a string
//
// input  : "no such file"
// output : "No such file"
func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// Flushes all buffers to out/err stream
func (a *App) flush() {
	a.bufOut.Flush()
	a.bufErr.Flush()
}
