package disgoslash

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/models"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSlashCommandDo(t *testing.T) {
	res, err := SlashCommandDo(&models.InteractionRequest{})
	require.NoError(t, err)
	require.Equal(t, InteractionResponse, res)
}
