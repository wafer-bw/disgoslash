package disgoslash

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var guildID = "1234567890"
var mockClient = &mockClientInterface{}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCoverMockRequest(t *testing.T) {
	// This covers the generated request method inside client_mock.go
	mockClient.On("request", http.MethodGet, "", nil).Return(http.StatusOK, nil, nil)
	status, data, err := mockClient.request(http.MethodGet, "", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.Nil(t, data)
}
