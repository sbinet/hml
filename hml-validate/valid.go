package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Validator struct {
	root string // work dir
	code []Code // list of submissions
}

func NewValidator(dir string, train bool) (Validator, error) {
	var err error
	v := Validator{
		root: dir,
		code: make([]Code, 0, 2),
	}

	// find all 'higgsml-pred' executables
	exes := make([]string, 0)
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(path), runscript) {
			exes = append(exes, path)
		}
		return err
	})

	if len(exes) <= 0 {
		return v, fmt.Errorf("hml: could not find any suitable executable in zip-file")
	}

	if len(exes) > 2 {
		return v, fmt.Errorf("hml: too many higgsml-pred executables (got=%d, max=%d)", len(exes), 2)
	}

	for _, exe := range exes {
		dir := filepath.Dir(exe)
		code, err := NewCode(dir, train)
		if err != nil {
			return v, err
		}
		code.Name = dir[len(v.root):]
		if code.Name[0] == '/' {
			code.Name = code.Name[1:]
		}
		v.code = append(v.code, code)
	}

	return v, err
}

func (v Validator) Run() error {
	var err error

	for i, code := range v.code {
		start := time.Now()
		printf("\n=== code submission #%d/%d (%s)...\n", i+1, len(v.code), code.Name)
		err = code.run()
		if err != nil {
			return err
		}
		printf("=== code submission #%d/%d (%s)... [ok] (delta=%v)\n",
			i+1, len(v.code), code.Name,
			time.Since(start),
		)
	}

	return err
}

type Code struct {
	Root    string // directory containing sources/binaries
	Name    string // name of this code submission (e.g. team/code)
	Train   string // path to training executable
	Pred    string // path to prediction executable
	Trained string // path to trained data parameters

	License string // path to LICENSE file
	Readme  string // path to README file

	DoTraining bool

	trainfile string // path to training.csv file
	testfile  string // path to test.csv file
}

func NewCode(dir string, train bool) (Code, error) {
	var err error

	code := Code{
		Root:       dir,
		Name:       filepath.Base(dir),
		Trained:    trainedname,
		DoTraining: train,
		trainfile:  trainfile,
		testfile:   testfile,
	}

	// FIXME: handle multiple-submissions zipfiles
	//        presumably: 1 directory per submission.

	exes := make([]string, 0)
	// find executables
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(path), "readme") {
			code.Readme = path
		}

		if strings.Contains(strings.ToLower(path), "license") {
			code.License = path
		}

		if strings.Contains(strings.ToLower(path), trainedname) {
			code.Trained = path
		}

		// FIXME: better way ?
		if !strings.Contains(fi.Mode().String(), "x") {
			return nil
		}
		exes = append(exes, path)
		// printf(">>> %s\n", path)
		if strings.HasSuffix(strings.ToLower(path), trainscript) {
			code.Train = path
		}
		if strings.HasSuffix(strings.ToLower(path), runscript) {
			code.Pred = path
		}
		return err
	})

	if len(exes) <= 0 {
		return code, fmt.Errorf("hml: could not find any suitable executable in zip-file")
	}

	if code.Train == "" && code.Pred == "" {
		// take first one
		code.Train = exes[0]
		code.Pred = exes[0]
	}

	if code.Train == "" && code.Pred != "" {
		code.Train = code.Pred
	}

	if code.Train != "" && code.Pred == "" {
		code.Pred = code.Train
	}

	if code.License == "" {
		return code, fmt.Errorf("hml: could not find a LICENSE file under [%s]", code.Root)
	}

	return code, err
}

func (code Code) run() error {
	var err error

	dir := filepath.Join(code.Root, ".higgsml-work")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	if code.DoTraining {
		err = code.run_training(dir)
		if err != nil {
			return err
		}
		printf("\n")
	}

	err = code.run_pred(dir)
	if err != nil {
		return err
	}

	return err
}

func (code Code) run_training(dir string) error {
	var err error
	errch := make(chan error)

	printf("::: run training...\n")

	trained := pdir(dir, trainedname)

	cmd := exec.Command(code.Train, code.trainfile, trained)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = code.Root

	start := time.Now()
	go func() {
		err = cmd.Start()
		if err != nil {
			errch <- err
			return
		}
		errch <- cmd.Wait()
	}()

	duration := 1 * time.Hour
	select {
	case <-time.After(duration):
		cmd.Process.Kill()
		return fmt.Errorf("hml: training timed out (%v)\n", duration)
	case err = <-errch:
		break
	}

	if err != nil {
		printf("::: run training... [ERR] (delta=%v)\n", time.Since(start))
		return err
	}

	printf("::: run training... [ok] (delta=%v)\n", time.Since(start))
	return err
}

func (code Code) run_pred(dir string) error {
	var err error
	errch := make(chan error)

	printf("::: run prediction...\n")

	trained := code.Trained
	if code.DoTraining {
		// if we ran the training, then use that file instead.
		trained = pdir(dir, trainedname)
	}

	results := pdir(dir, csv_results)
	cmd := exec.Command(code.Pred, code.testfile, trained, results)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = code.Root

	start := time.Now()
	go func() {
		err = cmd.Start()
		if err != nil {
			errch <- err
			return
		}
		errch <- cmd.Wait()
	}()

	duration := 1 * time.Hour
	select {
	case <-time.After(duration):
		cmd.Process.Kill()
		return fmt.Errorf("hml: prediction timed out (%v)\n", duration)
	case err = <-errch:
		break
	}

	if err != nil {
		printf("::: run prediction... [ERR] (delta=%v)\n", time.Since(start))
		return err
	}

	printf("::: run prediction... [ok] (delta=%v)\n", time.Since(start))
	return err
}

func pdir(dirs ...string) string {
	return filepath.Join(dirs...)
}
