hml
===

Tools to validate `Higgs-ML` challenge submissions.

## Layout of a `zip` submission file

`zip` submission files **SHALL** have the following structure:

```
Archive:  ./my-team.zip
 extracting: my-team/LICENSE
 extracting: my-team/README
 extracting: my-team/higgsml-run
 extracting: my-team/higgsml-train
 extracting: my-team/extra/stuff
```

ie: created from:

```
my-team
|-- LICENSE
|-- README
|-- higgsml-run
|-- higgsml-train
`-- extra
    `-- stuff
```

where:
- `higgsml-run` (with that exact spelling) is the executable script
  (or binary) which runs the prediction
- `higgsml-train` (with that exact spelling) is the executable script
  (or binary) which runs the training
- `LICENSE` (or `LICENSE.txt`) contains an `OSI`-approved license
- `README` (or `README.md` or `README.txt`) contains some
  documentation about the code.

(on `Windows (TM)`, the scripts should be called `higgsml-run.bat` and
`higgsml-train.bat`)

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

### `zip` file generation

A tool is available to generate a `zip` file according to the above
rules:
 [hml-mk-zip](https://github.com/sbinet/hml/blob/master/hml-mk-zip/main.go)

```sh
$ go get github.com/sbinet/hml/hml-mk-zip
$ hml-mk-zip /path/to/mk-zip.go my-team >| my-team.zip

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

## `hml-validate`

`hml-validate` is a tool to validate the content of a `zip` submission
file, to make sure that `zip` file will be usable by `HEP` physicists.

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
  inflating: my-team/higgsml-pred
  inflating: my-team/higgsml-training

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
 hml-validate zipfile [[training.csv] test.csv]
  -train=false: run the training
```

### Example

An example of the expected `zip` file's content (and directory layout)
can be found [here](https://github.com/sbinet/hml/tree/master/testdata/team-3)
