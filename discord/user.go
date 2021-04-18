package discord

// https://discord.com/developers/docs/resources/user

// User - A discord user
type User struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	Avatar        string      `json:"avatar"`
	Bot           bool        `json:"bot"`
	System        bool        `json:"system"`
	MFAEnabled    bool        `json:"mfa_enabled"`
	Locale        string      `json:"locale"`
	Verified      bool        `json:"verified"`
	Email         string      `json:"email"`
	Flags         int         `json:"flags"`
	PremiumType   PremiumType `json:"premium_type"`
	PublicFlags   int         `json:"public_flags"`
}

// PremiumType - The type of premium subscription
type PremiumType uint8

// PremiumType Enum
const (
	PremiumTypeNone PremiumType = iota
	PremiumTypeNitroClassic
	PremiumTypeNitro
)
