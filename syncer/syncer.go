package syncer

import (
	"log"

	"github.com/wafer-bw/disgoslash/client"
	"github.com/wafer-bw/disgoslash/config"
	"github.com/wafer-bw/disgoslash/slashcommands"
)

// deps defines `Syncer` dependencies
type deps struct {
	client client.Client
}

// impl implements `Syncer` properties
type impl struct {
	deps *deps
	conf *config.Config
}

// Syncer interfaces `Syncer` methods
type Syncer interface {
	Run(guildIDs []string, slashCommandMap slashcommands.Map) []error
}

type unregisterTarget struct {
	guildID   string
	commandID string
	name      string
}

// New returns a new `Syncer` interface
func New(creds *config.Credentials) Syncer {
	conf := config.New(creds)
	client := client.New(creds)
	return construct(&deps{client: client}, conf)
}

func construct(deps *deps, conf *config.Config) Syncer {
	return &impl{deps: deps, conf: conf}
}

// Run will reregister all of the provided slash commands
func (impl *impl) Run(guildIDs []string, slashCommandMap slashcommands.Map) []error {
	allErrs := []error{}
	unregisterTargets, errs := impl.getCommandsToUnregister(guildIDs, slashCommandMap)
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	allErrs = append(allErrs, impl.unregisterCommands(unregisterTargets)...)
	allErrs = append(allErrs, impl.registerCommands(slashCommandMap)...)
	return allErrs
}

func (impl *impl) getCommandsToUnregister(guildIDs []string, commandMap slashcommands.Map) ([]unregisterTarget, []error) {
	errs := []error{}
	log.Println("Collecting outdated commands...")
	uniqueGuildIDs := impl.getUniqueGuildIDs(guildIDs, commandMap)
	unregisterTargets := []unregisterTarget{}
	for _, guildID := range uniqueGuildIDs {
		log.Printf("\t- Guild: %s\n", guildText(guildID))
		commands, err := impl.deps.client.ListApplicationCommands(guildID)
		if err != nil {
			log.Printf("\t\t- ERROR: %s\n", err.Error())
			errs = append(errs, err)
		} else {
			log.Printf("\t\t- SUCCESS")
		}
		for _, command := range commands {
			unregisterTargets = append(unregisterTargets, unregisterTarget{
				guildID:   guildID,
				commandID: command.ID,
				name:      command.Name,
			})
		}
	}
	return unregisterTargets, errs
}

func (impl *impl) registerCommands(commandMap slashcommands.Map) []error {
	errs := []error{}
	log.Println("Registering new commands...")
	for _, command := range commandMap {
		for _, guildID := range command.GuildIDs {
			log.Printf("\t- Guild: %s, Command: %s\n", guildText(guildID), command.Name)
			err := impl.deps.client.CreateApplicationCommand(guildID, command.AppCommand)
			if err != nil {
				log.Printf("\t\t- ERROR: %s\n", err.Error())
				errs = append(errs, err)
			} else {
				log.Printf("\t\t- SUCCESS")
			}
		}
	}
	return errs
}

func (impl *impl) unregisterCommands(unregisterTargets []unregisterTarget) []error {
	errs := []error{}
	log.Println("Unregistering outdated commands...")
	for _, target := range unregisterTargets {
		log.Printf("\t- Guild: %s, Command: %s\n", guildText(target.guildID), target.name)
		err := impl.deps.client.DeleteApplicationCommand(target.guildID, target.commandID)
		if err != nil {
			log.Printf("\t\t- ERROR: %s\n", err.Error())
			errs = append(errs, err)
		} else {
			log.Printf("\t\t- SUCCESS")
		}
	}
	return errs
}

func (impl *impl) getUniqueGuildIDs(guildIDs []string, commands slashcommands.Map) []string {
	uniqueGuildIDsMap := map[string]struct{}{
		"": {}, // include global
	}
	for _, id := range guildIDs {
		if _, ok := uniqueGuildIDsMap[id]; !ok {
			uniqueGuildIDsMap[id] = struct{}{}
		}
	}
	for _, command := range commands {
		for _, guildID := range command.GuildIDs {
			if _, ok := uniqueGuildIDsMap[guildID]; !ok {
				uniqueGuildIDsMap[guildID] = struct{}{}
			}
		}
	}
	uniqueGuildIDs := []string{}
	for id := range uniqueGuildIDsMap {
		uniqueGuildIDs = append(uniqueGuildIDs, id)
	}
	return uniqueGuildIDs
}

func guildText(guildID string) string {
	if guildID == "" {
		return "GLOBAL"
	}
	return guildID
}
