package disgoslash

import "time"

// APIErrorResponse - Discord API error response object
type APIErrorResponse struct {
	Message    string  `json:"message"`
	Code       int     `json:"code"`
	RetryAfter float32 `json:"retry_after"`
	Global     bool    `json:"global"`
}

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
type PremiumType int

// PremiumType Enum
const (
	PremiumTypeNone         PremiumType = 0
	PremiumTypeNitroClassic PremiumType = 1
	PremiumTypeNitro        PremiumType = 2
)

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
	Emojis                      []interface{}              `json:"emojis"`   // todo struct https://discord.com/developers/docs/resources/emoji#emoji-object
	Features                    []interface{}              `json:"features"` // todo enum https://discord.com/developers/docs/resources/guild#guild-object-guild-features
	MFALevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               string                     `json:"application_id"`
	SystemChannelID             string                     `json:"system_channel_id"`
	SystemChannelFlags          int                        `json:"system_channel_flags"`
	RulesChannelID              string                     `json:"rules_channel_id"`
	JoinedAt                    time.Time                  `json:"joined_at"`
	Large                       bool                       `json:"large"`
	Unavailable                 bool                       `json:"unavailable"`
	MemberCount                 int                        `json:"member_count"`
	VoiceStates                 []*interface{}             `json:"voice_states"` // todo struct https://discord.com/developers/docs/resources/voice#voice-state-object
	Members                     []*GuildMember             `json:"members"`      //
	Channels                    []*interface{}             `json:"channels"`     // todo struct https://discord.com/developers/docs/resources/channel#channel-object
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
	WelcomeScreen               interface{}                `json:"welcome_screen"` // todo struct https://discord.com/developers/docs/resources/guild#welcome-screen-object
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
type MFALevel int

// MFALevel Enum
const (
	MFALevelNone     MFALevel = 0
	MFALevelElevated MFALevel = 1
)

// VerificationLevel - The type of Verification Level
type VerificationLevel int

// VerificationLevel Enum
const (
	VerificationLevelNone     VerificationLevel = 0
	VerificationLevelLow      VerificationLevel = 1
	VerificationLevelMedium   VerificationLevel = 2
	VerificationLevelHigh     VerificationLevel = 3
	VerificationLevelVeryHigh VerificationLevel = 4
)

// NotificationLevel - The type of Notification Level
type NotificationLevel int

// NotificationLevel Enum
const (
	NotificationLevelAllMessages  NotificationLevel = 0
	NotificationLevelOnlyMentions NotificationLevel = 1
)

// ExplicitContentFilterLevel - The type of Notification Level
type ExplicitContentFilterLevel int

// ExplicitContentFilterLevel Enum
const (
	ExplicitContentFilterLevelDisabled            ExplicitContentFilterLevel = 0
	ExplicitContentFilterLevelMembersWithoutRoles ExplicitContentFilterLevel = 1
	ExplicitContentFilterLevelAllMembers          ExplicitContentFilterLevel = 2
)

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

// InteractionRequest - The base request model sent when a user invokes a command
type InteractionRequest struct {
	ID        string                             `json:"id"`
	Type      InteractionType                    `json:"type"`
	Data      *ApplicationCommandInteractionData `json:"data"`
	GuildID   string                             `json:"guild_id"`
	ChannelID string                             `json:"channel_id"`
	Member    *GuildMember                       `json:"member"`
	Token     string                             `json:"token"`
	Version   int                                `json:"version"`
}

// InteractionType - The type of the interaction
type InteractionType int

// InteractionType Enum
const (
	InteractionTypePing               InteractionType = 1
	InteractionTypeApplicationCommand InteractionType = 2
)

// InteractionResponse - The base model of a response to an interaction request
type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data"`
}

// InteractionResponseType - The type of the response
type InteractionResponseType int

