package disgoslash

// MockSlashCommandName mocks a command name
var MockSlashCommandName = "hello"

// MockSlashCommandResponseContent mocks a command response message
var MockSlashCommandResponseContent = "Hello World!"

// MockInteractionResponse mocks an interaciton response object
var MockInteractionResponse = &InteractionResponse{
	Type: InteractionResponseTypeChannelMessageWithSource,
	Data: &InteractionApplicationCommandCallbackData{Content: MockSlashCommandResponseContent},
}

// SlashCommandDo mocks a command `Do` function
func SlashCommandDo(request *InteractionRequest) (*InteractionResponse, error) {
	return MockInteractionResponse, nil
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
