package disgoslash

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Handler implements a `Handler` interface's properties
type Handler struct {
	SlashCommandMap SlashCommandMap
	PublicKey       string
}

var pongResponse = &InteractionResponse{
	Type: InteractionResponseTypePong,
}

// Handle handles incoming HTTP requests
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

func (handler *Handler) resolve(r *http.Request) (*InteractionRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if !verify(body, r.Header, handler.PublicKey) {
		return nil, ErrUnauthorized
	}

	return handler.unmarshal(body)
}

func (handler *Handler) triage(interaction *InteractionRequest) (*InteractionResponse, error) {
	switch interaction.Type {
	case InteractionTypePing:
		return pongResponse, nil
	case InteractionTypeApplicationCommand:
		return handler.execute(interaction)
	default:
		return nil, ErrInvalidInteractionType
	}
}

func (handler *Handler) execute(interaction *InteractionRequest) (*InteractionResponse, error) {
	slashCommand, ok := handler.SlashCommandMap[interaction.Data.Name]
	if !ok {
		return nil, ErrNotImplemented
	}
	return slashCommand.Do(interaction)
}

func (handler *Handler) respond(w http.ResponseWriter, body []byte, err error) {
	switch err {
	case nil:
		if _, err = w.Write(body); err != nil {
			handler.respond(w, nil, err)
		}
	case ErrUnauthorized:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case ErrInvalidInteractionType:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrNotImplemented:
		http.Error(w, err.Error(), http.StatusNotImplemented)
	default:
		log.Printf("ERROR: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *Handler) unmarshal(data []byte) (*InteractionRequest, error) {
	interaction := &InteractionRequest{}
	if err := json.Unmarshal(data, interaction); err != nil {
		return nil, err
	}
	return interaction, nil
}

func (handler *Handler) marshal(response *InteractionResponse) ([]byte, error) {
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
