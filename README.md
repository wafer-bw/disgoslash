# disgoslash
A Golang serverless Discord slash command application package

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

## Example Project
This small example project uses [Vercel serverless functions](https://vercel.com/docs/serverless-functions/supported-languages#go) to execute and respond to Discord slash commands. You can create a Vercel account [here](https://vercel.com/).

1. Create project directory
    ```sh
    mkdir slashhello
    cd slashhello
    go mod init example.com/slashhello 
    mkdir api
    touch api/index.go
    ```
2. Add the serverless function for Vercel to host and a simple slash command  
    `slashhello/api/index.go`
    ```golang
    package api

    import (
        "net/http"

        "github.com/wafer-bw/disgoslash"
    )

    var command = &discord.ApplicationCommand{
        Name:        "hello",
        Description: "Says hello to the user",
        Options: []*discord.ApplicationCommandOption{
            {
                Type:        discord.ApplicationCommandOptionTypeString,
                Name:        "Name",
                Description: "Enter your name",
                Required:    true,
            },
        },
    }

    func hello(request *discord.InteractionRequest) *discord.InteractionResponse {
        return &discord.InteractionResponse{
            Type: discord.InteractionResponseTypeChannelMessageWithSource,
            Data: &discord.InteractionApplicationCommandCallbackData{
                Content: "Hello " + request.Data.Options[0].Value + "!",
            },
        }
    }

    var slashCommand = disgoslash.NewSlashCommand(command, hello, true, "YOUR_GUILD_ID")

    func Handler(w http.ResponseWriter, r *http.Request) {
        handler := &disgoslash.Handler{
            SlashCommandMap: disgoslash.NewSlashCommandMap(slashCommand),
            Creds: &discord.Credentials{
                PublicKey: "YOUR_APPLICATION_PUBLIC_KEY",
                ClientID:  "YOUR_APPLICATION_CLIENT_ID",
                Token:     "YOUR_APPLICATION_BOT_TOKEN",
            },
        }
        handler.Handle(w, r)
    }
    ```
3. Deploy to vercel
    ```sh
    npm i -g vercel
    vercel --prod
    # Follow the prompts. It will look something like this when it finishes:
    # > Vercel CLI 21.2.3
    # > ✅  Production: https://yourprojectnamehere.vercel.app [copied to clipboard] [32s]
    ```
4. On your Discord application's General Information page set the "Interactions Endpoint URL":
    ```
    https://yourprojectname.vercel.app/api
    ```
5. Go into your Discord server and run the slash command:
    ```
    /hello Bob
    ```
    Your bot should respond:
    > Hello Bob!


## References
* [discordgo](https://github.com/bwmarrin/discordgo) - ed25519 auth
