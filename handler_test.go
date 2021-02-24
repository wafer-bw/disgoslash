package disgoslash

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/discord"
)

//nolint
var publicKey, privateKey, _ = ed25519.GenerateKey(nil)
var url = "http://localhost/api"
var interactionName = "interaction"
var response = &discord.InteractionResponse{
	Type: discord.InteractionResponseTypeChannelMessageWithSource,
	Data: &discord.InteractionApplicationCommandCallbackData{Content: "Hello World!"},
}
var do = func(request *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	return response, nil
}
var handler = &Handler{
	PublicKey:       hex.EncodeToString(publicKey),
	SlashCommandMap: NewSlashCommandMap(NewSlashCommand(interactionName, &discord.ApplicationCommand{Name: interactionName, Description: "desc"}, do, true, []string{"11111"})),
}
var handlerFunc = http.HandlerFunc(handler.Handle)

func TestHandle(t *testing.T) {

	t.Run("success/respond to ping", func(t *testing.T) {
		requestBody := `{"type":1}`
		body, resp, err := httpTestRequest(http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
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

		body, resp, err := httpTestRequest(http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})

	t.Run("failure/unimplemented interaction", func(t *testing.T) {
		data, err := json.Marshal(&discord.InteractionRequest{
			Type: discord.InteractionTypeApplicationCommand,
			Data: &discord.ApplicationCommandInteractionData{Name: interactionName + "Z"},
		})
		require.NoError(t, err)
		requestBody := string(data)

		body, resp, err := httpTestRequest(http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotImplemented, resp.StatusCode, string(body))
	})

	t.Run("failure/invalid interaction type", func(t *testing.T) {
		requestBody := `{"type": 3}`

		body, resp, err := httpTestRequest(http.MethodGet, url, getAuthHeaders(requestBody), requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, string(body))
	})

	t.Run("failure/unauthorized", func(t *testing.T) {
		requestBody := `{"type": 1}`
		headers := getAuthHeaders(requestBody)
		headers["X-Signature-Ed25519"] = "X"

		body, resp, err := httpTestRequest(http.MethodGet, url, headers, requestBody)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode, string(body))
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

func httpTestRequest(method string, url string, headers map[string]string, body string) ([]byte, *http.Response, error) {
	request := httptest.NewRequest(method, url, strings.NewReader(body))
	request.Header.Set("accept", "application/json")
	for key, val := range headers {
		request.Header.Set(key, val)
	}
	recorder := httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, err := ioutil.ReadAll(recorder.Body)
	return responseBody, response, err
}
