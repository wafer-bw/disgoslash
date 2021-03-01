package disgoslash

import (
	"strings"

	"github.com/wafer-bw/disgoslash/discord"
)

// SlashCommand object properties
type SlashCommand struct {
	// The work to do when a slash command is invoked by a user
	Action Action

	// 1-32 character name matching ^[\w-]{1,32}$
	// Ex: "/tableflip"
	Name string

	// The IDs of the guilds (servers) the application command
	// should be registered to. To register the command globally
	// include a blank string ("").
	GuildIDs []string

	// The application command object schema which will be
	// registered to your Discord servers
	ApplicationCommand *discord.ApplicationCommand
}

// Action is the function executed when a
// slash command interaction request is received
// by the Handler.
//
// This is where you put the work to be done when
// a user invokes the slash command.
//
// An action function does not include an error
// in its response signature because errors should
// be handled gracefully and a message informing the
// user that something went wrong should be included
// in the interaction response.
type Action func(request *discord.InteractionRequest) *discord.InteractionResponse

// SlashCommandMap using each slash command's
// application command name as a key
type SlashCommandMap map[string]SlashCommand

// NewSlashCommand creates a new SlashCommand
func NewSlashCommand(appCommand *discord.ApplicationCommand, action Action, global bool, guildIDs []string) SlashCommand {
	if guildIDs == nil {
		guildIDs = []string{}
	}
	if global {
		guildIDs = append(guildIDs, "")
	}
	return SlashCommand{
		Name:               strings.ToLower(appCommand.Name),
		ApplicationCommand: appCommand,
		Action:             action,
		GuildIDs:           guildIDs,
	}
}

// NewSlashCommandMap creates a new SlashCommandMap
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
