package main

var (
	Def = struct {
		BuildScript string
		RunScript   string
		TrainScript string
		TrainedName string
		Results     string

		WorkDir string
	}{
		BuildScript: "higgsml-build" + binExe,
		RunScript:   "higgsml-run" + binExe,
		TrainScript: "higgsml-train" + binExe,
		TrainedName: "trained.dat",
		Results:     "submission.csv",

		WorkDir: ".higgsml-work",
	}
)
