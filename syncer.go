package disgoslash

import (
	"log"

	"github.com/wafer-bw/disgoslash/discord"
)

// Syncer is used to automatically update your Discord application's slash commands
type Syncer struct {
	Creds           *discord.Credentials
	SlashCommandMap SlashCommandMap
	GuildIDs        []string
	client          ClientInterface
}

type unregisterTarget struct {
	guildID   string
	commandID string
	name      string
}

// Sync your Discord application's slash commands...
//
// registers new slash commands
// unregisters old slash commands
// reregisters existing slash commands
func (syncer *Syncer) Sync() []error {
	if syncer.client == nil {
		syncer.client = NewClient(syncer.Creds)
	}
	allErrs := []error{}
	unregisterTargets, errs := syncer.getCommandsToUnregister()
	if len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	allErrs = append(allErrs, syncer.unregisterCommands(unregisterTargets)...)
	allErrs = append(allErrs, syncer.registerCommands(syncer.SlashCommandMap)...)
	return allErrs
}

func (syncer *Syncer) getCommandsToUnregister() ([]unregisterTarget, []error) {
	errs := []error{}
	log.Println("Collecting outdated commands...")
	uniqueGuildIDs := syncer.getUniqueGuildIDs(syncer.GuildIDs, syncer.SlashCommandMap)
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

func (syncer *Syncer) unregisterCommands(unregisterTargets []unregisterTarget) []error {
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

func (syncer *Syncer) registerCommands(commandMap SlashCommandMap) []error {
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

func (syncer *Syncer) getUniqueGuildIDs(guildIDs []string, commands SlashCommandMap) []string {
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
