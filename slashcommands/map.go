package slashcommands

import (
	"strings"
)

// Map of slash commands using the lowercase slash command name as keys
type Map map[string]SlashCommand

// NewMap returns a new `Map`
func NewMap(slashCommands ...SlashCommand) Map {
	scm := Map{}
	scm.add(slashCommands...)
	return scm
}

func (scm Map) add(slashCommandsSlice ...SlashCommand) {
	for _, command := range slashCommandsSlice {
		scm[strings.ToLower(command.Name)] = command
	}
}
