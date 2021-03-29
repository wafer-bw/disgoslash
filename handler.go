package disgoslash

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wafer-bw/disgoslash/discord"
	"github.com/wafer-bw/disgoslash/errs"
)

// Handler is used to handle Discord slash command interaction requests.
type Handler struct {
	SlashCommandMap SlashCommandMap
	Creds           *discord.Credentials
}

type response struct {
	w    http.ResponseWriter
	body []byte
	err  error
}

var pongResponse = &discord.InteractionResponse{
	Type: discord.InteractionResponseTypePong,
}

// Handle incoming interaction requests from Discord guilds,
// executing the SlashCommand's Action and responding with
// its InteractionResponse.
//
// 400 - An invalid Discord Interaction Type was passed in the request.
//
// 401 - Authorization failed.
//
// 500 - Something unexpected went wrong OR the Action did not respond
// within discord's maximum response time of 3 seconds.
//
// 501 - A SlashCommand that does not exist in the SlashCommandMap was
// requested.
func (handler *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	deadline := time.Now().Add(discord.MaxResponseTime)
	ctx, cancel := context.WithDeadline(r.Context(), deadline)
	defer cancel()

	responseChannel := make(chan response, 1)
	go handler.handle(responseChannel, w, r)
	select {
	case response := <-responseChannel:
		handler.respond(response)
	case <-ctx.Done():
		handler.respond(response{w: w, body: nil, err: ctx.Err()})
	}
}

func (handler *Handler) handle(ch chan response, w http.ResponseWriter, r *http.Request) {
	interactionRequest, err := handler.resolve(r)
	if err != nil {
		ch <- response{w: w, body: nil, err: err}
		return
	}

	interactionResponse, err := handler.execute(interactionRequest)
	if err != nil {
		ch <- response{w: w, body: nil, err: err}
		return
	}

	body, err := handler.marshal(interactionResponse)
	if err != nil {
		ch <- response{w: w, body: nil, err: err}
		return
	}

	ch <- response{w: w, body: body, err: nil}
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

func (handler *Handler) execute(interaction *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	switch interaction.Type {
	case discord.InteractionTypePing:
		return pongResponse, nil
	case discord.InteractionTypeApplicationCommand:
		return handler.doAction(interaction)
	default:
		return nil, errs.ErrInvalidInteractionType
	}
}

func (handler *Handler) doAction(interaction *discord.InteractionRequest) (*discord.InteractionResponse, error) {
	slashCommand, ok := handler.SlashCommandMap[interaction.Data.Name]
	if !ok {
		return nil, errs.ErrNotImplemented
	}
	response := slashCommand.Action(interaction)
	if response == nil {
		return nil, errs.ErrNilInteractionResponse
	}
	return response, nil
}

func (handler *Handler) respond(resp response) {
	if resp.err != nil {
		log.Println(resp.err)
	}

	switch resp.err {
	case nil:
		if _, err := resp.w.Write(resp.body); resp.err != nil {
			handler.respond(response{w: resp.w, body: nil, err: err})
		}
	case errs.ErrInvalidInteractionType:
		http.Error(resp.w, resp.err.Error(), http.StatusBadRequest)
	case errs.ErrUnauthorized:
		http.Error(resp.w, resp.err.Error(), http.StatusUnauthorized)
	case errs.ErrNotImplemented:
		http.Error(resp.w, resp.err.Error(), http.StatusNotImplemented)
	default:
		http.Error(resp.w, resp.err.Error(), http.StatusInternalServerError)
	}
}

func (handler *Handler) unmarshal(data []byte) (*discord.InteractionRequest, error) {
	interaction := &discord.InteractionRequest{}
	if err := json.Unmarshal(data, interaction); err != nil {
		return nil, err
	}

	if interaction.Data == nil {
		return interaction, nil
	}

	slashCommand, ok := handler.SlashCommandMap[interaction.Data.Name]
	if !ok {
		return nil, errs.ErrNotImplemented
	}

	handler.unmarshalOptions(slashCommand.ApplicationCommand.Options, interaction.Data.Options)
	return interaction, nil
}

func (handler *Handler) unmarshalOptions(commandOptions []*discord.ApplicationCommandOption, interactionOptions []*discord.ApplicationCommandInteractionDataOption) {
	for i, commandOption := range commandOptions {
		interactionOption := interactionOptions[i]
		handler.unmarshalOption(commandOption, interactionOption)
	}
}

func (handler *Handler) unmarshalOption(commandOption *discord.ApplicationCommandOption, interactionOption *discord.ApplicationCommandInteractionDataOption) {
	switch commandOption.Type {
	case discord.ApplicationCommandOptionTypeInteger:
		interactionOption.Integer = new(int)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.Integer); err != nil {
			log.Println("error:", err)
		}
	case discord.ApplicationCommandOptionTypeBoolean:
		interactionOption.Boolean = new(bool)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.Boolean); err != nil {
			log.Println("error:", err)
		}
	case discord.ApplicationCommandOptionTypeString:
		interactionOption.String = new(string)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.String); err != nil {
			fmt.Println("error:", err)
		}
	case discord.ApplicationCommandOptionTypeUser:
		interactionOption.User = new(discord.User)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.User); err != nil {
			fmt.Println("error:", err)
		}
	case discord.ApplicationCommandOptionTypeRole:
		interactionOption.Role = new(discord.Role)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.Role); err != nil {
			fmt.Println("error:", err)
		}
	case discord.ApplicationCommandOptionTypeSubCommand:
		interactionOption.SubCommand = new(string)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.SubCommand); err != nil {
			fmt.Println("error:", err)
		}
	case discord.ApplicationCommandOptionTypeSubCommandGroup:
		interactionOption.SubCommandGroup = new(string)
		if err := json.Unmarshal(interactionOption.Value, interactionOption.SubCommandGroup); err != nil {
			fmt.Println("error:", err)
		}
	}
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
