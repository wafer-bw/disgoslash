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
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	require.Nil(t, err)

	t.Run("success", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey)
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authorizationImpl.Verify([]byte(body), headers)
		require.True(t, res)
	})

	t.Run("failure/modified message parts", func(t *testing.T) {
		msg := []byte(timestamp + "baddata" + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey)
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/blank signature timestamp", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey)
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", "")
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/blank signature ed25519", func(t *testing.T) {
		headers := http.Header{}
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey)
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", "")

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/non-hex public key", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey) + "Z"
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/non-hex signature", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey)
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize])+"Z")

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/wrong length public key", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey) + "1111"
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize]))

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})

	t.Run("failure/wrong length signature", func(t *testing.T) {
		msg := []byte(timestamp + body)
		headers := http.Header{}
		signature := ed25519.Sign(privkey, msg)
		conf := GetConf()
		conf.Credentials.PublicKey = hex.EncodeToString(pubkey)
		authorizationImpl := NewAuth(conf)

		headers.Set("X-Signature-Timestamp", timestamp)
		headers.Set("X-Signature-Ed25519", hex.EncodeToString(signature[:ed25519.SignatureSize])+"1111")

		res := authorizationImpl.Verify([]byte(body), headers)
		require.False(t, res)
	})
}
