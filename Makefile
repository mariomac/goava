GOCMD       ?= go
all: fmt lint test

fmt:
	@echo "=== [ fmt ]: formatting code..."
	goimports -w ./

lint:
	@echo "=== [ lint ]: Validating source code running golint..."
	golangci-lint run

test:
	@echo "=== [ test ]: Running unit tests..."
	$(GOCMD) test -mod vendor -coverpkg=./... -coverprofile coverage.out -a -race ./pkg/...

.PHONY: all build clean lint compile test fmt