package discord

const BaseURL string = "https://discord.com/api"
const APIVersion string = "v8"

// Credentials required from your Discord application & bot
type Credentials struct {
	PublicKey string
	ClientID  string
	Token     string
}

// APIErrorResponse - Discord API error response object
type APIErrorResponse struct {
	Message    string  `json:"message"`
	Code       int     `json:"code"`
	RetryAfter float32 `json:"retry_after"`
	Global     bool    `json:"global"`
}
