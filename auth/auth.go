package auth

import (
	"crypto/ed25519"
	"encoding/hex"
	"net/http"

	"github.com/wafer-bw/disgoslash/config"
)

// Deps defines `Authorization` dependencies
type Deps struct{}

// impl implements `Authorization` properties
type impl struct {
	deps *Deps
	conf *config.Config
}

// Authorization interfaces `Authorization` methods
type Authorization interface {
	Verify(rawBody []byte, headers http.Header) bool
}

// New returns a new `Authorization` interface
func New(deps *Deps, conf *config.Config) Authorization {
	return &impl{deps: deps, conf: conf}
}

// Verify verifies that requests from discord are authorized using ed25519
// https://discord.com/developers/docs/interactions/slash-commands#security-and-authorization
func (impl *impl) Verify(rawBody []byte, headers http.Header) bool {
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

	keyBytes, err := hex.DecodeString(impl.conf.Credentials.PublicKey)
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
