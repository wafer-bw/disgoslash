package disgoslash

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/models"
)

var guildID = "1234567890"

// todo - TestNewClient()

func TestListApplicationCommands(t *testing.T) {
	t.Run("success/global", func(t *testing.T) {
		commandName := `hello`
		mockResponse := fmt.Sprintf(`[{"id": "12345","application_id": "12345","name": "%s","description": "says hello"}]`, commandName)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// nolint
			w.Write([]byte(mockResponse))
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		commands, err := client.ListApplicationCommands("")
		require.NoError(t, err)
		require.Equal(t, commands[0].Name, commandName)
	})
	t.Run("success/guild", func(t *testing.T) {
		commandName := `hello`
		mockResponse := fmt.Sprintf(`[{"id": "12345","application_id": "12345","name": "%s","description": "says hello"}]`, commandName)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// nolint
			w.Write([]byte(mockResponse))
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		commands, err := client.ListApplicationCommands(guildID)
		require.NoError(t, err)
		require.Equal(t, commands[0].Name, commandName)
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockResponse := `{"message": "401: Unauthorized", "code": 0}`
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			// nolint
			w.Write([]byte(mockResponse))
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		_, err := client.ListApplicationCommands(guildID)
		require.Error(t, err)
	})
}

func TestDeleteApplicationCommand(t *testing.T) {
	t.Run("success/global", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.DeleteApplicationCommand("", "12345")
		require.NoError(t, err)
	})
	t.Run("success/guild", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.DeleteApplicationCommand("12345", "12345")
		require.NoError(t, err)
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.DeleteApplicationCommand("12345", "12345")
		require.Error(t, err)
	})
}

func TestCreateApplicationCommand(t *testing.T) {
	t.Run("success/global", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.CreateApplicationCommand("", &models.ApplicationCommand{})
		require.NoError(t, err)
	})
	t.Run("success/guild", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.CreateApplicationCommand("12345", &models.ApplicationCommand{})
		require.NoError(t, err)
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.CreateApplicationCommand("12345", &models.ApplicationCommand{})
		require.Error(t, err)
	})
	t.Run("failure/already exists", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer func() { mockServer.Close() }()
		Conf.discordAPI.baseURL = mockServer.URL
		client := constructClient(Conf)

		err := client.CreateApplicationCommand("12345", &models.ApplicationCommand{})
		require.Error(t, err)
		require.Equal(t, err, ErrAlreadyExists)
	})
}
