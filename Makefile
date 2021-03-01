get:
	go get -t -v -d ./...
.PHONY: get

mocks:
	mockery --name "clientInterface" --structname "mockClientInterface" --filename "client_mock.go" --inpackage
.PHONY: mocks

test:
	go test -coverprofile=cover.out
.PHONY: test

testv:
	go test -v -coverprofile=cover.out
.PHONY: test

test-ci:
	go test -covermode=count -coverprofile=coverage.out
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

doc:
	godoc -http=:6060 -play
.PHONY: doc

precommit:
	make get
	make mocks
	make tidy
	make test
	make lint
	make fmt
.PHONY: precommit
