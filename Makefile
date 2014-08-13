## simple makefile to log workflow
.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test


build:
	@go build $(GOFLAGS) ./...

install:
	@go get $(GOFLAGS) ./...

test: install
	@go test $(GOFLAGS) ./...
	(cd testdata && ./generate)
	@echo ""
	@echo ""
	@echo "=== team-1 ==="
	@hml-validate testdata/higgsml-test-team1.zip

	@echo ""
	@echo ""
	@echo "=== team-2 ==="
	@hml-validate testdata/higgsml-test-team2.zip

	@/bin/rm -rf testdata/team-2/go-higgsml

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
