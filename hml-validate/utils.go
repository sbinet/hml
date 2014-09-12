package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func path_exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func copytree(dstdir, srcdir string) error {
	var err error

	if !path_exists(dstdir) {
		err = os.MkdirAll(dstdir, 0755)
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(srcdir, func(path string, info os.FileInfo, err error) error {
		rel := ""
		rel, err = filepath.Rel(srcdir, path)
		out := filepath.Join(dstdir, rel)
		fmode := info.Mode()
		if fmode.IsDir() {
			err = os.MkdirAll(out, fmode.Perm())
			if err != nil {
				return err
			}
		} else if fmode.IsRegular() {
			dst, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR, fmode.Perm())
			if err != nil {
				return err
			}
			defer func() { err = dst.Close() }()
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		} else if (fmode & os.ModeSymlink) != 0 {
			rlink, err := os.Readlink(path)
			if err != nil {
				return err
			}
			err = os.Symlink(rlink, out)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unhandled mode (%v) for path [%s]", fmode, path)
		}
		return nil
	})
	return err
}

func copyfile(dstname, srcname string) error {
	src, err := os.Open(srcname)
	if err != nil {
		return err
	}
	defer src.Close()

	fi, err := src.Stat()
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(dstname, os.O_CREATE|os.O_WRONLY, fi.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	err = dst.Close()
	if err != nil {
		return err
	}

	return err
}
