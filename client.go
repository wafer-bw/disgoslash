package disgoslash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	discord "github.com/wafer-bw/disgoslash/discord"
)

// Client implements a `ClientInterface` interface's properties
type Client struct {
	apiURL  string
	headers map[string]string
}

type conf struct {
	baseURL     string
	apiVersion  string
	contentType string
}

// ClientInterface methods
type ClientInterface interface {
	ListApplicationCommands(guildID string) ([]*discord.ApplicationCommand, error)
	CreateApplicationCommand(guildID string, command *discord.ApplicationCommand) error
	DeleteApplicationCommand(guildID string, commandID string) error
}

// NewClient creates a new `ClientInterface` instance
func NewClient(creds *discord.Credentials) ClientInterface {
	return constructClient(creds, newConf())
}

func constructClient(creds *discord.Credentials, dapi *conf) ClientInterface {
	return &Client{
		apiURL: fmt.Sprintf("%s/%s/applications/%s", dapi.baseURL, dapi.apiVersion, creds.ClientID),
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bot %s", creds.Token),
			"Content-Type":  dapi.contentType,
		},
	}
}

func newConf() *conf {
	return &conf{
		baseURL:     "https://discord.com/api",
		apiVersion:  "v8",
		contentType: "application/json",
	}
}

// ListApplicationCommands // todo
func (client *Client) ListApplicationCommands(guildID string) ([]*discord.ApplicationCommand, error) {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", client.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", client.apiURL, guildID)
	}
	return client.listApplicationCommands(url)
}

// CreateApplicationCommand // todo
func (client *Client) CreateApplicationCommand(guildID string, command *discord.ApplicationCommand) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", client.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", client.apiURL, guildID)
	}
	return client.createApplicationCommand(url, command)
}

// DeleteApplicationCommand // todo
func (client *Client) DeleteApplicationCommand(guildID string, commandID string) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands/%s", client.apiURL, commandID)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands/%s", client.apiURL, guildID, commandID)
	}
	return client.deleteApplicationCommands(url)
}

func (client *Client) listApplicationCommands(url string) ([]*discord.ApplicationCommand, error) {
	status, data, err := request(http.MethodGet, url, client.headers, nil)
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", status, string(data))
	}
	commands := &[]*discord.ApplicationCommand{}
	if err := unmarshal(data, commands); err != nil {
		return nil, err
	}
	return *commands, nil
}

func (client *Client) createApplicationCommand(url string, command *discord.ApplicationCommand) error {
	body, err := marshal(command)
	if err != nil {
		return err
	}
	if status, data, err := request(http.MethodPost, url, client.headers, body); err != nil {
		return err
	} else if status == http.StatusOK {
		return ErrAlreadyExists
	} else if status != http.StatusCreated {
		return fmt.Errorf("%d - %s", status, string(data))
	}
	return nil
}

func (client *Client) deleteApplicationCommands(url string) error {
	if status, data, err := request(http.MethodDelete, url, client.headers, nil); err != nil {
		return err
	} else if status != http.StatusNoContent {
		return fmt.Errorf("%d - %s", status, string(data))
	}
	return nil
}

func unmarshal(body []byte, v interface{}) error {
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

func marshal(v interface{}) (io.Reader, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func request(method string, url string, headers map[string]string, body io.Reader) (int, []byte, error) {
	attempts := 0
	maxAttempts := 3

	for attempts < maxAttempts {
		attempts++

		client := &http.Client{}
		request, err := http.NewRequest(method, url, body)
		if err != nil {
			return 0, nil, err
		}

		for key, val := range headers {
			request.Header.Set(key, val)
		}
		response, err := client.Do(request)
		if err != nil {
			return 0, nil, err
		}

		switch response.StatusCode {
		case http.StatusForbidden:
			return 0, nil, ErrForbidden
		case http.StatusUnauthorized:
			return 0, nil, ErrUnauthorized
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, nil, err
		}

		waitTime, err := determineRetry(response.StatusCode, data)
		if err != nil {
			return 0, nil, err
		}

		if waitTime <= 0 {
			return response.StatusCode, data, nil
		}
		time.Sleep(waitTime)
	}
	return 0, nil, ErrMaxRetries
}

func determineRetry(statusCode int, data []byte) (time.Duration, error) {
	if statusCode != http.StatusTooManyRequests {
		return 0, nil
	}
	responseErr := &discord.APIErrorResponse{}
	if err := unmarshal(data, responseErr); err != nil {
		return 0, err
	}
	return time.Duration(responseErr.RetryAfter) * time.Second, nil
}
