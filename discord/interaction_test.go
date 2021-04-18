package discord

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplicationCommandInteractionDataOption(t *testing.T) {
	t.Run("success/get string", func(t *testing.T) {
		expect := "abc"
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		actual, ok := cmd.StringValue()
		require.True(t, ok)
		require.Equal(t, expect, actual)
	})
	t.Run("success/get int", func(t *testing.T) {
		expect := float64(123)
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		actual, ok := cmd.IntValue()
		require.True(t, ok)
		require.Equal(t, int(expect), actual)
	})
	t.Run("failure/get int", func(t *testing.T) {
		expect := "abc"
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		_, ok := cmd.IntValue()
		require.False(t, ok)
	})
	t.Run("success/get bool", func(t *testing.T) {
		expect := true
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		actual, ok := cmd.BoolValue()
		require.True(t, ok)
		require.Equal(t, expect, actual)
	})
	t.Run("success/get channel id", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		actual, ok := cmd.ChannelIDValue()
		require.True(t, ok)
		require.Equal(t, expect, actual)
	})
	t.Run("success/get role id", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		actual, ok := cmd.RoleIDValue()
		require.True(t, ok)
		require.Equal(t, expect, actual)
	})
	t.Run("success/get user id", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: expect}
		actual, ok := cmd.UserIDValue()
		require.True(t, ok)
		require.Equal(t, expect, actual)
	})
}
