package disgoslash

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wafer-bw/disgoslash/discord"
	"github.com/wafer-bw/disgoslash/errs"
)

// Handler is used to handle Discord slash command interaction requests.
type Handler struct {
	SlashCommandMap SlashCommandMap
	Creds           *discord.Credentials
}

var pongResponse = &discord.InteractionResponse{
	Type: discord.InteractionResponseTypePong,
}

// Handle incoming interaction requests from Discord guilds
// executing the SlashCommand's Action and returning the
// interaction response.
func (handler *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	interactionRequest, err := handler.resolve(r)
	if err != nil {
		handler.respond(w, nil, err)
		return
	}

	interactionResponse, err := handler.triage(interactionRequest)
	if err != nil {
		handler.respond(w, nil, err)
		return
	}

	body, err := handler.marshal(interactionResponse)
	if err != nil {
		handler.respond(w, nil, err)
		return
	}

	handler.respond(w, body, nil)
}

func (handler *Handler) resolve(r *http.Request) (*discord.InteractionRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if !verify(body, r.Header, handler.Creds.PublicKey) {
		return nil, errs.ErrUnauthorized
	}

	return handler.unmarshal(body)
}

func (handler *Handler) triage(interaction *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	switch interaction.Type {
	case discord.InteractionTypePing:
		return pongResponse, nil
	case discord.InteractionTypeApplicationCommand:
		return handler.execute(interaction)
	default:
		return nil, errs.ErrInvalidInteractionType
	}
}

func (handler *Handler) execute(interaction *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	slashCommand, ok := handler.SlashCommandMap[interaction.Data.Name]
	if !ok {
		return nil, errs.ErrNotImplemented
	}
	return slashCommand.Action(interaction), nil
}

func (handler *Handler) respond(w http.ResponseWriter, body []byte, err error) {
	switch err {
	case nil:
		if _, err = w.Write(body); err != nil {
			handler.respond(w, nil, err)
		}
	case errs.ErrUnauthorized:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errs.ErrInvalidInteractionType:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errs.ErrNotImplemented:
		http.Error(w, err.Error(), http.StatusNotImplemented)
	default:
		log.Printf("ERROR: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *Handler) unmarshal(data []byte) (*discord.InteractionRequest, error) {
	interaction := &discord.InteractionRequest{}
	if err := json.Unmarshal(data, interaction); err != nil {
		return nil, err
	}
	return interaction, nil
}

func (handler *Handler) marshal(response *discord.InteractionResponse) ([]byte, error) {
	return json.Marshal(response)
}

func verify(rawBody []byte, headers http.Header, publicKey string) bool {
	signature := headers.Get("x-signature-ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	} else if len(sig) != ed25519.SignatureSize {
		return false
	}

	timestamp := headers.Get("x-signature-timestamp")
	if timestamp == "" {
		return false
	}

	keyBytes, err := hex.DecodeString(publicKey)
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
