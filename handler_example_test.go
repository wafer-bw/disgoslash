package disgoslash_test

import (
	"net/http"

	"github.com/wafer-bw/disgoslash"
	discord "github.com/wafer-bw/disgoslash/discord"
)

// hi is the function that runs when a user uses the hello slash command
func hi(request *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	username := request.Data.Options[0].Value
	return &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{
			Content: "Hi " + username + "!",
		},
	}, nil
}

func ExampleHandler_Handle() {
	creds := &discord.Credentials{
		PublicKey: "YOUR_DISCORD_APPLICATION_PUBLIC_KEY",
		ClientID:  "YOUR_DISCORD_APPLICATION_CLIENT_ID",
		Token:     "YOUR_DISCORD_BOT_TOKEN",
	}

	globalCommand := true
	guildIDs := []string{"YOUR_GUILD_(SERVER)_ID", "ANOTHER_GUILD_ID"}

	// command object that will be registered to your Discord server
	// https://discord.com/developers/docs/interactions/slash-commands#applicationcommand
	command := &discord.ApplicationCommand{
		Name:        "hi",
		Description: "Says hi to the user who uses it.",
		Options: []*discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeString,
				Name:        "Name",
				Description: "Enter your name",
				Required:    true,
			},
		},
	}

	slashCommand := disgoslash.NewSlashCommand(command, hi, globalCommand, guildIDs)
	slashCommands := disgoslash.NewSlashCommandMap(slashCommand)

	handler := &disgoslash.Handler{
		SlashCommandMap: slashCommands,
		Creds:           creds,
	}

	http.HandleFunc("/", handler.Handle)
}
