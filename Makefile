.DEFAULT_GOAL := build

export GO111MODULE=on
export CGO_ENABLED=0
export BINARY=jelliflix
export BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Runs the linter
	$(GOPATH)/bin/golangci-lint run

.PHONY: test
test: ## Run the test suite
	CGO_ENABLED=1 go test -race -coverprofile="coverage.txt" ./...

.PHONY: clean
clean: ## Remove the binary
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi

.PHONY: build
build: ## Build the binary
	go build -a -gcflags='-N -l' -installsuffix cgo -o $(BINARY)

.PHONY: packed
packed: ## Build a packed version of the binary
	build
	upx --best --lzma $(BINARY)

.PHONY: docker
docker: ## Build the docker image with packed binaries
	docker build -t $(BINARY):latest -t $(BINARY):$(BUILD) -f Dockerfile .
