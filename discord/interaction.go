package discord

import (
	"encoding/json"
)

// https://discord.com/developers/docs/interactions/slash-commands

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
type InteractionType uint8

// InteractionType Enum
const (
	InteractionTypePing = InteractionType(iota + 1)
	InteractionTypeApplicationCommand
)

// InteractionResponse - The base model of a response to an interaction request
type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data"`
}

// InteractionResponseType - The type of the response
type InteractionResponseType uint8

// InteractionResponseType Enum
const (
	InteractionResponseTypePong = InteractionResponseType(iota + 1)
	// InteractionResponseTypeAcknowledge has been deprecated by Discord
	InteractionResponseTypeAcknowledge
	// InteractionResponseTypeChannelMessage has been deprecated by Discord
	InteractionResponseTypeChannelMessage
	InteractionResponseTypeChannelMessageWithSource
	InteractionResponseTypeAcknowledgeWithSource
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
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

// StringValue casts the option value to a string
func (option ApplicationCommandInteractionDataOption) StringValue() (value string, ok bool) {
	value, ok = option.Value.(string)
	return value, ok
}

// IntValue casts the option value to a string
func (option ApplicationCommandInteractionDataOption) IntValue() (value int, ok bool) {
	tmpValue, ok := option.Value.(float64)
	if !ok {
		return 0, ok
	}
	value = int(tmpValue)
	return value, ok
}

// BoolValue casts the option value to a string
func (option ApplicationCommandInteractionDataOption) BoolValue() (value bool, ok bool) {
	value, ok = option.Value.(bool)
	return value, ok
}

// UserIDValue is an alias for StringValue that casts the option value to a user ID string
func (option ApplicationCommandInteractionDataOption) UserIDValue() (value string, ok bool) {
	return option.StringValue()
}

// RoleIDValue is an alias for StringValue that casts the option value to a user ID string
func (option ApplicationCommandInteractionDataOption) RoleIDValue() (value string, ok bool) {
	return option.StringValue()
}

// ChannelIDValue is an alias for StringValue that casts the option value to a user ID string
func (option ApplicationCommandInteractionDataOption) ChannelIDValue() (value string, ok bool) {
	return option.StringValue()
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
	ID                string                      `json:"id"`
	ApplicationID     string                      `json:"application_id"`
	Name              string                      `json:"name"`                         // 1-32 character name matching ^[\w-]{1,32}$
	Description       string                      `json:"description"`                  // 1-100 character description
	Options           []*ApplicationCommandOption `json:"options"`                      // the parameters for the command
	DefaultPermission bool                        `json:"default_permission,omitempty"` // whether the command is enabled by default when the app is added to a guild (defaults to true)
}

// ApplicationCommandOption - The parameters for the command
type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType      `json:"type"`
	Name        string                            `json:"name"`        // 1-32 character name matching ^[\w-]{1,32}$
	Description string                            `json:"description"` // 1-100 character description
	Required    bool                              `json:"required"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices"`
	Options     []*ApplicationCommandOption       `json:"options"`
}

// ApplicationCommandOptionChoice - User choice for `string` and/or `int` type options
// Value will always be unmarshalled as a string.
type ApplicationCommandOptionChoice struct {
	Name  string `json:"name"`
	Value string `json:"value,string"`
}

// ApplicationCommandOptionType - Types of command options
type ApplicationCommandOptionType uint8

// ApplicationCommandOptionType Enum
const (
	ApplicationCommandOptionTypeSubCommand = ApplicationCommandOptionType(iota + 1)
	ApplicationCommandOptionTypeSubCommandGroup
	ApplicationCommandOptionTypeString
	ApplicationCommandOptionTypeInteger
	ApplicationCommandOptionTypeBoolean
	ApplicationCommandOptionTypeUser
	ApplicationCommandOptionTypeChannel
	ApplicationCommandOptionTypeRole
)

type GuildApplicationCommandPermissions struct {
	ID            string          `json:"id"`             // the id of the command
	ApplicationID string          `json:"application_id"` // the id of the application the command belongs to
	GuildID       string          `json:"guild_id"`       // the id of the guild
	Permissions   json.RawMessage `json:"permissions"`    // the permissions for the command in the guild
}

type ApplicationCommandPermissions struct {
	ID         string      `json:"id"`         // the id of the role or user
	Type       interface{} `json:"type"`       // role or user
	Permission bool        `json:"permission"` // true to allow, false to disallow
}

// ApplicationCommandPermissionType - Types of application command permissions
type ApplicationCommandPermissionType uint8

// ApplicationCommandPermissionType Enum
const (
	ApplicationCommandPermissionTypeSubRole = ApplicationCommandPermissionType(iota + 1)
	ApplicationCommandPermissionTypeSubUser
)
