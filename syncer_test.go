package disgoslash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var mockClient = &MockClient{}
var syncerImpl = constructSyncer(mockClient, GetConf())

// todo - TestNewSyncer()

func TestRun(t *testing.T) {
	applicationCommands := []*ApplicationCommand{
		{ID: "A", Name: "testCommandA", Description: "desc"},
		{ID: "B", Name: "testCommandB", Description: "desc"},
	}
	slashCommandMap := NewSlashCommandMap(
		NewSlashCommand("testCommandA", applicationCommands[0], do, true, []string{"12345"}),
		NewSlashCommand("testCommandB", applicationCommands[1], do, false, []string{"67890"}),
	)

	t.Run("success", func(t *testing.T) {
		mockClient.On("ListApplicationCommands", "").Return([]*ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("ListApplicationCommands", "12345").Return([]*ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("ListApplicationCommands", "67890").Return([]*ApplicationCommand{applicationCommands[1]}, nil).Times(1)

		mockClient.On("DeleteApplicationCommand", "", "A").Return(nil).Times(1)
		mockClient.On("DeleteApplicationCommand", "12345", "A").Return(nil).Times(1)
		mockClient.On("DeleteApplicationCommand", "67890", "B").Return(nil).Times(1)

		mockClient.On("CreateApplicationCommand", "", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("CreateApplicationCommand", "12345", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("CreateApplicationCommand", "67890", applicationCommands[1]).Return(nil).Times(1)

		errs := syncerImpl.Run([]string{"", "12345"}, slashCommandMap)
		require.Equal(t, 0, len(errs))
	})
	t.Run("failure/has errors", func(t *testing.T) {
		mockClient.On("ListApplicationCommands", "").Return([]*ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("ListApplicationCommands", "12345").Return([]*ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("ListApplicationCommands", "67890").Return(nil, ErrForbidden).Times(1)

		mockClient.On("DeleteApplicationCommand", "", "A").Return(ErrMaxRetries).Times(1)
		mockClient.On("DeleteApplicationCommand", "12345", "A").Return(nil).Times(1)

		mockClient.On("CreateApplicationCommand", "", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("CreateApplicationCommand", "12345", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("CreateApplicationCommand", "67890", applicationCommands[1]).Return(ErrForbidden).Times(1)

		errs := syncerImpl.Run([]string{"", "12345"}, slashCommandMap)
		require.Equal(t, 3, len(errs))
	})
}
