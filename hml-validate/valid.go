package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Validate struct {
	root  string // workdir
	Train string // name of training executable
	Prod  string // name of production executable

	Readme string // name of the README file
}

func NewValidate(dir string) (Validate, error) {
	var err error

	v := Validate{root: dir}

	exes := make([]string, 0)
	// find executables
	err = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		if strings.Contains(strings.ToLower(path), "readme") {
			v.Readme = path
		}

		// FIXME: better way ?
		if !strings.Contains(fi.Mode().String(), "x") {
			return nil
		}
		exes = append(exes, path)
		// printf(">>> %s\n", path)
		if strings.Contains(strings.ToLower(path), "train") {
			v.Train = path
		}
		if strings.Contains(strings.ToLower(path), "prod") {
			v.Prod = path
		}
		return err
	})

	if len(exes) <= 0 {
		return v, fmt.Errorf("hml: could not find any suitable executable in zip-file")
	}

	if v.Train == "" && v.Prod == "" {
		// take first one
		v.Train = exes[0]
		v.Prod = exes[0]
	}

	return v, err
}

func (v Validate) Run() error {
	var err error

	printf("\n")
	err = v.run_training()
	if err != nil {
		return err
	}

	printf("\n")
	err = v.run_pred()
	if err != nil {
		return err
	}

	return err
}

func (v Validate) run_training() error {
	var err error
	printf("::: run training...\n")
	dir := filepath.Join(v.root, "hml-train")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command(v.Train, "training.csv", "trained.dat")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = dir

	errch := make(chan error)
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
		printf("::: run training... [ERR]\n")
		break
	}

	if err != nil {
		printf("::: run training... [ERR]\n")
		return err
	}

	printf("::: run training... [ok]\n")
	return err
}

func (v Validate) run_pred() error {
	var err error
	printf("::: run prediction...\n")
	dir := filepath.Join(v.root, "hml-prediction")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command(v.Prod, "test.csv", "trained.dat", "scores_test.csv")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Dir = dir

	errch := make(chan error)
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
		printf("::: run prediction... [ERR]\n")
		return err
	}

	printf("::: run prediction... [ok]\n")
	return err
}
