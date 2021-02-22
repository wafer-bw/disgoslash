package disgoslash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		require.NotPanics(t, func() {
			NewConfig(&Credentials{
				Token:     "a",
				ClientID:  "b",
				PublicKey: "c",
			})
		})
	})
	t.Run("failure/panics", func(t *testing.T) {
		require.Panics(t, func() { NewConfig(nil) })
	})
}

func TestFindBlankEnvVars(t *testing.T) {
	blanks := findBlankEnvVars(envVars{ClientID: "123abc"})
	for _, b := range blanks {
		require.NotEqual(t, "ClientID", b)
	}
}

func TestGetEnvVars(t *testing.T) {
	require.Panics(t, func() { getEnvVars() })
}

func TestHaveNoBlankEnvVars(t *testing.T) {
	require.Panics(t, func() { ensureNoBlankEnvVars(envVars{}) })
}
