package syncer

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/errs"
	clientMocks "github.com/wafer-bw/disgoslash/generatedmocks/client"
	"github.com/wafer-bw/disgoslash/mocks"
	"github.com/wafer-bw/disgoslash/models"
	"github.com/wafer-bw/disgoslash/slashcommands"
)

var clientMock = &clientMocks.Client{}
var syncerImpl = construct(&deps{client: clientMock}, mocks.Conf)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRun(t *testing.T) {
	applicationCommands := []*models.ApplicationCommand{
		{ID: "A", Name: "testCommandA", Description: "desc"},
		{ID: "B", Name: "testCommandB", Description: "desc"},
	}
	slashCommandMap := slashcommands.NewMap(
		slashcommands.New("testCommandA", applicationCommands[0], mocks.SlashCommandDo, true, []string{"12345"}),
		slashcommands.New("testCommandB", applicationCommands[1], mocks.SlashCommandDo, false, []string{"67890"}),
	)

	t.Run("success", func(t *testing.T) {
		clientMock.On("ListApplicationCommands", "").Return([]*models.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		clientMock.On("ListApplicationCommands", "12345").Return([]*models.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		clientMock.On("ListApplicationCommands", "67890").Return([]*models.ApplicationCommand{applicationCommands[1]}, nil).Times(1)

		clientMock.On("DeleteApplicationCommand", "", "A").Return(nil).Times(1)
		clientMock.On("DeleteApplicationCommand", "12345", "A").Return(nil).Times(1)
		clientMock.On("DeleteApplicationCommand", "67890", "B").Return(nil).Times(1)

		clientMock.On("CreateApplicationCommand", "", applicationCommands[0]).Return(nil).Times(1)
		clientMock.On("CreateApplicationCommand", "12345", applicationCommands[0]).Return(nil).Times(1)
		clientMock.On("CreateApplicationCommand", "67890", applicationCommands[1]).Return(nil).Times(1)

		errs := syncerImpl.Run([]string{"", "12345"}, slashCommandMap)
		require.Equal(t, 0, len(errs))
	})
	t.Run("failure/has errors", func(t *testing.T) {
		clientMock.On("ListApplicationCommands", "").Return([]*models.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		clientMock.On("ListApplicationCommands", "12345").Return([]*models.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		clientMock.On("ListApplicationCommands", "67890").Return(nil, errs.ErrForbidden).Times(1)

		clientMock.On("DeleteApplicationCommand", "", "A").Return(errs.ErrMaxRetries).Times(1)
		clientMock.On("DeleteApplicationCommand", "12345", "A").Return(nil).Times(1)

		clientMock.On("CreateApplicationCommand", "", applicationCommands[0]).Return(nil).Times(1)
		clientMock.On("CreateApplicationCommand", "12345", applicationCommands[0]).Return(nil).Times(1)
		clientMock.On("CreateApplicationCommand", "67890", applicationCommands[1]).Return(errs.ErrForbidden).Times(1)

		errs := syncerImpl.Run([]string{"", "12345"}, slashCommandMap)
		require.Equal(t, 3, len(errs))
	})
}
