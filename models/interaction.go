package models

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
