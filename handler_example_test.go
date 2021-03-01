package disgoslash_test

import (
	"net/http"

	"github.com/wafer-bw/disgoslash"
	discord "github.com/wafer-bw/disgoslash/discord"
)

func ExampleHandler_Handle() {
	creds := &discord.Credentials{
		PublicKey: "YOUR_DISCORD_APPLICATION_PUBLIC_KEY",
		ClientID:  "YOUR_DISCORD_APPLICATION_CLIENT_ID",
		Token:     "YOUR_DISCORD_BOT_TOKEN",
	}

	handler := &disgoslash.Handler{
		SlashCommandMap: disgoslash.SlashCommandMap{},
		Creds:           creds,
	}
	http.HandleFunc("/", handler.Handle)
}
