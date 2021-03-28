package disgoslash

import (
	"log"

	"github.com/wafer-bw/disgoslash/discord"
)

// Syncer is used to automatically update slash commands on Discord guilds (servers).
type Syncer struct {
	Creds           *discord.Credentials
	SlashCommandMap SlashCommandMap
	GuildIDs        []string
	client          clientInterface
}

type unregisterTarget struct {
	guildID   string
	commandID string
	name      string
}

// Sync your Discord application's slash commands...
//
// Registers new commands, unregisters old commands, and reregisters existing commands.
//
// In order for a command to be registered
// to a guild (server), the bot will need to be granted
// the "appliations.commands" scope for that server.
//
// A global command will be registered to all servers
// the bot has been granted access to.
func (syncer *Syncer) Sync() []error {
	if syncer.client == nil {
		syncer.client = newClient(syncer.Creds)
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
		log.Printf("\tGuild: %s\n", guildText(guildID))
		commands, err := syncer.client.list(guildID)
		if err != nil {
			log.Printf("\t\terror: %s\n", err.Error())
			errs = append(errs, err)
		} else {
			log.Printf("\t\tsuccess")
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
		log.Printf("\tGuild: %s, Command: %s\n", guildText(target.guildID), target.name)
		err := syncer.client.delete(target.guildID, target.commandID)
		if err != nil {
			log.Printf("\t\terror: %s\n", err.Error())
			errs = append(errs, err)
		} else {
			log.Printf("\t\tsuccess")
		}
	}
	return errs
}

func (syncer *Syncer) registerCommands(commandMap SlashCommandMap) []error {
	errs := []error{}
	log.Println("Registering new commands...")
	for _, command := range commandMap {
		for _, guildID := range command.GuildIDs {
			log.Printf("\tGuild: %s, Command: %s\n", guildText(guildID), command.Name)
			err := syncer.client.create(guildID, command.ApplicationCommand)
			if err != nil {
				log.Printf("\t\terror: %s\n", err.Error())
				errs = append(errs, err)
			} else {
				log.Printf("\t\tsuccess")
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
