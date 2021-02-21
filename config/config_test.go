package config

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		require.NotPanics(t, func() {
			New(&Credentials{
				Token:     "a",
				ClientID:  "b",
				PublicKey: "c",
			})
		})
	})
	t.Run("failure/panics", func(t *testing.T) {
		require.Panics(t, func() { New(nil) })
	})
}

func TestFindBlankEnvVars(t *testing.T) {
	blanks := findBlankEnvVars(EnvVars{ClientID: "123abc"})
	for _, b := range blanks {
		require.NotEqual(t, "ClientID", b)
	}
}

func TestGetEnvVars(t *testing.T) {
	require.Panics(t, func() { getEnvVars() })
}

func TestHaveNoBlankEnvVars(t *testing.T) {
	require.Panics(t, func() { ensureNoBlankEnvVars(EnvVars{}) })
}
