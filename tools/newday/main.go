// Binary newday creates directories and files for quickly bootstrapping a new
// puzzle's solution. It is intended to be invoked from the root directory of this repository, like so:
//
//	go run ./tools/newday -year=2017 -day=13
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/golang/glog"

	// For embedding templates.
	_ "embed"
)

var (
	year = flag.Int("year", 2015, "The year.")
	day  = flag.Int("day", 1, "The day.")
	dir  = flag.String("dir", "", "The directory under which files will be written. The files will have the paths <dir>/yearYYYY/dayDD/dayDD{,_test}.go. If this flag is empty the current directory will be used.")
)

var (
	//go:embed dayDD.go.tmpl
	pkgTmplRaw string
	pkgTmpl    = template.Must(template.New("pkg").Parse(pkgTmplRaw))
	//go:embed dayDD_test.go.tmpl
	testTmplRaw string
	testTmpl    = template.Must(template.New("test").Parse(testTmplRaw))
)

type tmplArgs struct {
	Year, Day int
}

func (t tmplArgs) PaddedDay() string {
	return fmt.Sprintf("%02d", t.Day)
}

type target struct {
	Year, Day int
	Dir       string
}

func (t target) packageDirectory() string {
	return filepath.Join(t.Dir, fmt.Sprintf("year%d", t.Year), fmt.Sprintf("day%02d", t.Day))
}

func (t target) tmplArgs() tmplArgs {
	return tmplArgs{
		Year: t.Year,
		Day:  t.Day,
	}
}

func (t target) ensurePackageDirectoryExists() error {
	dir := t.packageDirectory()
	glog.V(1).Infof("Ensuring package directory exists: %q", dir)
	return os.MkdirAll(dir, fs.FileMode(0o755))
}

func (t target) writePackageFile() error {
	var buf bytes.Buffer
	if err := pkgTmpl.Execute(&buf, t.tmplArgs()); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format executed template: %v", err)
	}
	p := filepath.Join(t.packageDirectory(), fmt.Sprintf("day%02d.go", t.Day))
	if err := os.WriteFile(p, formatted, fs.FileMode(0o644)); err != nil {
		return fmt.Errorf("write formatted source: %v", err)
	}
	glog.V(1).Infof("Wrote package file: %q", p)
	return nil
}

func (t target) writeTestFile() error {
	var buf bytes.Buffer
	if err := testTmpl.Execute(&buf, t.tmplArgs()); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format executed template: %v", err)
	}
	p := filepath.Join(t.packageDirectory(), fmt.Sprintf("day%02d_test.go", t.Day))
	if err := os.WriteFile(p, formatted, fs.FileMode(0o644)); err != nil {
		return fmt.Errorf("write formatted source: %v", err)
	}
	glog.V(1).Infof("Wrote test file: %q", p)
	return nil
}

func (t target) WriteFiles() error {
	glog.V(1).Info("Creating directories and writing files...")
	if err := t.ensurePackageDirectoryExists(); err != nil {
		return fmt.Errorf("ensure package directory exists: %v", err)
	}
	if err := t.writePackageFile(); err != nil {
		return fmt.Errorf("write package file: %v", err)
	}
	if err := t.writeTestFile(); err != nil {
		return fmt.Errorf("write test file: %v", err)
	}
	return nil
}

func errmain() error {
	if *year < 2015 {
		return fmt.Errorf("-year=%d is invalid; must be at least 2015", *year)
	}
	if *day < 1 || *day > 25 {
		return fmt.Errorf("-day=%d is invalid; must be in the range [1, 25]", *day)
	}
	outputDir := *dir
	if outputDir == "" {
		d, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("default to current directory: %v", err)
		}
		glog.V(1).Infof("-dir was empty; falling back to current directory: %q.", d)
		outputDir = d
	} else {
		glog.V(1).Infof("Files will be written to directory: %q", outputDir)
	}

	t := target{
		Year: *year,
		Day:  *day,
		Dir:  outputDir,
	}
	if err := t.WriteFiles(); err != nil {
		return fmt.Errorf("write files: %v", err)
	}

	return nil
}

func main() {
	flag.Parse()
	if err := errmain(); err != nil {
		glog.Exit(err)
	}
}
