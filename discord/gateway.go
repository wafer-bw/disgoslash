package discord

// Presence - A user's current state on a guild
type Presence struct {
	User         User           `json:"user"`
	GuildID      string         `json:"guild_id"`
	Status       PresenceStatus `json:"status"`
	Activities   []interface{}  `json:"activities"` // todo - https://discord.com/developers/docs/topics/gateway#activity-object
	ClientStatus ClientStatus   `json:"client_status"`
}

// PresenceStatus - The type of PresenceStatus
type PresenceStatus string

// PresenceStatus Enum
const (
	Idle         PresenceStatus = "idle"
	DoNotDisturb PresenceStatus = "dnd"
	Online       PresenceStatus = "online"
	Offline      PresenceStatus = "offline"
)

// ClientStatus - Active session statuses for a user per platform
type ClientStatus struct {
	Web     PresenceStatus
	Mobile  PresenceStatus
	Desktop PresenceStatus
}
