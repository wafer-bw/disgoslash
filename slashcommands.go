package disgoslash

import (
	"strings"

	"github.com/wafer-bw/disgoslash/discord"
)

// SlashCommand properties
type SlashCommand struct {
	Do         Do
	Name       string
	GuildIDs   []string
	AppCommand *discord.ApplicationCommand
}

// Do work
type Do func(request *discord.InteractionRequest) (*discord.InteractionResponse, error)

// SlashCommandMap of slash commands using the lowercase slash command name as keys
type SlashCommandMap map[string]SlashCommand

// NewSlashCommand creates a new `SlashCommand` object
func NewSlashCommand(name string, appCommand *discord.ApplicationCommand, do Do, global bool, guildIDs []string) SlashCommand {
	if guildIDs == nil {
		guildIDs = []string{}
	}
	if global {
		guildIDs = append(guildIDs, "")
	}
	return SlashCommand{
		Name:       strings.ToLower(name),
		AppCommand: appCommand,
		Do:         do,
		GuildIDs:   guildIDs,
	}
}

// NewSlashCommandMap returns a new `SlashCommandMap`
func NewSlashCommandMap(slashCommands ...SlashCommand) SlashCommandMap {
	scm := SlashCommandMap{}
	scm.add(slashCommands...)
	return scm
}

func (scm SlashCommandMap) add(slashCommandsSlice ...SlashCommand) {
	for _, command := range slashCommandsSlice {
		scm[strings.ToLower(command.Name)] = command
	}
}
