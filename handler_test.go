package disgoslash

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/discord"
	"github.com/wafer-bw/disgoslash/errs"
)

//nolint
var publicKey, privateKey, _ = ed25519.GenerateKey(nil)
var url = "http://localhost/api"

func TestHandle(t *testing.T) {
	interactionName := "interaction"
	testResponse := &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{Content: "Hello World!"},
	}
	do := func(request *discord.InteractionRequest) *discord.InteractionResponse {
		return testResponse
	}
	handler := &Handler{
		Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
		SlashCommandMap: NewSlashCommandMap(NewSlashCommand(&discord.ApplicationCommand{Name: interactionName, Description: "desc"}, do, true, []string{"11111"})),
	}
	handlerFunc := http.HandlerFunc(handler.Handle)
	t.Run("success/respond to ping", func(t *testing.T) {
		requestBody := `{"type":1}`
		body, resp, err := httpTestRequest(handlerFunc, http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})
	t.Run("success/run interaction", func(t *testing.T) {
		interaction := &discord.InteractionRequest{
			Type: discord.InteractionTypeApplicationCommand,
			Data: &discord.ApplicationCommandInteractionData{Name: interactionName},
		}
		data, err := json.Marshal(interaction)
		require.NoError(t, err)
		requestBody := string(data)

		body, resp, err := httpTestRequest(handlerFunc, http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})
	t.Run("failure/interaction took too long", func(t *testing.T) {
		longDo := func(_ *discord.InteractionRequest) *discord.InteractionResponse {
			time.Sleep(discord.MaxResponseTime + 500*time.Millisecond)
			return testResponse
		}
		longHandler := &Handler{
			Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
			SlashCommandMap: NewSlashCommandMap(NewSlashCommand(&discord.ApplicationCommand{Name: interactionName, Description: "desc"}, longDo, true, []string{"11111"})),
		}

		interaction := &discord.InteractionRequest{
			Type: discord.InteractionTypeApplicationCommand,
			Data: &discord.ApplicationCommandInteractionData{Name: interactionName},
		}
		data, err := json.Marshal(interaction)
		require.NoError(t, err)
		requestBody := string(data)

		body, resp, err := httpTestRequest(http.HandlerFunc(longHandler.Handle), http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
		require.Equal(t, errs.ErrTookTooLong, strings.TrimSuffix(string(body), "\n"))
	})
	t.Run("failure/unimplemented interaction", func(t *testing.T) {
		data, err := json.Marshal(&discord.InteractionRequest{
			Type: discord.InteractionTypeApplicationCommand,
			Data: &discord.ApplicationCommandInteractionData{Name: interactionName + "Z"},
		})
		require.NoError(t, err)
		requestBody := string(data)

		body, resp, err := httpTestRequest(handlerFunc, http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotImplemented, resp.StatusCode, string(body))
	})
	t.Run("failure/invalid interaction type", func(t *testing.T) {
		requestBody := `{"type": 3}`

		body, resp, err := httpTestRequest(handlerFunc, http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, string(body))
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		requestBody := `{"type": 1}`
		headers := getAuthHeaders(requestBody)
		headers["X-Signature-Ed25519"] = "X"

		body, resp, err := httpTestRequest(handlerFunc, http.MethodGet, url, headers, requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode, string(body))
	})
}

func TestUnmarshal(t *testing.T) {
	commandName := "interaction"
	do := func(request *discord.InteractionRequest) *discord.InteractionResponse {
		return nil
	}
	t.Run("success/unmarshal string option", func(t *testing.T) {
		optionName := "username"
		optionVal := "wafer-bw"
		applicationCommand := &discord.ApplicationCommand{
			Name:        commandName,
			Description: "desc",
			Options: []*discord.ApplicationCommandOption{
				{Name: optionName, Type: discord.ApplicationCommandOptionTypeString},
			},
		}
		handler := &Handler{
			Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
			SlashCommandMap: NewSlashCommandMap(NewSlashCommand(applicationCommand, do, true, []string{"11111"})),
		}
		interaction := &discord.InteractionRequest{Data: &discord.ApplicationCommandInteractionData{
			Name: commandName,
			Options: []*discord.ApplicationCommandInteractionDataOption{
				{Name: optionName, Value: json.RawMessage(fmt.Sprintf(`"%s"`, optionVal))},
			},
		}}
		interactionBytes, err := json.Marshal(interaction)
		require.NoError(t, err)

		interactionRequest, err := handler.unmarshal(interactionBytes)
		require.NoError(t, err)
		require.Equal(t, optionVal, *interactionRequest.Data.Options[0].String)
	})
	t.Run("success/unmarshal integer option", func(t *testing.T) {
		optionName := "count"
		optionVal := 100
		applicationCommand := &discord.ApplicationCommand{
			Name:        commandName,
			Description: "desc",
			Options: []*discord.ApplicationCommandOption{
				{Name: optionName, Type: discord.ApplicationCommandOptionTypeInteger},
			},
		}
		handler := &Handler{
			Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
			SlashCommandMap: NewSlashCommandMap(NewSlashCommand(applicationCommand, do, true, []string{"11111"})),
		}

		interaction := &discord.InteractionRequest{Data: &discord.ApplicationCommandInteractionData{
			Name: commandName,
			Options: []*discord.ApplicationCommandInteractionDataOption{
				{Name: optionName, Value: json.RawMessage(fmt.Sprintf(`%d`, optionVal))},
			},
		}}
		interactionBytes, err := json.Marshal(interaction)
		require.NoError(t, err)

		interactionRequest, err := handler.unmarshal(interactionBytes)
		require.NoError(t, err)
		require.Equal(t, optionVal, *interactionRequest.Data.Options[0].Integer)
	})
	t.Run("success/unmarshal bool option", func(t *testing.T) {
		optionName := "enabled"
		optionVal := true
		applicationCommand := &discord.ApplicationCommand{
			Name:        commandName,
			Description: "desc",
			Options: []*discord.ApplicationCommandOption{
				{Name: optionName, Type: discord.ApplicationCommandOptionTypeBoolean},
			},
		}
		handler := &Handler{
			Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
			SlashCommandMap: NewSlashCommandMap(NewSlashCommand(applicationCommand, do, true, []string{"11111"})),
		}

		interaction := &discord.InteractionRequest{Data: &discord.ApplicationCommandInteractionData{
			Name: commandName,
			Options: []*discord.ApplicationCommandInteractionDataOption{
				{Name: optionName, Value: json.RawMessage(fmt.Sprintf(`%t`, optionVal))},
			},
		}}
		interactionBytes, err := json.Marshal(interaction)
		require.NoError(t, err)

		interactionRequest, err := handler.unmarshal(interactionBytes)
		require.NoError(t, err)
		require.Equal(t, optionVal, *interactionRequest.Data.Options[0].Boolean)
	})
	t.Run("success/unmarshal user option", func(t *testing.T) {
		optionName := "enabled"
		optionVal := &discord.User{ID: "0", Username: "wafer"}
		applicationCommand := &discord.ApplicationCommand{
			Name:        commandName,
			Description: "desc",
			Options: []*discord.ApplicationCommandOption{
				{Name: optionName, Type: discord.ApplicationCommandOptionTypeUser},
			},
		}
		handler := &Handler{
			Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
			SlashCommandMap: NewSlashCommandMap(NewSlashCommand(applicationCommand, do, true, []string{"11111"})),
		}
		rawJSON, err := json.Marshal(optionVal)
		require.NoError(t, err)
		interaction := &discord.InteractionRequest{Data: &discord.ApplicationCommandInteractionData{
			Name: commandName,
			Options: []*discord.ApplicationCommandInteractionDataOption{
				{Name: optionName, Value: json.RawMessage(string(rawJSON))},
			},
		}}
		interactionBytes, err := json.Marshal(interaction)
		require.NoError(t, err)

		interactionRequest, err := handler.unmarshal(interactionBytes)
		require.NoError(t, err)
		require.Equal(t, *optionVal, *interactionRequest.Data.Options[0].User)
	})
	t.Run("success/unmarshal role option", func(t *testing.T) {
		optionName := "enabled"
		optionVal := &discord.Role{ID: "0", Name: "admin"}
		applicationCommand := &discord.ApplicationCommand{
			Name:        commandName,
			Description: "desc",
			Options: []*discord.ApplicationCommandOption{
				{Name: optionName, Type: discord.ApplicationCommandOptionTypeRole},
			},
		}
		handler := &Handler{
			Creds:           &discord.Credentials{PublicKey: hex.EncodeToString(publicKey)},
			SlashCommandMap: NewSlashCommandMap(NewSlashCommand(applicationCommand, do, true, []string{"11111"})),
		}
		rawJSON, err := json.Marshal(optionVal)
		require.NoError(t, err)
		interaction := &discord.InteractionRequest{Data: &discord.ApplicationCommandInteractionData{
			Name: commandName,
			Options: []*discord.ApplicationCommandInteractionDataOption{
				{Name: optionName, Value: json.RawMessage(string(rawJSON))},
			},
		}}
		interactionBytes, err := json.Marshal(interaction)
		require.NoError(t, err)

		interactionRequest, err := handler.unmarshal(interactionBytes)
		require.NoError(t, err)
		require.Equal(t, *optionVal, *interactionRequest.Data.Options[0].Role)
	})
}

func TestVerify(t *testing.T) {
	body := "body"
	timestamp := "1500000000"
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		msg := []byte(timestamp + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey))
		require.True(t, res)
	})
	t.Run("failure/modified message parts", func(t *testing.T) {
		msg := []byte(timestamp + "baddata" + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey))
		require.False(t, res)
	})
	t.Run("failure/blank signature timestamp", func(t *testing.T) {
		msg := []byte(timestamp + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", "")
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey))
		require.False(t, res)
	})
	t.Run("failure/blank signature ed25519", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", "")

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey))
		require.False(t, res)
	})
	t.Run("failure/non-hex public key", func(t *testing.T) {
		msg := []byte(timestamp + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey)+"Z")
		require.False(t, res)
	})
	t.Run("failure/non-hex signature", func(t *testing.T) {
		msg := []byte(timestamp + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize])+"Z")

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey))
		require.False(t, res)
	})
	t.Run("failure/wrong length public key", func(t *testing.T) {
		msg := []byte(timestamp + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey)+"1111")
		require.False(t, res)
	})
	t.Run("failure/wrong length signature", func(t *testing.T) {
		msg := []byte(timestamp + body)
		signature := ed25519.Sign(privateKey, msg)
		headers := http.Header{}
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize])+"1111")

		res := verify([]byte(body), headers, hex.EncodeToString(publicKey))
		require.False(t, res)
	})
}

func getAuthHeaders(body string) map[string]string {
	timestamp := "1000000000"
	msg := []byte(timestamp + body)
	signature := ed25519.Sign(privateKey, msg)
	headers := map[string]string{
		"X-Signature-Timestamp": timestamp,
		"X-Signature-Ed25519":   hex.EncodeToString(signature[:ed25519.SignatureSize]),
	}
	return headers
}

func httpTestRequest(handler http.HandlerFunc, method string, url string, headers map[string]string, body string) ([]byte, *http.Response, error) {
	request := httptest.NewRequest(method, url, strings.NewReader(body))
	request.Header.Set("accept", "application/json")
	for key, val := range headers {
		request.Header.Set(key, val)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	testResponse := recorder.Result()
	responseBody, err := ioutil.ReadAll(recorder.Body)
	return responseBody, testResponse, err
}
