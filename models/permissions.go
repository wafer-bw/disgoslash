package models

// Role - A set of permissions attached to a group of users
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       int       `json:"color"`
	Hoist       bool      `json:"hoist"`
	Position    int       `json:"position"`
	Permissions string    `json:"permissions"`
	Managed     bool      `json:"managed"`
	Mentionable bool      `json:"mentionable"`
	Tags        *RoleTags `json:"tags"`
}

// RoleTags - A set of tags applied to a `Role`
type RoleTags struct {
	BotID             string `json:"bot_id"`
	IntegrationID     string `json:"integration_id"`
	PremiumSubscriber bool   `json:"premium_subscriber"`
}
