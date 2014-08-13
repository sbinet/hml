//+build ignore

package main

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
)

func main() {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	files := os.Args[1:]
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

	_, err = io.Copy(os.Stdout, buf)
	if err != nil {
		log.Fatal(err)
	}
}
