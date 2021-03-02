package disgoslash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/discord"
	"github.com/wafer-bw/disgoslash/errs"
)

func TestNewClient(t *testing.T) {
	c := newClient(&discord.Credentials{PublicKey: "a", ClientID: "b", Token: "c"})
	require.IsType(t, &client{}, c)
}

func TestList(t *testing.T) {
	t.Run("success/global", func(t *testing.T) {
		commandName := `hello`
		mockResponse := fmt.Sprintf(`[{"id": "12345","application_id": "12345","name": "%s","description": "says hello"}]`, commandName)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte(mockResponse))
			require.NoError(t, err)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		commands, err := client.list("")
		require.NoError(t, err)
		require.Equal(t, commands[0].Name, commandName)
	})
	t.Run("success/guild", func(t *testing.T) {
		commandName := `hello`
		mockResponse := fmt.Sprintf(`[{"id": "12345","application_id": "12345","name": "%s","description": "says hello"}]`, commandName)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte(mockResponse))
			require.NoError(t, err)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		commands, err := client.list(guildID)
		require.NoError(t, err)
		require.Equal(t, commands[0].Name, commandName)
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		_, err := client.list(guildID)
		require.Error(t, err)
	})
	t.Run("failure/internal server error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		_, err := client.list(guildID)
		require.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success/global", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.delete("", "12345")
		require.NoError(t, err)
	})
	t.Run("success/guild", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.delete("12345", "12345")
		require.NoError(t, err)
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.delete("12345", "12345")
		require.Error(t, err)
	})
	t.Run("failure/internal server error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.delete("12345", "12345")
		require.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("success/global", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.create("", &discord.ApplicationCommand{})
		require.NoError(t, err)
	})
	t.Run("success/guild", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.create("12345", &discord.ApplicationCommand{})
		require.NoError(t, err)
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.create("12345", &discord.ApplicationCommand{})
		require.Error(t, err)
	})
	t.Run("failure/already exists", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.create("12345", &discord.ApplicationCommand{})
		require.Error(t, err)
		require.Equal(t, err, errs.ErrAlreadyExists)
	})
	t.Run("failure/internal server error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		err := client.create("12345", &discord.ApplicationCommand{})
		require.Error(t, err)
	})
}

func TestRequest(t *testing.T) {
	t.Run("success/rate limited", func(t *testing.T) {
		retryAfter := 0.25
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			mockResponse := &discord.APIErrorResponse{Message: "You are being rate limited.", RetryAfter: retryAfter}
			if retryAfter != 0 {
				retryAfter = 0
				w.WriteHeader(http.StatusTooManyRequests)
				mockResponseData, err := json.Marshal(mockResponse)
				require.NoError(t, err)
				_, err = w.Write([]byte(mockResponseData))
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		status, data, err := client.request(http.MethodGet, mockServer.URL, nil)
		require.NoError(t, err)
		require.Equal(t, "", string(data))
		require.Equal(t, http.StatusOK, status)
	})
	t.Run("failure/max retries hit", func(t *testing.T) {
		retryAfter := 0.01
		attempt := 0
		maxTestAttempts := maxAttempts + 1

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			mockResponse := &discord.APIErrorResponse{Message: "You are being rate limited.", RetryAfter: retryAfter}
			if attempt < maxTestAttempts {
				attempt++
				w.WriteHeader(http.StatusTooManyRequests)
				mockResponseData, err := json.Marshal(mockResponse)
				require.NoError(t, err)
				_, err = w.Write([]byte(mockResponseData))
				require.NoError(t, err)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		_, _, err := client.request(http.MethodGet, mockServer.URL, nil)
		require.Error(t, err)
		require.Equal(t, errs.ErrMaxRetries, err)
	})
	t.Run("failure/forbidden", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer func() { mockServer.Close() }()
		client := constructClient(&discord.Credentials{}, mockServer.URL, apiVersion)

		_, _, err := client.request(http.MethodGet, mockServer.URL, nil)
		require.Error(t, err)
		require.Equal(t, errs.ErrForbidden, err)
	})
}
