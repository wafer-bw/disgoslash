package discord

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplicationCommandUnmarshal(t *testing.T) {
	t.Run("success/permissions omitempty", func(t *testing.T) {
		notContain := "default_permission"
		cmd := ApplicationCommand{Name: "name"}
		actual, err := json.Marshal(cmd)
		require.NoError(t, err)
		require.NotContains(t, string(actual), notContain)
	})
}

func TestApplicationCommandInteractionDataOption(t *testing.T) {
	t.Run("success/get string", func(t *testing.T) {
		expect := "abc"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`"%s"`, expect))}
		actual, err := cmd.GetString()
		require.NoError(t, err)
		require.Equal(t, expect, *actual)
	})
	t.Run("failure/get string", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(expect)}
		_, err := cmd.GetString()
		require.Error(t, err)
	})
	t.Run("success/get int", func(t *testing.T) {
		expect := 123
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`%d`, expect))}
		actual, err := cmd.GetInt()
		require.NoError(t, err)
		require.Equal(t, expect, *actual)
	})
	t.Run("failure/get string", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`"%s"`, expect))}
		_, err := cmd.GetInt()
		require.Error(t, err)
	})
	t.Run("success/get bool", func(t *testing.T) {
		expect := true
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`%t`, expect))}
		actual, err := cmd.GetBool()
		require.NoError(t, err)
		require.Equal(t, expect, *actual)
	})
	t.Run("failure/get string", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`"%s"`, expect))}
		_, err := cmd.GetBool()
		require.Error(t, err)
	})
	t.Run("success/get channel", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`"%s"`, expect))}
		actual, err := cmd.GetChannelID()
		require.NoError(t, err)
		require.Equal(t, expect, *actual)
	})
	t.Run("success/get role", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`"%s"`, expect))}
		actual, err := cmd.GetRoleID()
		require.NoError(t, err)
		require.Equal(t, expect, *actual)
	})
	t.Run("success/get user", func(t *testing.T) {
		expect := "123"
		cmd := ApplicationCommandInteractionDataOption{Value: json.RawMessage(fmt.Sprintf(`"%s"`, expect))}
		actual, err := cmd.GetUserID()
		require.NoError(t, err)
		require.Equal(t, expect, *actual)
	})
}
