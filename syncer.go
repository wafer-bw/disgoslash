package disgoslash

import (
	"log"
)

// syncer implements a `Syncer` interface's properties
type syncer struct {
	client Client
	conf   *Config
}

// Syncer interface methods
type Syncer interface {
	Run(guildIDs []string, slashCommandMap SlashCommandMap) []error
}

type unregisterTarget struct {
	guildID   string
	commandID string
	name      string
}

// NewSyncer creates a new `Syncer` instance
func NewSyncer(creds *Credentials) Syncer {
	conf := NewConfig(creds)
	client := NewClient(creds)
	return constructSyncer(client, conf)
}

func constructSyncer(client Client, conf *Config) Syncer {
	return &syncer{client: client, conf: conf}
}

// Run will reregister all of the provided slash commands
func (syncer *syncer) Run(guildIDs []string, slashCommandMap SlashCommandMap) []error {
	allErrs := []error{}
	unregisterTargets, errs := syncer.getCommandsToUnregister(guildIDs, slashCommandMap)
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	allErrs = append(allErrs, syncer.unregisterCommands(unregisterTargets)...)
	allErrs = append(allErrs, syncer.registerCommands(slashCommandMap)...)
	return allErrs
}

func (syncer *syncer) getCommandsToUnregister(guildIDs []string, commandMap SlashCommandMap) ([]unregisterTarget, []error) {
	errs := []error{}
	log.Println("Collecting outdated commands...")
	uniqueGuildIDs := syncer.getUniqueGuildIDs(guildIDs, commandMap)
	unregisterTargets := []unregisterTarget{}
	for _, guildID := range uniqueGuildIDs {
		log.Printf("\t- Guild: %s\n", guildText(guildID))
		commands, err := syncer.client.ListApplicationCommands(guildID)
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

func (syncer *syncer) registerCommands(commandMap SlashCommandMap) []error {
	errs := []error{}
	log.Println("Registering new commands...")
	for _, command := range commandMap {
		for _, guildID := range command.GuildIDs {
			log.Printf("\t- Guild: %s, Command: %s\n", guildText(guildID), command.Name)
			err := syncer.client.CreateApplicationCommand(guildID, command.AppCommand)
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

func (syncer *syncer) unregisterCommands(unregisterTargets []unregisterTarget) []error {
	errs := []error{}
	log.Println("Unregistering outdated commands...")
	for _, target := range unregisterTargets {
		log.Printf("\t- Guild: %s, Command: %s\n", guildText(target.guildID), target.name)
		err := syncer.client.DeleteApplicationCommand(target.guildID, target.commandID)
		if err != nil {
			log.Printf("\t\t- ERROR: %s\n", err.Error())
			errs = append(errs, err)
		} else {
			log.Printf("\t\t- SUCCESS")
		}
	}
	return errs
}

func (syncer *syncer) getUniqueGuildIDs(guildIDs []string, commands SlashCommandMap) []string {
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
