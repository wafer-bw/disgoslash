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
	InteractionResponseTypePong InteractionResponseType = 1
	// InteractionResponseTypeAcknowledge has been deprecated by Discord
	InteractionResponseTypeAcknowledge InteractionResponseType = 2
	// InteractionResponseTypeChannelMessage has been deprecated by Discord
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
	Value   json.RawMessage                            `json:"value"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

// GetString returns the unmarshalled Value of the option as a string pointer,
// if the Value was null in the source JSON it will the pointer will be nil.
//
// This exists for clarity and should be used when the interaction's
// corresponding ApplicationCommandOptionType is ApplicationCommandOptionTypeString
// in which case the Discord API will send a string Value in the interaction request
func (option ApplicationCommandInteractionDataOption) GetString() (*string, error) {
	val := new(string)
	if err := json.Unmarshal(option.Value, val); err != nil {
		return nil, err
	}
	return val, nil
}

// GetInt returns the unmarshalled Value of the option as an int pointer,
// if the Value was null in the source JSON it will the pointer will be nil.
//
// This exists for clarity and should be used when the interaction's
// corresponding ApplicationCommandOptionType is ApplicationCommandOptionTypeInteger
// in which case the Discord API will send an int Value in the interaction request
func (option ApplicationCommandInteractionDataOption) GetInt() (*int, error) {
	val := new(int)
	if err := json.Unmarshal(option.Value, val); err != nil {
		return nil, err
	}
	return val, nil
}

// GetBool returns the unmarshalled Value of the option as a bool pointer,
// if the Value was null in the source JSON it will the pointer will be nil.
//
// This exists for clarity and should be used when the interaction's
// correspoinding ApplicationCommandOptionType is ApplicationCommandOptionTypeBoolean
// in which case the Discord API will send a bool Value in the interaction request
func (option ApplicationCommandInteractionDataOption) GetBool() (*bool, error) {
	val := new(bool)
	if err := json.Unmarshal(option.Value, val); err != nil {
		return nil, err
	}
	return val, nil
}

// GetUserID is an alias for ApplicationCommandInteractionDataOption.GetString
// returning the unmarshalled Value of the option as a string pointer.
//
// This exists for clarity and should be used when the interaction's
// correspoinding ApplicationCommandOptionType is ApplicationCommandOptionTypeUser
// in which case the Discord API will send the User ID as a string Value in the interaction request
func (option ApplicationCommandInteractionDataOption) GetUserID() (*string, error) {
	return option.GetString()
}

// GetRoleID is an alias for ApplicationCommandInteractionDataOption.GetString
// returning the unmarshalled Value of the option as a string pointer.
//
// This exists for clarity and should be used when the interaction's
// correspoinding ApplicationCommandOptionType is ApplicationCommandOptionTypeRole
// in which case the Discord API will send the Role ID as a string Value in the interaction request
func (option ApplicationCommandInteractionDataOption) GetRoleID() (*string, error) {
	return option.GetString()
}

// GetChannelID is an alias for ApplicationCommandInteractionDataOption.GetString
// returning the unmarshalled Value of the option as a string pointer.
//
// This exists for clarity and should be used when the interaction's
// correspoinding ApplicationCommandOptionType is ApplicationCommandOptionTypeChannel
// in which case the Discord API will send the Channel ID as a string Value in the interaction request
func (option ApplicationCommandInteractionDataOption) GetChannelID() (*string, error) {
	return option.GetString()
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
type ApplicationCommandPermissionType int

// ApplicationCommandPermissionType Enum
const (
	ApplicationCommandPermissionTypeSubRole ApplicationCommandPermissionType = 1
	ApplicationCommandPermissionTypeSubUser ApplicationCommandPermissionType = 2
)
