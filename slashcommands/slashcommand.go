package slashcommands

import (
	"strings"

	"github.com/wafer-bw/disgoslash/models"
)

// SlashCommand properties
type SlashCommand struct {
	Do         Do
	Name       string
	GuildIDs   []string
	AppCommand *models.ApplicationCommand
}

// Do work
type Do func(request *models.InteractionRequest) (*models.InteractionResponse, error)

// New `SlashCommand`
func New(name string, appCommand *models.ApplicationCommand, do Do, global bool, guildIDs []string) SlashCommand {
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
