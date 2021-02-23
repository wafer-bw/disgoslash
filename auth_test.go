package disgoslash

import (
	"crypto/ed25519"
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// todo - TestNewAuth()

func TestVerify(t *testing.T) {
	body := "body"
	timestamp := "1500000000"
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey))

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authImpl.Verify([]byte(body), headers)
		require.True(t, res)
	})

	t.Run("failure/modified message parts", func(t *testing.T) {
		msg := []byte(timestamp + "baddata" + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey))

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/blank signature timestamp", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey))

		headers.Set("X-Signature-Timestamp", "")
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/blank signature ed25519", func(t *testing.T) {
		headers := http.Header{}
		authImpl := NewAuth(hex.EncodeToString(publicKey))

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", "")

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/non-hex public key", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey) + "Z")

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/non-hex signature", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey))

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize])+"Z")

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/wrong length public key", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey) + "1111")

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/wrong length signature", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privateKey, msg)
		authImpl := NewAuth(hex.EncodeToString(publicKey))
		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize])+"1111")

		res := authImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})
}
