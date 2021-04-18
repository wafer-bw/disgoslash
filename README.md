# disgoslash
A Golang package that provides an easy way to create serverless Discord slash command applications

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
This small example project uses [Vercel serverless functions](https://vercel.com/docs/serverless-functions/supported-languages#go) to execute and respond to a Discord slash command.

### Prerequisites
* [Golang](https://golang.org/dl/)
* [Discord Account](https://discord.com/)
* [Vercel Account](https://vercel.com/)

### Setup

#### Preliminary Note
This guide does not cover using environment variables, GitHub secrets, or Vercel secrets.
Instead, secrets are used directly within the example code to keep it simpler. Please only use
the credentials of the new application created while following this guide. When you
are finished with the guide, delete the application and create a new one and employ
environment variables and/or secrets within your code to protect your application secrets.

#### Discord Server
1. Log into your Discord account and create a new server.
2. Enable "Developer Mode" in your Discord App Settings under "Appearance".
3. In the new server you created, right click on it's name and click on "Copy ID", this is your `Guild ID`.

#### Discord Application
1. Create a new application [here](https://discord.com/developers/applications).
2. On this page you will find both the `Client ID` and `Public Key`.
    * These two secrets should not be tracked in version control

#### Discord Application Bot
1. Within the application you created above, click on "Bot" in the left pane.
2. Click on "Add Bot".
3. On this page you will find the bot's `Token`.
    * This secret should not be tracked in version control

#### Grant Application Access
1. Within the application you created above, click on "OAuth2" in the left pane.
2. Under "Scopes" check "applications.commands"
3. Copy the URL in the textbox and go to that URL in a new browser window/tab
4. In the "Add to Server" dropdown select the server you created above.

#### Add Secrets
The following secrets you have collected above should not be tracked in version control:
    * `Client ID`
    * `Public Key`
    * `Token`
1. Open [examples/vercel/api/index.go](./examples/vercel/api/index.go)
2. Fill in the secrets from above by making the following replacements
    * `CLIENT_ID` -> Your application `Client ID`
    * `PUBLIC_KEY` -> Your application `Public Key`
    * `TOKEN` -> Your application `Token`
    * `GUILD_ID` -> Your server `Guild ID`

#### Deploy to Vercel
1. Ensure you are in the examples/vercel directory [here](./examples/vercel)
2. Deploy to vercel
    ```sh
    vercel
    #> Vercel CLI 21.3.3
    #> ? Set up and deploy “~/disgoslash/examples/vercel”? [Y/n]
    y
    #> ? Which scope do you want to deploy to?
    Your Name
    #> ? Link to existing project? [y/N]
    n
    #> ? What’s your project’s name?
    hellotest
    #> ? In which directory is your code located?
    ./
    #> No framework detected. Default Project Settings:
    #> - Build Command: `npm run vercel-build` or `npm run build`
    #> - Output Directory: `public` if it exists, or `.`
    #> - Development Command: None
    #> ? Want to override the settings? [y/N]
    n
    #> 🔗  Linked to yourusername/hellotest (created .vercel and added it to .gitignore)
    #> 🔍  Inspect: https://vercel.com/yourusername/hellotest/adshjfkdashjfdsal [1s]
    #> ✅  Production: https://hellotest-tau.vercel.app [copied to clipboard] [45s]
    #> 📝  Deployed to production. Run `vercel --prod` to overwrite later (https://vercel.link/2F).
    #> 💡  To change the domain or build command, go to https://vercel.com/yourusername/hellotest/settings
    ```
3. Copy the "Production" URL from the output of running `vercel`.
4. On your Discord application's General Information page set the "Interactions Endpoint URL" to the URL you copied above, adding `/api` to the end.
    ```
    https://hellotest-tau.vercel.app/api
    ```
5. After saving those changes you should see a green bar with a message like "All your edits have been carefully recorded."

#### Register Slash Command
1. Change into the [`sync`](./sync) directory
    ```sh
    cd sync
    ```
2. Run sync
    ```sh
    go run sync.go
    #> Collecting outdated commands...
    #>     Guild: GLOBAL
    #>             success
    #>     Guild: 000000000000000000
    #>             success
    #> Unregistering outdated commands...
    #>     Guild: GLOBAL, Command: hello
    #>             success
    #> Registering new commands...
    #>     Guild: 000000000000000000, Command: hello
    #>             success
    #>     Guild: GLOBAL, Command: hello
    #>             success
    ```

#### Use Slash Command
In your discord server type `/hello`, press tab, enter a name, then hit enter. The command should run and the bot should respond.
```sh
/hello name: Bob
#> Hello Bob!
```

## Outstanding Features
- [ ] Stable release version
- [ ] Internal client does not have an Edit method, it only uses Discord's create and delete endpoints right now
- [ ] Support application command permissions
- [ ] Some models in [./discord](./discord) are not translated to Go structs yet
- [ ] Exporter to export ApplicationCommands to JSON files
- [ ] CLI tool to list, delete, create, update application commands
