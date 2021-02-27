package disgoslash_test

import (
	"github.com/wafer-bw/disgoslash"
	discord "github.com/wafer-bw/disgoslash/discord"
)

// A Syncer requires your Discord credentials, a map of slash commands
// created by disgoslash.NewSlashCommandMap() which accepts  slash commands
// created by disgoslash.NewSlashCommand(), and the guild (server) IDs that
// you have commands registered with already.
func ExampleSyncer_Sync() {
	guildIDs := []string{"YOUR_GUILD_(SERVER)_ID", "ANOTHER_GUILD_ID"}
	creds := &discord.Credentials{
		PublicKey: "YOUR_DISCORD_APPLICATION_PUBLIC_KEY",
		ClientID:  "YOUR_DISCORD_APPLICATION_CLIENT_ID",
		Token:     "YOUR_DISCORD_BOT_TOKEN",
	}

	syncer := &disgoslash.Syncer{
		SlashCommandMap: disgoslash.SlashCommandMap{}, // Replace this with your slash command map,
		GuildIDs:        guildIDs,
		Creds:           creds,
	}

	syncer.Sync()
}
