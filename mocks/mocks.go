package mocks

import (
	"github.com/wafer-bw/disgoslash/config"
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
func GetConf() *config.Config {
	return &config.Config{
		Credentials: &config.Credentials{
			PublicKey: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			ClientID:  "abc123",
			Token:     "abc123",
		},
		DiscordAPI: &config.DiscordAPI{
			BaseURL:     "https://discord.com/api",
			APIVersion:  "v8",
			ContentType: "application/json",
		},
	}
}

// Conf mocks the `config.Config` object
var Conf = GetConf()
