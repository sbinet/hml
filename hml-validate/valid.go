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
		if strings.Contains(strings.ToLower(path), Def.RunScript) {
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

	// create output dir to collect results
	outdir := "higgsml-output"
	err = os.MkdirAll(outdir, 0755)
	if err != nil {
		return err
	}

	for i, code := range v.code {
		start := time.Now()
		printf("\n=== code submission #%d/%d (%s)...\n", i+1, len(v.code), code.Name)
		err = code.run(outdir)
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
	Build   string // path to code-building executable
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
		Trained:    Def.TrainedName,
		DoTraining: train,
		trainfile:  trainfile,
		testfile:   testfile,
	}

	// FIXME: handle multiple-submissions zipfiles
	//        presumably: 1 directory per submission.

	exes := make([]string, 0)
	// find executables
	find_execs := func() error {
		return filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				return nil
			}
			if strings.Contains(strings.ToLower(path), "readme") {
				code.Readme = path
			}

			if strings.Contains(strings.ToLower(path), "license") {
				code.License = path
			}

			if strings.Contains(strings.ToLower(path), Def.TrainedName) {
				code.Trained = path
			}

			// FIXME: better way ?
			if !strings.Contains(fi.Mode().String(), "x") {
				return nil
			}
			exes = append(exes, path)
			// printf(">>> %s\n", path)
			if strings.HasSuffix(strings.ToLower(path), Def.TrainScript) {
				code.Train = path
			}
			if strings.HasSuffix(strings.ToLower(path), Def.RunScript) {
				code.Pred = path
			}
			if strings.HasSuffix(strings.ToLower(path), Def.BuildScript) {
				code.Build = path
			}
			return err
		})
	}

	err = find_execs()
	if err != nil {
		return code, err
	}

	// whether we need to build the code first.
	if code.Build != "" && *g_build {
		err = code.build(dir)
		if err != nil {
			return code, err
		}
		printf("\n")

		// presumably, train-script and run-script need to be re-discovered...
		exes = exes[:0]
		err = find_execs()
		if err != nil {
			return code, err
		}
	}

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

func (code Code) run(outdir string) error {
	var err error

	dir := filepath.Join(code.Root, Def.WorkDir)
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

	err = code.collect(outdir)
	if err != nil {
		return err
	}

	return err
}

func (code Code) build(dir string) error {
	var err error
	errch := make(chan error)

	printf("::: building code...\n")

	cmd := exec.Command(code.Build)
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
		err = code.kill(cmd.Process)
		err = fmt.Errorf("hml: building code timed out (%v)\nerr=%v\n", duration, err)
	case err = <-errch:
	}

	if err != nil {
		printf("::: building code... [ERR] (delta=%v)\n", time.Since(start))
		return err
	}

	printf("::: building code... [ok] (delta=%v)\n", time.Since(start))
	return err
}

func (code Code) run_training(dir string) error {
	var err error
	errch := make(chan error)

	printf("::: run training...\n")

	trained := pdir(dir, Def.TrainedName)

	cmd := code.command(code.Train, code.trainfile, trained)

	start := time.Now()
	go func() {
		err = cmd.Start()
		if err != nil {
			errch <- err
			return
		}
		errch <- cmd.Wait()
	}()

	duration := *g_traintime
	select {
	case <-time.After(duration):
		err = code.kill(cmd.Process)
		err = fmt.Errorf("hml: training timed out (%v)\nerr=%v\n", duration, err)
	case err = <-errch:
	}

	if err != nil {
		printf("::: run training... [ERR] (delta=%v)\n", time.Since(start))
		return err
	}

	_, err = os.Stat(trained)
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
		trained = pdir(dir, Def.TrainedName)
	}

	results := pdir(dir, Def.Results)
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

	duration := *g_predtime
	select {
	case <-time.After(duration):
		err = code.kill(cmd.Process)
		err = fmt.Errorf("hml: prediction timed out (%v)\nerr=%v\n", duration, err)
	case err = <-errch:
	}

	if err != nil {
		printf("::: run prediction... [ERR] (delta=%v)\n", time.Since(start))
		return err
	}

	_, err = os.Stat(results)
	if err != nil {
		printf("::: run prediction... [ERR] (delta=%v)\n", time.Since(start))
		return err
	}

	printf("::: run prediction... [ok] (delta=%v)\n", time.Since(start))
	return err
}

func (code Code) collect(outdir string) error {
	var err error

	// collect result
	srcdir := pdir(code.Root, Def.WorkDir)
	dstdir := pdir(outdir, code.Name)

	err = os.MkdirAll(dstdir, 0755)
	if err != nil {
		return err
	}

	files := []string{
		pdir(srcdir, Def.Results),
	}

	if code.DoTraining {
		// if we ran the training, then collect the trained data file as well
		files = append(files, pdir(code.Root, Def.WorkDir, Def.TrainedName))
	}

	for _, src := range files {
		dst := pdir(dstdir, filepath.Base(src))
		err = copyfile(dst, src)
		if err != nil {
			return err
		}
	}

	return err
}

func pdir(dirs ...string) string {
	return filepath.Join(dirs...)
}
