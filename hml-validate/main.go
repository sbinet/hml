package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	g_train = flag.Bool("train", false, "switch to run the training during validation")

	trainfile = "training.csv"
	testfile  = "test.csv"
)

func printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, args...)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: hml-validate zipfile-or-directory [test.csv [training.csv]]

where:
  - test.csv is a test file (taken from $PWD if not given.)
  - training.csv is a training file (taken from $PWD if not given.)
    training.csv is needed iff -train is enabled.

ex:
 $ hml-validate my-team
 $ hml-validate my-team/
 $ hml-validate my-team.zip
 $ hml-validate my-team.tar.gz
 $ hml-validate my-team /path/to/test.csv
 $ hml-validate my-team /path/to/test.csv /path/to/training.csv

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	printf("::: higgsml-validate...\n")
	if flag.NArg() < 1 {
		printf("**error** hml-validate needs the path to a zip file\n")
		flag.Usage()
		os.Exit(1)
	}

	if !*g_train {
		trainfile = ""
	}

	if flag.NArg() > 1 {
		testfile = flag.Arg(1)
	}
	if flag.NArg() > 2 {
		testfile = flag.Arg(1)
		trainfile = flag.Arg(2)
	}

	for _, file := range []*string{&testfile, &trainfile} {
		// printf("file[%d]=%q\n", i, *file)
		if *file == "" {
			continue
		}
		if !filepath.IsAbs(*file) {
			name, err := filepath.Abs(*file)
			if err == nil {
				*file = name
			}
		}
		_, err := os.Lstat(*file)
		if err != nil {
			printf("**error** no such file [%s]\n", *file)
			flag.Usage()
			os.Exit(1)
		}
		// printf("file[%d]=%q\n", i, *file)
	}

	rc := run()
	os.Exit(rc)
}

func run() int {

	dir, err := ioutil.TempDir("", "higgsml-validate-")
	if err != nil {
		printf("**error** creating tmpdir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(dir)

	fname := flag.Arg(0)

	fi, err := os.Lstat(fname)
	if err != nil {
		printf("**error** stat-ing [%s]: %v\n", fname, err)
		return 1
	}

	if fi.IsDir() {
		err = process_dir(dir, fname)
		if err != nil {
			printf("**error** processing directory [%s]: %v\n", fname, err)
			return 1
		}
	} else {
		err = process_archive(dir, fname)
		if err != nil {
			printf("**error** processing archive [%s]: %v\n", fname, err)
			return 1
		}
	}

	v, err := NewValidator(dir, *g_train)
	if err != nil {
		printf("**error** validating: %v\n", err)
		return 1
	}

	// printf("validate: %#v\n", v)
	err = v.Run()
	if err != nil {
		printf("**error** running validation: %v\n", err)
		return 1
	}

	return 0
}

func process_dir(outdir, indir string) error {
	return copytree(filepath.Join(outdir, filepath.Base(indir)), indir)
}

func process_archive(dir, fname string) error {
	var err error
	printf("Archive: %s\n", fname)

	r, err := zip.OpenReader(fname)
	if err != nil {
		printf("**error** %v\n", err)
		return err
	}
	defer r.Close()

	// printf("comment: %q\n", r.Comment)

	err = unzip(dir, r)
	if err != nil {
		printf("**error** unzipping [%s] under [%s]: %v\n", fname, dir, err)
		return err
	}
	return err
}

func unzip(tmpdir string, r *zip.ReadCloser) error {

	var err error

	for _, f := range r.File {
		printf("  inflating: %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			printf("**error** opening [%s]: %v\n", f.Name, err)
			return err
		}
		defer rc.Close()

		ofname := filepath.Join(tmpdir, f.Name)
		// printf("extracting into: [%s] (%v)\n", ofname, f.Mode())
		odir := filepath.Dir(ofname)
		err = os.MkdirAll(odir, 0755)
		if err != nil {
			printf("**error** creating output dir [%s]: %v\n", odir, err)
			return err
		}
		w, err := os.Create(ofname)
		if err != nil {
			printf("**error** creating output file [%s]: %v\n", ofname, err)
			return err
		}
		defer func(w *os.File) {
			err = w.Close()
			if err != nil {
				printf("**error** closing output file [%s]: %v\n", ofname, err)
				os.Exit(1)
			}
		}(w)
		_, err = io.Copy(w, rc)
		if err != nil {
			printf("**error** copying to [%s]: %v\n", ofname, err)
			return err
		}

		err = w.Chmod(f.Mode())
		if err != nil {
			printf("**error** setting chmod [%v] for file [%s]: %v\n", f.Mode(), ofname, err)
			return err
		}
	}

	return err
}
