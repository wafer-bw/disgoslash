package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wafer-bw/disgoslash/auth"
	"github.com/wafer-bw/disgoslash/config"
	"github.com/wafer-bw/disgoslash/errs"
	"github.com/wafer-bw/disgoslash/models"
	"github.com/wafer-bw/disgoslash/slashcommands"
)

// deps defines `Handler` dependencies
type deps struct {
	slashCommandsMap slashcommands.Map
	auth             auth.Authorization
}

// impl implements `Handler` properties
type impl struct {
	deps *deps
	conf *config.Config
}

// Handler interfaces `Handler` methods
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// New returns a new `Handler` interface
func New(slashCommandMap slashcommands.Map, creds *config.Credentials) Handler {
	conf := config.New(creds)
	return construct(&deps{
		auth:             auth.New(&auth.Deps{}, conf),
		slashCommandsMap: slashCommandMap,
	}, conf)
}

// construct a new `Handler` interface
func construct(deps *deps, conf *config.Config) Handler {
	return &impl{deps: deps, conf: conf}
}

var pongResponse = &models.InteractionResponse{
	Type: models.InteractionResponseTypePong,
}

// Handle handles incoming HTTP requests
func (impl *impl) Handle(w http.ResponseWriter, r *http.Request) {
	interactionRequest, err := impl.resolve(r)
	if err != nil {
		impl.respond(w, nil, err)
		return
	}

	interactionResponse, err := impl.triage(interactionRequest)
	if err != nil {
		impl.respond(w, nil, err)
		return
	}

	body, err := impl.marshal(interactionResponse)
	if err != nil {
		impl.respond(w, nil, err)
		return
	}

	impl.respond(w, body, nil)
}

func (impl *impl) resolve(r *http.Request) (*models.InteractionRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if !impl.deps.auth.Verify(body, r.Header) {
		return nil, errs.ErrUnauthorized
	}

	return impl.unmarshal(body)
}

func (impl *impl) triage(interaction *models.InteractionRequest) (*models.InteractionResponse, error) {
	switch interaction.Type {
	case models.InteractionTypePing:
		return pongResponse, nil
	case models.InteractionTypeApplicationCommand:
		return impl.execute(interaction)
	default:
		return nil, errs.ErrInvalidInteractionType
	}
}

func (impl *impl) execute(interaction *models.InteractionRequest) (*models.InteractionResponse, error) {
	slashCommand, ok := impl.deps.slashCommandsMap[interaction.Data.Name]
	if !ok {
		return nil, errs.ErrNotImplemented
	}
	return slashCommand.Do(interaction)
}

func (impl *impl) respond(w http.ResponseWriter, body []byte, err error) {
	switch err {
	case nil:
		if _, err = w.Write(body); err != nil {
			impl.respond(w, nil, err)
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

func (impl *impl) unmarshal(data []byte) (*models.InteractionRequest, error) {
	interaction := &models.InteractionRequest{}
	if err := json.Unmarshal(data, interaction); err != nil {
		return nil, err
	}
	return interaction, nil
}

func (impl *impl) marshal(response *models.InteractionResponse) ([]byte, error) {
	return json.Marshal(response)
}
