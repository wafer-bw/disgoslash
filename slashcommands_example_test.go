package disgoslash_test

import (
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
)

var slashCommand disgoslash.SlashCommand
var anotherSlashCommand disgoslash.SlashCommand
var slashCommandMap disgoslash.SlashCommandMap

func action(request *discord.InteractionRequest) *discord.InteractionResponse {
	name, _ := request.Data.Options[0].StringValue()
	return &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{
			Content: "Hello " + name + "!",
		},
	}
}

func ExampleNewSlashCommand() {
	isGlobal := true
	guildIDs := []string{"GUILD_ID", "ANOTHER_GUILD_ID"}
	applicationCommand := &discord.ApplicationCommand{
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
		DefaultPermission: true,
	}

	slashCommand = disgoslash.NewSlashCommand(applicationCommand, action, isGlobal, guildIDs)
}

func ExampleNewSlashCommandMap() {
	slashCommands := []disgoslash.SlashCommand{slashCommand, anotherSlashCommand}
	slashCommandMap = disgoslash.NewSlashCommandMap(slashCommands...)
}
