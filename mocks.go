package disgoslash

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
