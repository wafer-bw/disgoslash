package models

import "time"

// https://discord.com/developers/docs/resources/channel

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
