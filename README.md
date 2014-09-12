hml
===

Tools to validate `Higgs-ML` challenge submissions.

- [hml](#user-content-hml)
	- [Layout of a zip submission file](#user-content-layout-of-a-zip-submission-file)
		- [zip file generation](#user-content-zip-file-generation)
	- [hml-validate](#user-content-hml-validate)
		- [Install hml-validate](#user-content-install-hml-validate)
		- [Run hml-validate](#user-content-run-hml-validate)
		- [hml-validate help](#user-content-hml-validate-help)
		- [Example](#user-content-example)


## Layout of a `zip` submission file

`zip` submission files **SHALL** have the following structure:

```
Archive:  ./my-team.zip
 extracting: my-team/LICENSE
 extracting: my-team/README
 extracting: my-team/higgsml-build
 extracting: my-team/higgsml-run
 extracting: my-team/higgsml-train
 extracting: my-team/extra/stuff
```

ie: created from:

```
my-team
|-- LICENSE
|-- README
|-- higgsml-build    (optional)
|-- higgsml-run
|-- higgsml-train
`-- extra            (optional)
    `-- stuff        (optional)
```

where:
- `higgsml-build` (with that exact spelling) is the executable script
  (or binary) which builds the code from sources (OPTIONAL)
- `higgsml-run` (with that exact spelling) is the executable script
  (or binary) which runs the prediction
- `higgsml-train` (with that exact spelling) is the executable script
  (or binary) which runs the training
- `LICENSE` (or `LICENSE.txt`) contains an `OSI`-approved license
- `README` (or `README.md` or `README.txt`) contains some
  documentation about the code.

(on `Windows (TM)`, the scripts should be called `higgsml-build.bat`,
`higgsml-run.bat` and `higgsml-train.bat`)

In case you submit 2 codes, each code should have its own directory
structure. *e.g.:*

```
Archive:  ./my-team.zip
 extracting: my-team/code-1/LICENSE
 extracting: my-team/code-1/README
 extracting: my-team/code-1/higgsml-run
 extracting: my-team/code-1/higgsml-train

 extracting: my-team/code-2/LICENSE
 extracting: my-team/code-2/README
 extracting: my-team/code-2/higgsml-run
 extracting: my-team/code-2/higgsml-train
```

*i.e.:* created from:

```
my-team
|-- code-1
|   |-- LICENSE
|   |-- README
|   |-- higgsml-run
|   `-- higgsml-train
`-- code-2
    |-- LICENSE
    |-- README
    |-- higgsml-run
    `-- higgsml-train
```

Environment configuration, if needed, should be performed in the
`higgsml-xyz` scripts (*e.g.* setting up `$PYTHONPATH` or
`$LD_LIBRARY_PATH` environment variables.)

If you ship the sources (and not just a binary) the directory
**SHALL** contain a file `higgsml-build` (or `higgsml-build.bat` on
`Windows (TM)`) which can be called with no argument, and will run the
build procedure in-place (*i.e.* not out of tree) so the `higgsml-run`
and `higgsml-train` can find all the necessary assets at runtime.
The `higgsml-build` script should not try to fetch additional sources
from outside the code directory (*i.e.* no outbound connection allowed.)

### `zip` file generation

A tool is available to generate a `zip` file according to the above
rules:
 [hml-mk-zip](https://github.com/sbinet/hml/blob/master/hml-mk-zip/main.go)

```sh
$ go get github.com/sbinet/hml/hml-mk-zip
$ hml-mk-zip my-team.zip my-team

deflating: my-team/code-1/LICENSE
deflating: my-team/code-1/README.md
deflating: my-team/code-1/go-higgsml
deflating: my-team/code-1/higgsml-run
deflating: my-team/code-1/higgsml-train
deflating: my-team/code-2/LICENSE
deflating: my-team/code-2/README.md
deflating: my-team/code-2/higgsml-run
deflating: my-team/code-2/higgsml-train
deflating: my-team/code-2/higgsml_simplest_v2.py
```

### Example

`github.com/sbinet/hml` has a couple of testcases (with `python` and
`go` programs).
Let's try to make a proper `zip` file:

```sh
$ cd /somewhere
$ git clone git://github.com/sbinet/hml
$ cd hml/testdata
$ hml-mk-zip team-3.zip team-3
2014/09/12 09:33:22 deflating: team-3/LICENSE
2014/09/12 09:33:22 deflating: team-3/README.md
2014/09/12 09:33:22 deflating: team-3/higgsml-run
2014/09/12 09:33:22 deflating: team-3/higgsml-simplest.py
2014/09/12 09:33:22 deflating: team-3/higgsml-train
```

## `hml-validate`

`hml-validate` is a tool to validate the content of a `zip` submission
file, to make sure that `zip` file will be usable by `HEP` physicists.

`hml-validate` will run off a `zip` submission file.

If the `zip` submission file contains a `higgsml-build` script, it
will be run prior to anything else, to generate the needed binaries
and assets.

`higgsml-build` **SHALL** be called with no argument.


`hml-validate` will then run (when instructed to do so by the `-train`
switch from the command line):
 
 ```sh
$ higgsml-train training.csv trained.dat
 ```

to create a `trained.dat` file from the `training.csv` sample.

Finally, `hml-validate` will run:

```sh
$ higgsml-run test.csv trained.dat submission.csv
```

to create the `submission.csv` file from the test sample and the
training parameters.

When everything is successful, it will collect the results (the
`submission.csv` for each code) under a new `higgsml-output`
directory.

### Install `hml-validate`

```sh
$ go get github.com/sbinet/hml/hml-validate
```

### Run `hml-validate`

```sh
$ hml-validate my-team.zip
::: higgsml-validate...
Archive: my-team.zip
  inflating: my-team/README.md
  inflating: my-team/higgsml-run
  inflating: my-team/higgsml-train

::: run prediction...
::: higgs-ml [prediction]...
::: args: test.csv trained.dat scores_test.csv
::: compute the score for the test file entries [test.csv]
::: loop again on test file to load BDT score pairs
::: sort on the score
::: build a map key=id, value=rank
::: you can now submit [scores_test.csv] to Kaggle website
::: timing: 15.862965585s
::: bye.
::: higgs-ml [prediction]... [ok]
::: run prediction... [ok] (delta=15.903641153s)
```

### `hml-validate` help

```sh
$ hml-validate -help
::: higgsml-validate...
Usage of hml-validate:
 hml-validate zipfile [test.csv [training.csv]]
  -train=false: run the training
```

### Example

An example of the expected `zip` file's content (and directory layout)
can be found [here](https://github.com/sbinet/hml/tree/master/testdata/team-3)
