hml
===

Tools to validate `Higgs-ML` challenge submissions.

## Test a `zip` submission file

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
