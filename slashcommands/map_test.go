package slashcommands

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/mocks"
	"github.com/wafer-bw/disgoslash/models"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNewMap(t *testing.T) {
	slashCommandMap := NewMap(
		New(mocks.SlashCommandName, &models.ApplicationCommand{Name: mocks.SlashCommandName, Description: "desc"}, mocks.SlashCommandDo, true, []string{"11111"}),
	)
	require.Equal(t, 1, len(slashCommandMap))
}
