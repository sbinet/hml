//+build ignore

package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// Create a new zip archive.
	w := zip.NewWriter(os.Stdout)

	// Add some files to the archive.
	files := make([]string, 0)
	for _, path := range os.Args[1:] {
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
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

}
