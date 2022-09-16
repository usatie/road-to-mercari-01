package imgconv_test

import (
	"fmt"
	"imgconv"
	"os"
	"testing"
)

const (
	outputDir   = "../output"
	testdataDir = "../testdata"
)

func cleanupOutputDir(b *testing.B) {
	b.Helper()

	if err := os.RemoveAll(outputDir); err != nil {
		b.Fatal("err %w", err)
	}
	if err := os.Mkdir(outputDir, 0755); err != nil {
		b.Fatal("err %w", err)
	}
}

func Test(t *testing.T) {
	t.Parallel()
}

// Benchmarks
func BenchmarkConvert_Jpg2Png(b *testing.B) {
	app := &imgconv.App{Output: os.Stdout, ErrOutput: nil}
	for i := 0; i < 10; i++ {
		app.Run([]string{"./convert", "-d", b.TempDir(), testdataDir + "/jpgs"})
	}
}

func BenchmarkConvert_Png2Jpg(b *testing.B) {
	app := &imgconv.App{Output: os.Stdout, ErrOutput: nil}
	for i := 0; i < 10; i++ {
		app.Run([]string{"./convert", "-d", b.TempDir(), "-i=png", "-o=jpg", testdataDir + "/pngs"})
	}
}

// Error cases
// Root path is empty string
func ExampleApp_Run_Error_EmptyRootPath(t *testing.T) {
	app := &imgconv.App{Output: nil, ErrOutput: os.Stdout}
	fmt.Println(app.Run([]string{"./convert", ""}))
	// Output:
	// error: invalid argument
	// 1
}

// Root path is missing
func ExampleApp_Run_Error_MissingRootPath() {
	app := &imgconv.App{Output: nil, ErrOutput: os.Stdout}
	fmt.Println(app.Run([]string{"./convert"}))
	// Output:
	// error: invalid argument
	// 1
}

// Root path no such file or directory
func ExampleApp_Run_Error_NoSuchFileOrDir() {
	app := &imgconv.App{Output: nil, ErrOutput: os.Stdout}
	fmt.Println(app.Run([]string{"./convert", "nosuchdir"}))
	// Output: error: nosuchdir: no such file or directory
	// 1
}

// Path contains file(s) which are not image
func ExampleApp_Run_Error_InvalidFile() {
	app := &imgconv.App{Output: nil, ErrOutput: os.Stdout}
	fmt.Println(app.Run([]string{"./convert", "../testdata/texts"}))
	// Output:
	// error: ../testdata/texts/hello.txt is not a valid file
	// 1
}

// Path contains file(s) whose permission is bad
func ExampleApp_Run_Error_PermissionDenied() {
	app := &imgconv.App{Output: nil, ErrOutput: os.Stdout}
	fmt.Println(app.Run([]string{"./convert", "../testdata/secrets/bear.jpg"}))
	// Output:
	// error: open ../testdata/secrets/bear.jpg: permission denied
	// 1
}
