package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wafer-bw/disgoslash/config"
	"github.com/wafer-bw/disgoslash/errs"
	"github.com/wafer-bw/disgoslash/models"
)

// impl implements `Client` properties
type impl struct {
	conf    *config.Config
	apiURL  string
	headers map[string]string
}

// Client interfaces `Client` methods
type Client interface {
	ListApplicationCommands(guildID string) ([]*models.ApplicationCommand, error)
	CreateApplicationCommand(guildID string, command *models.ApplicationCommand) error
	DeleteApplicationCommand(guildID string, commandID string) error
}

// New returns a new `Client` interface
func New(creds *config.Credentials) Client {
	conf := config.New(creds)
	return construct(conf)
}

// construct a new `Client` interface
func construct(conf *config.Config) Client {
	return &impl{
		conf:   conf,
		apiURL: fmt.Sprintf("%s/%s/applications/%s", conf.DiscordAPI.BaseURL, conf.DiscordAPI.APIVersion, conf.Credentials.ClientID),
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bot %s", conf.Credentials.Token),
			"Content-Type":  conf.DiscordAPI.ContentType,
		},
	}
}

func (impl *impl) ListApplicationCommands(guildID string) ([]*models.ApplicationCommand, error) {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", impl.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", impl.apiURL, guildID)
	}
	return impl.listApplicationCommands(url)
}

func (impl *impl) CreateApplicationCommand(guildID string, command *models.ApplicationCommand) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", impl.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", impl.apiURL, guildID)
	}
	return impl.createApplicationCommand(url, command)
}

func (impl *impl) DeleteApplicationCommand(guildID string, commandID string) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands/%s", impl.apiURL, commandID)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands/%s", impl.apiURL, guildID, commandID)
	}
	return impl.deleteApplicationCommands(url)
}

func (impl *impl) listApplicationCommands(url string) ([]*models.ApplicationCommand, error) {
	status, data, err := httpRequest(http.MethodGet, url, impl.headers, nil)
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

func (impl *impl) createApplicationCommand(url string, command *models.ApplicationCommand) error {
	body, err := marshal(command)
	if err != nil {
		return err
	}
	if status, data, err := httpRequest(http.MethodPost, url, impl.headers, body); err != nil {
		return err
	} else if status == http.StatusOK {
		return errs.ErrAlreadyExists
	} else if status != http.StatusCreated {
		return fmt.Errorf("%d - %s", status, string(data))
	}
	return nil
}

func (impl *impl) deleteApplicationCommands(url string) error {
	if status, data, err := httpRequest(http.MethodDelete, url, impl.headers, nil); err != nil {
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

func httpRequest(method string, url string, headers map[string]string, body io.Reader) (int, []byte, error) {
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
			return 0, nil, errs.ErrForbidden
		case http.StatusUnauthorized:
			return 0, nil, errs.ErrUnauthorized
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
	return 0, nil, errs.ErrMaxRetries
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
