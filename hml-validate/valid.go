package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	runscript   = "higgsml-run"
	trainscript = "higgsml-train"
)

type Validate struct {
	root    string // work dir
	assets  string // submission dir
	Train   string // path of training executable
	Pred    string // path of prediction executable
	Trained string // path to trained data

	Readme string // name of the README file

	DoTraining bool

	trainfile string // path to training.csv file
	testfile  string // path to test.csv file
}

func NewValidate(dir string, train bool) (Validate, error) {
	var err error

	v := Validate{
		root:       dir,
		Trained:    "trained.dat",
		DoTraining: train,
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
			v.Readme = path
		}

		if strings.Contains(strings.ToLower(path), "trained.dat") {
			v.Trained = path
		}

		// FIXME: better way ?
		if !strings.Contains(fi.Mode().String(), "x") {
			return nil
		}
		exes = append(exes, path)
		// printf(">>> %s\n", path)
		if strings.Contains(strings.ToLower(path), "higgsml-train") {
			v.Train = path
		}
		if strings.Contains(strings.ToLower(path), "higgsml-pred") {
			v.Pred = path
		}
		return err
	})

	if len(exes) <= 0 {
		return v, fmt.Errorf("hml: could not find any suitable executable in zip-file")
	}

	if v.Train == "" && v.Pred == "" {
		// take first one
		v.Train = exes[0]
		v.Pred = exes[0]
	}

	if v.Train == "" && v.Pred != "" {
		v.Train = v.Pred
	}

	if v.Train != "" && v.Pred == "" {
		v.Pred = v.Train
	}

	v.assets = filepath.Dir(v.Pred)

	return v, err
}

func (v Validate) Run() error {
	var err error

	// printf("root:   [%s]\n", v.root)
	// printf("assets: [%s]\n", v.assets)

	dir := filepath.Join(v.assets, ".higgsml-work")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	if v.DoTraining {
		printf("\n")
		err = v.run_training(dir)
		if err != nil {
			return err
		}
	}

	printf("\n")
	err = v.run_pred(dir)
	if err != nil {
		return err
	}

	return err
}

func (v Validate) run_training(dir string) error {
	var err error
	errch := make(chan error)

	printf("::: run training...\n")

	cmd := exec.Command(v.Train, v.wdir("training.csv"), pdir(dir, "trained.dat"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = v.assets

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

func (v Validate) run_pred(dir string) error {
	var err error
	errch := make(chan error)

	printf("::: run prediction...\n")

	cmd := exec.Command(v.Pred, v.wdir("test.csv"), v.Trained, pdir(dir, "scores_test.csv"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = v.assets

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

func (v Validate) wdir(fname string) string {
	pwd, err := os.Getwd()
	if err != nil {
		return fname
	}
	return filepath.Join(pwd, fname)
}

func pdir(dirs ...string) string {
	return filepath.Join(dirs...)
}
