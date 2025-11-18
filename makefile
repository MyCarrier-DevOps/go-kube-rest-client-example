SHELL:=/bin/bash

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

lint: install-tools
	@echo "Linting all modules..."; \
	golangci-lint run --config .github/.golangci.yml --timeout 5m ./...; 

.PHONY: fmt
fmt:
	@echo "Formatting module..."; \
	golangci-lint fmt --config .github/.golangci.yml ./...;

.PHONY: bump
bump:
	@echo "Bumping module..."; \
	go get -u && go mod tidy; \

.PHONY: tidy
tidy:
	@echo "Tidying module..."; \
	go mod tidy;


.PHONY: install-tools
install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b `go env GOPATH`/bin v2.5.0
