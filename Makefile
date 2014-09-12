## simple makefile to log workflow
.PHONY: all test test-all clean build install gen

GOFLAGS ?= $(GOFLAGS:)

all: install


build:
	@go build $(GOFLAGS) ./...

install:
	@go get github.com/sbinet/go-higgsml
	@go get $(GOFLAGS) ./...

gen:
	(cd testdata && ./generate)

test-all: clean install gen
	@go get github.com/sbinet/go-higgsml
	@go test $(GOFLAGS) ./...

	@echo ""
	@echo ""
	@echo "=== team-0 === (train)"
	@hml-validate -train testdata/higgsml-test-team0.zip

	@echo ""
	@echo ""
	@echo "=== team-0 ==="
	@hml-validate testdata/higgsml-test-team0.zip

	@echo ""
	@echo ""
	@echo "=== team-1 === (train)"
	@hml-validate -train testdata/higgsml-test-team1.zip

	@echo ""
	@echo ""
	@echo "=== team-1 ==="
	@hml-validate testdata/higgsml-test-team1.zip

	@echo ""
	@echo ""
	@echo "=== team-2 === (train)"
	@hml-validate -train testdata/higgsml-test-team2.zip

	@echo ""
	@echo ""
	@echo "=== team-2 ==="
	@hml-validate testdata/higgsml-test-team2.zip

	@echo ""
	@echo ""
	@echo "=== team-3 === (train)"
	@hml-validate -train testdata/higgsml-test-team3.zip

	@echo ""
	@echo ""
	@echo "=== team-3 ==="
	@hml-validate testdata/higgsml-test-team3.zip

	@echo ""
	@echo ""
	@echo "=== team-4 === (train)"
	@hml-validate -train testdata/higgsml-test-team4.zip

	@echo ""
	@echo ""
	@echo "=== team-4 ==="
	@hml-validate testdata/higgsml-test-team4.zip

test: clean install gen
	@go get github.com/sbinet/go-higgsml
	@go test $(GOFLAGS) ./...

	@echo ""
	@echo ""
	@echo "=== team-0 ==="
	@hml-validate testdata/higgsml-test-team0.zip

	@echo ""
	@echo ""
	@echo "=== team-1 ==="
	@hml-validate testdata/higgsml-test-team1.zip

	@echo ""
	@echo ""
	@echo "=== team-4 ==="
	@hml-validate testdata/higgsml-test-team4.zip

clean:
	@go clean $(GOFLAGS) -i ./...
	@rm -rf higgsml-output

## EOF
