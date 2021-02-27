package disgoslash

import (
	"strings"

	"github.com/wafer-bw/disgoslash/discord"
)

// SlashCommand object properties
type SlashCommand struct {
	// The work to do when a slash command is invoked by a user
	Action Action

	// The name of the slash command which the user will invoke
	// Ex: "/slashcommandname"
	Name string

	// The IDs of the guilds (servers) the application command
	// should be registered to
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
type Action func(request *discord.InteractionRequest) (*discord.InteractionResponse, error)

// SlashCommandMap using each slash command's
// application command lowercase name as a key
type SlashCommandMap map[string]SlashCommand

// NewSlashCommand creates a new SlashCommand which must
// be added to a SlashCommandMap via NewSlashCommandMap()
// to be used by the Syncer and/or Handler.
//
// In order for an application command to be registered
// to a guild (server), the bot will need to be granted
// the "appliations.commands" scope for that server.
//
// A global command will be registered to all servers
// the bot has been granted access to.
//
// The guild (server) IDs list indicates which guilds to
// register the application command to.
//
// https://discord.com/developers/docs/interactions/slash-commands#authorizing-your-application
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

// NewSlashCommandMap creates a SlashCommandMap to be passed
// to either...
//
// A Syncer to register your application commands.
//
// A Handler which receives interaction
// requests coming from your Discord server(s) when a
// user invokes a slash command and executes the slash command's
// Action function.
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
