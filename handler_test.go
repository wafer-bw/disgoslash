package disgoslash

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var url = "http://localhost/api"
var mockAuth = &MockAuth{}
var interactionName = "interaction"
var handlerImpl = constructHandler(mockAuth, NewSlashCommandMap(
	NewSlashCommand(interactionName, &ApplicationCommand{Name: interactionName, Description: "desc"}, SlashCommandDo, true, []string{"11111"}),
), Conf)
var handlerFunc = http.HandlerFunc(handlerImpl.Handle)

func TestNewHandler(t *testing.T) {
	require.NotNil(t, handlerImpl)
	require.IsType(t, &handler{}, handlerImpl)
}

func TestHandle(t *testing.T) {
	headers := map[string]string{"Accept": "application/json"}
	t.Run("success/respond to ping", func(t *testing.T) {
		mockAuth.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		body, resp, err := httpTestRequest(http.MethodGet, url, headers, `{"type": 1}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})
	t.Run("success/run interaction", func(t *testing.T) {
		mockAuth.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		interaction := &InteractionRequest{
			Type: InteractionTypeApplicationCommand,
			Data: &ApplicationCommandInteractionData{Name: interactionName},
		}
		data, err := json.Marshal(interaction)
		require.NoError(t, err)
		body, resp, err := httpTestRequest(http.MethodGet, url, headers, string(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})
	t.Run("failure/unimplemented interaction", func(t *testing.T) {
		mockAuth.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		interaction := &InteractionRequest{
			Type: InteractionTypeApplicationCommand,
			Data: &ApplicationCommandInteractionData{Name: interactionName + "Z"},
		}
		data, err := json.Marshal(interaction)
		require.NoError(t, err)
		body, resp, err := httpTestRequest(http.MethodGet, url, headers, string(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusNotImplemented, resp.StatusCode, string(body))
	})
	t.Run("failure/invalid interaction type", func(t *testing.T) {
		mockAuth.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		body, resp, err := httpTestRequest(http.MethodGet, url, headers, `{"type": 3}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, string(body))
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		mockAuth.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(false).Times(1)
		body, resp, err := httpTestRequest(http.MethodGet, url, headers, `{"type": 1}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode, string(body))
	})

}

func httpTestRequest(method string, url string, headers map[string]string, body string) ([]byte, *http.Response, error) {
	request := httptest.NewRequest(method, url, strings.NewReader(body))
	for key, val := range headers {
		request.Header.Set(key, val)
	}
	recorder := httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, err := ioutil.ReadAll(recorder.Body)
	return responseBody, response, err
}
