package disgoslash

import (
	"crypto/ed25519"
	"encoding/hex"
	"net/http"
)

// auth implements an `Auth` interface's properties
type auth struct {
	publicKey string
}

// Auth interfaces `Auth` methods
type Auth interface {
	Verify(rawBody []byte, headers http.Header) bool
}

// NewAuth returns a new `Auth` interface
func NewAuth(publicKey string) Auth {
	return &auth{publicKey: publicKey}
}

// Verify that the request from Discord is authorized using ed25519
// https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
func (auth *auth) Verify(rawBody []byte, headers http.Header) bool {
	signature := headers.Get("x-signature-ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp := headers.Get("x-signature-timestamp")
	if timestamp == "" {
		return false
	}

	keyBytes, err := hex.DecodeString(auth.publicKey)
	if err != nil {
		return false
	}

	key := ed25519.PublicKey(keyBytes)
	if len(key) != 32 {
		return false
	}

	msg := []byte(timestamp + string(rawBody))
	return ed25519.Verify(key, msg, sig)
}
