package disgoslash

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type envVars struct {
	PublicKey string `envconfig:"PUBLIC_KEY" required:"true" split_words:"true"`
	ClientID  string `envconfig:"CLIENT_ID" required:"true" split_words:"true"`
	Token     string `envconfig:"TOKEN" required:"true" split_words:"true"`
}

// Config holds all config data
type Config struct {
	Credentials *Credentials
	discordAPI  *discordAPI
}

type discordAPI struct {
	baseURL     string
	apiVersion  string
	contentType string
}

// Credentials config data
type Credentials struct {
	PublicKey string
	ClientID  string
	Token     string
}

var discordAPIConf = &discordAPI{
	baseURL:     "https://discord.com/api",
	apiVersion:  "v8",
	contentType: "application/json",
}

// NewConfig creates a new `Config` object; panics if unable.
// If `creds` is `nil` it will attempt to load environment variables
func NewConfig(creds *Credentials) *Config {
	if creds != nil {
		return &Config{Credentials: creds, discordAPI: discordAPIConf}
	}

	env := getEnvVars()
	ensureNoBlankEnvVars(env)
	return &Config{
		Credentials: &Credentials{
			PublicKey: env.PublicKey,
			ClientID:  env.ClientID,
			Token:     env.Token,
		},
		discordAPI: discordAPIConf,
	}
}

func getEnvVars() envVars {
	var env envVars
	err := envconfig.Process("", &env)
	if err != nil {
		panic(err)
	}
	return env
}

func ensureNoBlankEnvVars(env envVars) {
	blanks := findBlankEnvVars(env)
	if len(blanks) > 0 {
		panic(fmt.Errorf("the following environment variables are blank: %s", strings.Join(blanks, ", ")))
	}
}

func findBlankEnvVars(env envVars) []string {
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
