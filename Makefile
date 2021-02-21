get:
	go get -t -v -d ./...
.PHONY: get

mocks:
	rm -rf generatedmocks/*
	mockery --all --output="generatedmocks" --keeptree
.PHONY: mocks

test:
	go test -coverprofile=cover.out `go list ./... | grep -v ./generatedmocks`
.PHONY: test

testv:
	go test -v -coverprofile=cover.out `go list ./... | grep -v ./generatedmocks`
.PHONY: test

test-ci:
	go test -covermode=count -coverprofile=coverage.out `go list ./... | grep -v ./generatedmocks`
.PHONY: test

lint:
	golangci-lint run
.PHONY: lint

fmt:
	gofmt -s -w .
.PHONY: fmt

tidy:
	go mod tidy
.PHONY: tidy

precommit:
	make get
	make mocks
	make tidy
	make test
	make lint
	make fmt
.PHONY: precommit
