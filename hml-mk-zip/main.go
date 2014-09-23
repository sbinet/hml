package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, " %s out.zip file1 [file2 [dir1]...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nNeed more help ?\n => https://www.kaggle.com/c/higgs-boson/forums\n")
	}

	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "**error** %s needs at least 2 arguments.\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	out, err := os.Create(flag.Arg(0))
	if err != nil {
		log.Fatalf("could not create file [%s]: %v\n", flag.Arg(0), err)
	}
	defer out.Close()

	// Create a new zip archive.
	w := zip.NewWriter(out)

	// Add some files to the archive.
	files := make([]string, 0)
	for _, path := range flag.Args()[1:] {

		fi, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !fi.IsDir() {
			files = append(files, path)
			continue
		}
		// recurse
		err = filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				return nil
			}
			files = append(files, path)
			return err
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, fname := range files {
		log.Printf("deflating: %s\n", fname)
		r, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		fi, err := r.Stat()
		if err != nil {
			log.Fatal(err)
		}

		hdr, err := zip.FileInfoHeader(fi)
		if err != nil {
			log.Fatal(err)
		}

		// keep directory structure
		hdr.Name = fname

		f, err := w.CreateHeader(hdr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(f, r)
		if err != nil {
			log.Fatal(err)
		}

		err = r.Close()
		if err != nil {
			log.Fatal(err)
		}

	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

}
