package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	authMocks "github.com/wafer-bw/disgoslash/generatedmocks/auth"
	"github.com/wafer-bw/disgoslash/mocks"
	"github.com/wafer-bw/disgoslash/models"
	"github.com/wafer-bw/disgoslash/slashcommands"
)

var url = "http://localhost/api"
var authMock = &authMocks.Authorization{}
var interactionName = "interaction"
var handlerImpl = construct(&deps{
	auth: authMock,
	slashCommandsMap: slashcommands.NewMap(
		slashcommands.New(interactionName, &models.ApplicationCommand{Name: interactionName, Description: "desc"}, mocks.SlashCommandDo, true, []string{"11111"}),
	),
}, mocks.Conf)
var handler = http.HandlerFunc(handlerImpl.Handle)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNew(t *testing.T) {
	require.NotNil(t, handlerImpl)
	require.IsType(t, &impl{}, handlerImpl)
}

func TestHandle(t *testing.T) {
	headers := map[string]string{"Accept": "application/json"}
	t.Run("success/respond to ping", func(t *testing.T) {
		authMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		body, resp, err := httpRequest(http.MethodGet, url, headers, `{"type": 1}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})
	t.Run("success/run interaction", func(t *testing.T) {
		authMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		interaction := &models.InteractionRequest{
			Type: models.InteractionTypeApplicationCommand,
			Data: &models.ApplicationCommandInteractionData{Name: interactionName},
		}
		data, err := json.Marshal(interaction)
		require.NoError(t, err)
		body, resp, err := httpRequest(http.MethodGet, url, headers, string(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})
	t.Run("failure/unimplemented interaction", func(t *testing.T) {
		authMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		interaction := &models.InteractionRequest{
			Type: models.InteractionTypeApplicationCommand,
			Data: &models.ApplicationCommandInteractionData{Name: interactionName + "Z"},
		}
		data, err := json.Marshal(interaction)
		require.NoError(t, err)
		body, resp, err := httpRequest(http.MethodGet, url, headers, string(data))
		require.NoError(t, err)
		require.Equal(t, http.StatusNotImplemented, resp.StatusCode, string(body))
	})
	t.Run("failure/invalid interaction type", func(t *testing.T) {
		authMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true).Times(1)
		body, resp, err := httpRequest(http.MethodGet, url, headers, `{"type": 3}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, string(body))
	})
	t.Run("failure/unauthorized", func(t *testing.T) {
		authMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(false).Times(1)
		body, resp, err := httpRequest(http.MethodGet, url, headers, `{"type": 1}`)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode, string(body))
	})

}

func httpRequest(method string, url string, headers map[string]string, body string) ([]byte, *http.Response, error) {
	request := httptest.NewRequest(method, url, strings.NewReader(body))
	for key, val := range headers {
		request.Header.Set(key, val)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, err := ioutil.ReadAll(recorder.Body)
	return responseBody, response, err
}
