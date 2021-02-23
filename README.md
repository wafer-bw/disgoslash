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
* [Discord](https://discord.com/)
* [Discord Application](https://discord.com/developers/applications)

### Usage
todo

## TODOs
* `client`
    * EditGlobalApplicationCommand
    * EditGuildApplicationCommand
* `models`
    * Finish Guild Model
* Package
    * Remove interfaces and mocks where applicable
        * Determine what to do with `client`
    * Reduce lengths of names of `discord` types and prefix with `Discord`
    * Exported comment update pass
    * Coverage pass
    * todo pass
    * `doc.go`
    * Header comments
    * Examples
    * Stable Version Release

## References
* [discordgo](https://github.com/bwmarrin/discordgo) - ed25519 auth
