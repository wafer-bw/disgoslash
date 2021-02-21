package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// EnvVars defines expected & required environment variables
type EnvVars struct {
	PublicKey string `envconfig:"PUBLIC_KEY" required:"true" split_words:"true"`
	ClientID  string `envconfig:"CLIENT_ID" required:"true" split_words:"true"`
	Token     string `envconfig:"TOKEN" required:"true" split_words:"true"`
}

// Config holds all config data
type Config struct {
	Credentials *Credentials
	DiscordAPI  *DiscordAPI
}

// DiscordAPI config data
type DiscordAPI struct {
	BaseURL     string
	APIVersion  string
	ContentType string
}

// Credentials config data
type Credentials struct {
	PublicKey string
	ClientID  string
	Token     string
}

// DiscordAPIConf object
var DiscordAPIConf = &DiscordAPI{
	BaseURL:     "https://discord.com/api",
	APIVersion:  "v8",
	ContentType: "application/json",
}

// New returns a new `Config` struct; panics if unable
func New(creds *Credentials) *Config {
	if creds != nil {
		return &Config{Credentials: creds, DiscordAPI: DiscordAPIConf}
	}

	env := getEnvVars()
	ensureNoBlankEnvVars(env)
	return &Config{
		Credentials: &Credentials{
			PublicKey: env.PublicKey,
			ClientID:  env.ClientID,
			Token:     env.Token,
		},
		DiscordAPI: DiscordAPIConf,
	}
}

func getEnvVars() EnvVars {
	var env EnvVars
	err := envconfig.Process("", &env)
	if err != nil {
		panic(err)
	}
	return env
}

func ensureNoBlankEnvVars(env EnvVars) {
	blanks := findBlankEnvVars(env)
	if len(blanks) > 0 {
		panic(fmt.Errorf("the following environment variables are blank: %s", strings.Join(blanks, ", ")))
	}
}

func findBlankEnvVars(env EnvVars) []string {
	var blanks []string
	valueOfStruct := reflect.ValueOf(env)
	typeOfStruct := valueOfStruct.Type()
	for i := 0; i < valueOfStruct.NumField(); i++ {
		if valueOfStruct.Field(i).Interface() == "" {
			blanks = append(blanks, typeOfStruct.Field(i).Name)
		}
	}
	return blanks
}
