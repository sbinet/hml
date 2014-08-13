hml
===

Tools to validate `Higgs-ML` challenge submissions.

## Test a `zip` submission file

### Install `hml-validate`

```sh
$ go get github.com/sbinet/hml/...
```

### Run `hml-validate`

```sh
$ hml-validate my-team.zip
::: higgsml-validate...
Archive: my-team.zip
  inflating: my-team/README.md
  inflating: my-team/higgsml-pred
  inflating: my-team/higgsml-training

::: run training...
::: higgs-ml [training]...
::: args: training.csv trained.dat
::: higgs-ml [training]... [ok]

::: run pred...
::: higgs-ml [prediction]...
::: args: test.csv trained.dat scores_test.csv
::: higgs-ml [prediction]... [ok]
```