// InteractionResponseType Enum
const (
	InteractionResponseTypePong                     InteractionResponseType = 1
	InteractionResponseTypeAcknowledge              InteractionResponseType = 2
	InteractionResponseTypeChannelMessage           InteractionResponseType = 3
	InteractionResponseTypeChannelMessageWithSource InteractionResponseType = 4
	InteractionResponseTypeAcknowledgeWithSource    InteractionResponseType = 5
)

// InteractionApplicationCommandCallbackData - Optional response message payload
type InteractionApplicationCommandCallbackData struct {
	TTS             bool             `json:"tts"`
	Content         string           `json:"content"`
	Embeds          []*Embed         `json:"embeds"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions"`
}

// ApplicationCommandInteractionData - The command data payload
type ApplicationCommandInteractionData struct {
	ID      string                                     `json:"id"`
	Name    string                                     `json:"name"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

// ApplicationCommandInteractionDataOption - The params + values from the user
type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Value   string                                     `json:"value"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

// AllowedMentionType - The type of allowed mention
type AllowedMentionType string

// AllowedMentionType Enum
const (
	AllowedMentionTypeRoleMentions     AllowedMentionType = "roles"
	AllowedMentionTypeUserMentions     AllowedMentionType = "users"
	AllowedMentionTypeEveryoneMentions AllowedMentionType = "everyone"
)

// ApplicationCommand - The base commmand model that belongs to an application
type ApplicationCommand struct {
	ID            string                      `json:"id"`
	ApplicationID string                      `json:"application_id"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Options       []*ApplicationCommandOption `json:"options"`
}

// ApplicationCommandOption - The parameters for the command
type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType      `json:"type"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Required    bool                              `json:"required"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices"`
	Options     []*ApplicationCommandOption       `json:"options"`
}

// ApplicationCommandOptionChoice - User choice for `string` and/or `int` type options
type ApplicationCommandOptionChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ApplicationCommandOptionType - Types of command options
type ApplicationCommandOptionType int

// ApplicationCommandOptionType Enum
const (
	ApplicationCommandOptionTypeSubCommand      ApplicationCommandOptionType = 1
	ApplicationCommandOptionTypeSubCommandGroup ApplicationCommandOptionType = 2
	ApplicationCommandOptionTypeString          ApplicationCommandOptionType = 3
	ApplicationCommandOptionTypeInteger         ApplicationCommandOptionType = 4
	ApplicationCommandOptionTypeBoolean         ApplicationCommandOptionType = 5
	ApplicationCommandOptionTypeUser            ApplicationCommandOptionType = 6
	ApplicationCommandOptionTypeChannel         ApplicationCommandOptionType = 7
	ApplicationCommandOptionTypeRole            ApplicationCommandOptionType = 8
)

// Embed - an embed object
type Embed struct {
	Title       string     `json:"title"`
	Type        EmbedType  `json:"type"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Timestamp   time.Time  `json:"timestamp"`
	Color       int        `json:"color"`
	Footer      *Footer    `json:"footer"`
	Image       *Image     `json:"image"`
	Thumbnail   *Thumbnail `json:"thumbnail"`
	Video       *Video     `json:"video"`
	Provider    *Provider  `json:"provider"`
	Author      *Author    `json:"author"`
	Fields      []*Field   `json:"fields"`
}

// EmbedType - The type of the embed
type EmbedType string

// EmbedType Enum
const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGIFV    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

// Footer - Embed footer object
type Footer struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

// Image - Embed image object
type Image struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// Provider - Embed provider object
type Provider struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Thumbnail - Embed thumbnail object
type Thumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// Video - Embed video object
type Video struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// Author - Embed author object
type Author struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

// Field - Embed field object
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"Value"`
	Inline bool   `json:"inline"`
}

// AllowedMentions - Used to control mentions
type AllowedMentions struct {
	Parse       []AllowedMentionType `json:"parse"`
	Roles       []string             `json:"roles"`
	Users       []string             `json:"users"`
	RepliedUser bool                 `json:"replied_user"`
}
