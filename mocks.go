package disgoslash

import (
	"github.com/wafer-bw/disgoslash/models"
)

// SlashCommandName mocks a command name
var SlashCommandName = "hello"

// SlashCommandResponseContent mocks a command response message
var SlashCommandResponseContent = "Hello World!"

// InteractionResponse mocks an interaciton response object
var InteractionResponse = &models.InteractionResponse{
	Type: models.InteractionResponseTypeChannelMessageWithSource,
	Data: &models.InteractionApplicationCommandCallbackData{Content: SlashCommandResponseContent},
}

// SlashCommandDo mocks a command `Do` function
func SlashCommandDo(request *models.InteractionRequest) (*models.InteractionResponse, error) {
	return InteractionResponse, nil
}

// GetConf returns a new instance of a mocked config object
func GetConf() *Config {
	return &Config{
		Credentials: &Credentials{
			PublicKey: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			ClientID:  "abc123",
			Token:     "abc123",
		},
		discordAPI: &discordAPI{
			baseURL:     "https://discord.com/api",
			apiVersion:  "v8",
			contentType: "application/json",
		},
	}
}

// Conf mocks the `Config` object
var Conf = GetConf()
