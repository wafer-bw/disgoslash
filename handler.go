package disgoslash

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// handler implements a `Handler` interface's properties
type handler struct {
	slashCommandMap SlashCommandMap
	auth            Auth
}

// Handler interfaces `Handler` methods
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// New returns a new `Handler` interface
func New(slashCommandMap SlashCommandMap, creds *Credentials) Handler {
	return constructHandler(NewAuth(creds), slashCommandMap)
}

func constructHandler(auth Auth, slashCommandMap SlashCommandMap) Handler {
	return &handler{auth: auth, slashCommandMap: slashCommandMap}
}

var pongResponse = &InteractionResponse{
	Type: InteractionResponseTypePong,
}

// Handle handles incoming HTTP requests
func (handler *handler) Handle(w http.ResponseWriter, r *http.Request) {
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

func (handler *handler) resolve(r *http.Request) (*InteractionRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if !handler.auth.Verify(body, r.Header) {
		return nil, ErrUnauthorized
	}

	return handler.unmarshal(body)
}

func (handler *handler) triage(interaction *InteractionRequest) (*InteractionResponse, error) {
	switch interaction.Type {
	case InteractionTypePing:
		return pongResponse, nil
	case InteractionTypeApplicationCommand:
		return handler.execute(interaction)
	default:
		return nil, ErrInvalidInteractionType
	}
}

func (handler *handler) execute(interaction *InteractionRequest) (*InteractionResponse, error) {
	slashCommand, ok := handler.slashCommandMap[interaction.Data.Name]
	if !ok {
		return nil, ErrNotImplemented
	}
	return slashCommand.Do(interaction)
}

func (handler *handler) respond(w http.ResponseWriter, body []byte, err error) {
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

func (handler *handler) unmarshal(data []byte) (*InteractionRequest, error) {
	interaction := &InteractionRequest{}
	if err := json.Unmarshal(data, interaction); err != nil {
		return nil, err
	}
	return interaction, nil
}

func (handler *handler) marshal(response *InteractionResponse) ([]byte, error) {
	return json.Marshal(response)
}
