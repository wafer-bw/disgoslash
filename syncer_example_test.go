package disgoslash_test

import (
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
)

func ExampleSyncer_Sync() {
	guildIDs := []string{"YOUR_GUILD_(SERVER)_ID", "ANOTHER_GUILD_ID"}
	creds := &discord.Credentials{
		PublicKey: "YOUR_DISCORD_APPLICATION_PUBLIC_KEY",
		ClientID:  "YOUR_DISCORD_APPLICATION_CLIENT_ID",
		Token:     "YOUR_DISCORD_BOT_TOKEN",
	}

	syncer := &disgoslash.Syncer{
		SlashCommandMap: disgoslash.NewSlashCommandMap(disgoslash.SlashCommand{}),
		GuildIDs:        guildIDs,
		Creds:           creds,
	}
	syncer.Sync()
}
