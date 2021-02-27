package disgoslash_test

import (
	"net/http"

	"github.com/wafer-bw/disgoslash"
	discord "github.com/wafer-bw/disgoslash/discord"
)

// A Handler requires your Discord credentials and a map of slash commands
// created by disgoslash.NewSlashCommandMap() which accepts  slash commands
// created by disgoslash.NewSlashCommand().
func ExampleHandler_Handle() {
	creds := &discord.Credentials{
		PublicKey: "YOUR_DISCORD_APPLICATION_PUBLIC_KEY",
		ClientID:  "YOUR_DISCORD_APPLICATION_CLIENT_ID",
		Token:     "YOUR_DISCORD_BOT_TOKEN",
	}

	handler := &disgoslash.Handler{
		SlashCommandMap: disgoslash.SlashCommandMap{}, // Replace this with your slash command map
		Creds:           creds,
	}

	http.HandleFunc("/", handler.Handle)
}
