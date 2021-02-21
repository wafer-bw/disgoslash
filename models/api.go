package models

// APIErrorResponse - Discord API error response object
type APIErrorResponse struct {
	Message    string  `json:"message"`
	Code       int     `json:"code"`
	RetryAfter float32 `json:"retry_after"`
	Global     bool    `json:"global"`
}
