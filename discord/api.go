package discord

import "time"

// MaxResponseTime for a slash command before discord terminates the request
// Discord will error out the command if it takes more than 3 seconds.
const MaxResponseTime = 3 * time.Second

// BaseURL is the base Discord API url used by this package
const BaseURL string = "https://discord.com/api"

// APIVersion is the Discord API version used by this package
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
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
}
