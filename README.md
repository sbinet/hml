hml
===

Tools to validate `Higgs-ML` challenge submissions.

## Test a `zip` submission file

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
- `higgsml-run` is the script (or binary) which runs the prediction
- `higgsml-train` is the script (or binary) which runs the training


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
|   |-- LICENSE
    |-- README
    |-- higgsml-run
    `-- higgsml-train
```

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
