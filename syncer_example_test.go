package disgoslash_test

import (
	"github.com/wafer-bw/disgoslash"
	discord "github.com/wafer-bw/disgoslash/discord"
)

// command object that will be registered to your Discord server
// https://discord.com/developers/docs/interactions/slash-commands#applicationcommand
var command = &discord.ApplicationCommand{
	Name:        "hello",
	Description: "Says hello to the user who uses it.",
	Options: []*discord.ApplicationCommandOption{
		{
			Type:        discord.ApplicationCommandOptionTypeString,
			Name:        "Name",
			Description: "Enter your name",
			Required:    true,
		},
	},
}

// hello is the function that runs when a user uses the hello slash command
func hello(request *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	username := request.Data.Options[0].Value
	return &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{
			Content: "Hello " + username + "!",
		},
	}, nil
}

func ExampleSyncer_Sync() {
	creds := &discord.Credentials{
		PublicKey: "YOUR_DISCORD_APPLICATION_PUBLIC_KEY",
		ClientID:  "YOUR_DISCORD_APPLICATION_CLIENT_ID",
		Token:     "YOUR_DISCORD_BOT_TOKEN",
	}

	globalCommand := true
	guildIDs := []string{"YOUR_GUILD_(SERVER)_ID", "ANOTHER_GUILD_ID"}

	command := disgoslash.NewSlashCommand(command, hello, globalCommand, guildIDs)
	commands := disgoslash.NewSlashCommandMap(command)

	syncer := &disgoslash.Syncer{
		SlashCommandMap: commands,
		GuildIDs:        guildIDs,
		Creds:           creds,
	}

	syncer.Sync()
}
