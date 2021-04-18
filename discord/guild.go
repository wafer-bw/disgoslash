package discord

import (
	"encoding/json"
	"time"
)

// Guild - the base model of a guild / server
type Guild struct {
	ID                          string                     `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        string                     `json:"icon"`
	IconHash                    string                     `json:"icon_hash"`
	Splash                      string                     `json:"splash"`
	DiscoverySplash             string                     `json:"discovery_splash"`
	Owner                       bool                       `json:"owner"`
	OwnerID                     string                     `json:"owner_id"`
	Permissions                 string                     `json:"permissions"`
	Region                      string                     `json:"region"`
	AFKChannelID                string                     `json:"afk_channel_id"`
	AFKTimeout                  int                        `json:"afk_timeout"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             string                     `json:"widget_channel_id"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	DefaultMessageNotifications NotificationLevel          `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Roles                       []Role                     `json:"roles"`
	Emojis                      json.RawMessage            `json:"emojis"`   // TODO struct https://discord.com/developers/docs/resources/emoji#emoji-object
	Features                    json.RawMessage            `json:"features"` // TODO enum https://discord.com/developers/docs/resources/guild#guild-object-guild-features
	MFALevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               string                     `json:"application_id"`
	SystemChannelID             string                     `json:"system_channel_id"`
	SystemChannelFlags          int                        `json:"system_channel_flags"`
	RulesChannelID              string                     `json:"rules_channel_id"`
	JoinedAt                    time.Time                  `json:"joined_at"`
	Large                       bool                       `json:"large"`
	Unavailable                 bool                       `json:"unavailable"`
	MemberCount                 int                        `json:"member_count"`
	VoiceStates                 json.RawMessage            `json:"voice_states"` // TODO struct https://discord.com/developers/docs/resources/voice#voice-state-object
	Members                     []*GuildMember             `json:"members"`      //
	Channels                    json.RawMessage            `json:"channels"`     // TODO struct https://discord.com/developers/docs/resources/channel#channel-object
	Presences                   []*Presence                `json:"presences"`
	MaxPresences                int                        `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	VanityURLCode               string                     `json:"vanity_url_code"`
	Description                 string                     `json:"description"`
	Banner                      string                     `json:"banner"`
	PremiumTier                 int                        `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      string                     `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	ApproximateMemberCount      int                        `json:"approximate_member_count"`
	ApproximatePresenceCount    int                        `json:"approximate_presence_count"`
	WelcomeScreen               json.RawMessage            `json:"welcome_screen"` // TODO struct https://discord.com/developers/docs/resources/guild#welcome-screen-object
}

// GuildMember - The properties of a member of a guild
type GuildMember struct {
	User         *User     `json:"user"`
	Nick         string    `json:"nick"`
	Roles        []string  `json:"roles"`
	JoinedAt     time.Time `json:"joined_at"`
	PremiumSince time.Time `json:"premium_since"`
	Deaf         bool      `json:"deaf"`
	Mute         bool      `json:"mute"`
	Pending      bool      `json:"pending"`
	Permissions  string    `json:"permissions"`
}

// MFALevel - The type of MFA Level
type MFALevel uint8

// MFALevel Enum
const (
	MFALevelNone MFALevel = iota
	MFALevelElevated
)

// VerificationLevel - The type of Verification Level
type VerificationLevel uint8

// VerificationLevel Enum
const (
	VerificationLevelNone VerificationLevel = iota
	VerificationLevelLow
	VerificationLevelMedium
	VerificationLevelHigh
	VerificationLevelVeryHigh
)

// NotificationLevel - The type of Notification Level
type NotificationLevel uint8

// NotificationLevel Enum
const (
	NotificationLevelAllMessages NotificationLevel = iota
	NotificationLevelOnlyMentions
)

// ExplicitContentFilterLevel - The type of Notification Level
type ExplicitContentFilterLevel uint8

// ExplicitContentFilterLevel Enum
const (
	ExplicitContentFilterLevelDisabled ExplicitContentFilterLevel = iota
	ExplicitContentFilterLevelMembersWithoutRoles
	ExplicitContentFilterLevelAllMembers
)
