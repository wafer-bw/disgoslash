package disgoslash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wafer-bw/disgoslash/models"
)

// client implements a `Client` interface's properties
type client struct {
	conf    *Config
	apiURL  string
	headers map[string]string
}

// Client interface methods
type Client interface {
	ListApplicationCommands(guildID string) ([]*models.ApplicationCommand, error)
	CreateApplicationCommand(guildID string, command *models.ApplicationCommand) error
	DeleteApplicationCommand(guildID string, commandID string) error
}

// NewClient creates a new `Client` instance
func NewClient(creds *Credentials) Client {
	conf := NewConfig(creds)
	return constructClient(conf)
}

func constructClient(conf *Config) Client {
	return &client{
		conf:   conf,
		apiURL: fmt.Sprintf("%s/%s/applications/%s", conf.discordAPI.baseURL, conf.discordAPI.apiVersion, conf.Credentials.ClientID),
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bot %s", conf.Credentials.Token),
			"Content-Type":  conf.discordAPI.contentType,
		},
	}
}

func (client *client) ListApplicationCommands(guildID string) ([]*models.ApplicationCommand, error) {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", client.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", client.apiURL, guildID)
	}
	return client.listApplicationCommands(url)
}

func (client *client) CreateApplicationCommand(guildID string, command *models.ApplicationCommand) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", client.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", client.apiURL, guildID)
	}
	return client.createApplicationCommand(url, command)
}

func (client *client) DeleteApplicationCommand(guildID string, commandID string) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands/%s", client.apiURL, commandID)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands/%s", client.apiURL, guildID, commandID)
	}
	return client.deleteApplicationCommands(url)
}

func (client *client) listApplicationCommands(url string) ([]*models.ApplicationCommand, error) {
	status, data, err := request(http.MethodGet, url, client.headers, nil)
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", status, string(data))
	}
	commands := &[]*models.ApplicationCommand{}
	if err := unmarshal(data, commands); err != nil {
		return nil, err
	}
	return *commands, nil
}

func (client *client) createApplicationCommand(url string, command *models.ApplicationCommand) error {
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

func (client *client) deleteApplicationCommands(url string) error {
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
	responseErr := &models.APIErrorResponse{}
	if err := unmarshal(data, responseErr); err != nil {
		return 0, err
	}
	return time.Duration(responseErr.RetryAfter) * time.Second, nil
}
