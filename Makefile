## simple makefile to log workflow
.PHONY: all test clean build install gen

GOFLAGS ?= $(GOFLAGS:)

all: install test


build:
	@go build $(GOFLAGS) ./...

install:
	@go get $(GOFLAGS) ./...

gen:
	(cd testdata && ./generate)

test: install gen
	@go test $(GOFLAGS) ./...
	@echo ""
	@echo ""
	@echo "=== team-1 ==="
	@hml-validate -train testdata/higgsml-test-team1.zip

	@echo ""
	@echo ""
	@echo "=== team-1 ==="
	@hml-validate testdata/higgsml-test-team1.zip

	@echo ""
	@echo ""
	@echo "=== team-2 ==="
	@hml-validate -train testdata/higgsml-test-team2.zip

	@echo ""
	@echo ""
	@echo "=== team-2 ==="
	@hml-validate testdata/higgsml-test-team2.zip

	@/bin/rm -rf testdata/team-2/go-higgsml

	@echo ""
	@echo ""
	@echo "=== team-3 ==="
	@hml-validate -train testdata/higgsml-test-team3.zip

	@echo ""
	@echo ""
	@echo "=== team-3 ==="
	@hml-validate testdata/higgsml-test-team3.zip


clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
