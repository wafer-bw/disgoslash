# disgoslash
A serverless Discord slash command bot powered by Vercel written in Golang

![tests](https://github.com/wafer-bw/disgoslash/workflows/tests/badge.svg)
![lint](https://github.com/wafer-bw/disgoslash/workflows/lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/wafer-bw/disgoslash)](https://goreportcard.com/report/github.com/wafer-bw/disgoslash)
![CodeQL](https://github.com/wafer-bw/disgoslash/workflows/CodeQL/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/wafer-bw/disgoslash/badge.svg)](https://coveralls.io/github/wafer-bw/disgoslash)

## Getting Started

### Prerequisites
#### Primary
* [Golang](https://golang.org/dl/)
* [Vercel](https://vercel.com/)
* [Discord](https://discord.com/)
* [Discord Application](https://discord.com/developers/applications)

#### Dev
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
* [mockery](https://github.com/vektra/mockery)

### Setup
    ```sh
    git clone git@github.com:wafer-bw/disgoslash.git
    go get -t -v -d ./...
    ```

### Usage (POSIX)
```sh
# Get Dependencies
make get
# Tidy go.mod
make tidy
# Run tests
make test
# Run verbose tests
make testv
# Run linting
make lint
# Run formatting
make fmt
# Regenerate mocks
make mocks
# Run all the things you should before you make a commit
make precommit
```

### Developing Slash Commands
todo

## TODOs
* `client`
    * EditGlobalApplicationCommand
    * EditGuildApplicationCommand
* `models`
    * Finish Guild Model

## References
* [discordgo](https://github.com/bwmarrin/discordgo) - ed25519 auth
