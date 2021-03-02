# disgoslash
A Golang serverless Discord slash command application library

![tests](https://github.com/wafer-bw/disgoslash/workflows/tests/badge.svg)
![lint](https://github.com/wafer-bw/disgoslash/workflows/lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/wafer-bw/disgoslash)](https://goreportcard.com/report/github.com/wafer-bw/disgoslash)
![CodeQL](https://github.com/wafer-bw/disgoslash/workflows/CodeQL/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/wafer-bw/disgoslash/badge.svg)](https://coveralls.io/github/wafer-bw/disgoslash)
[![Go Reference](https://pkg.go.dev/badge/github.com/wafer-bw/disgoslash.svg)](https://pkg.go.dev/github.com/wafer-bw/disgoslash)

```sh
go get github.com/wafer-bw/disgoslash
```

## Getting Started

### Prerequisites
* [Golang](https://golang.org/dl/)
* [Discord](https://discord.com/)
    * Server `Guild ID`
* [Discord Application w/ Bot](https://discord.com/developers/applications)
    * Application `Public Key`
    * Application `Client ID`
    * Bot `Token`

## TODOs
* readme
    * table of contents
    * example project
* package
    * stable version release
* `client`
    * EditGlobalApplicationCommand
    * EditGuildApplicationCommand
* `syncer`
    * attempt to validate commands with some regex
* `models`
    * finish Guild Model
* import/export
    * `exporter`
        * export commands to json
    * `cmd/disgoslash.go`
        * list commands
        * unregister commands
        * register commands


## References
* [discordgo](https://github.com/bwmarrin/discordgo) - ed25519 auth
