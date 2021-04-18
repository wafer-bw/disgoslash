package disgoslash

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wafer-bw/disgoslash/discord"
)

func TestSync(t *testing.T) {
	testResponse := &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{Content: "Hello World!"},
	}
	do := func(request *discord.InteractionRequest) *discord.InteractionResponse {
		return testResponse
	}

	applicationCommands := []*discord.ApplicationCommand{
		{ID: "A", Name: "testCommandA", Description: "desc"},
		{ID: "B", Name: "testCommandB", Description: "desc"},
	}
	slashCommandMap := NewSlashCommandMap(
		NewSlashCommand(applicationCommands[0], do, true, []string{"12345"}),
		NewSlashCommand(applicationCommands[1], do, false, []string{"67890"}),
	)

	t.Run("success", func(t *testing.T) {
		syncer := &Syncer{SlashCommandMap: slashCommandMap, GuildIDs: []string{"12345"}, client: mockClient}

		mockClient.On("list", "").Return([]*discord.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("list", "12345").Return([]*discord.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("list", "67890").Return([]*discord.ApplicationCommand{applicationCommands[1]}, nil).Times(1)

		mockClient.On("delete", "", "A").Return(nil).Times(1)
		mockClient.On("delete", "12345", "A").Return(nil).Times(1)
		mockClient.On("delete", "67890", "B").Return(nil).Times(1)

		mockClient.On("create", "", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("create", "12345", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("create", "67890", applicationCommands[1]).Return(nil).Times(1)

		errs := syncer.Sync()
		require.Equal(t, 0, len(errs))
	})
	t.Run("failure/has errors", func(t *testing.T) {
		syncer := &Syncer{SlashCommandMap: slashCommandMap, GuildIDs: []string{"", "12345"}, client: mockClient}

		mockClient.On("list", "").Return([]*discord.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("list", "12345").Return([]*discord.ApplicationCommand{applicationCommands[0]}, nil).Times(1)
		mockClient.On("list", "67890").Return(nil, ErrForbidden).Times(1)

		mockClient.On("delete", "", "A").Return(ErrMaxRetries).Times(1)
		mockClient.On("delete", "12345", "A").Return(nil).Times(1)

		mockClient.On("create", "", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("create", "12345", applicationCommands[0]).Return(nil).Times(1)
		mockClient.On("create", "67890", applicationCommands[1]).Return(ErrForbidden).Times(1)

		errs := syncer.Sync()
		require.Equal(t, 3, len(errs))
	})
}
