# disgoslash
A Golang serverless Discord slash command helper library written for Vercel Serverless Functions

![tests](https://github.com/wafer-bw/disgoslash/workflows/tests/badge.svg)
![lint](https://github.com/wafer-bw/disgoslash/workflows/lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/wafer-bw/disgoslash)](https://goreportcard.com/report/github.com/wafer-bw/disgoslash)
![CodeQL](https://github.com/wafer-bw/disgoslash/workflows/CodeQL/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/wafer-bw/disgoslash/badge.svg)](https://coveralls.io/github/wafer-bw/disgoslash)

## Getting Started

### Prerequisites
* [Golang](https://golang.org/dl/)
* [Vercel](https://vercel.com/)
* [Discord Server](https://discord.com/)
    * Server `Guild ID`
* [Discord Application & Bot](https://discord.com/developers/applications)
    * Application `Public Key`
    * Application `Client ID`
    * Bot `Token`

## TODOs
* package
    * exported comment update pass
    * coverage pass
    * todo pass
    * `doc.go`
    * header comments
    * examples
    * stable version release
* usage
    * add usage section to readme
* `client`
    * cleanup struct usage
    * add tests for retry code
    * add tests for uncovered code
    * EditGlobalApplicationCommand
    * EditGuildApplicationCommand
* `models`
    * finish Guild Model
* `exporter`
    * export commands to json
* `cmd/disgoslash.go`
    * list commands
    * unregister commands
    * register commands


## References
* [discordgo](https://github.com/bwmarrin/discordgo) - ed25519 auth
