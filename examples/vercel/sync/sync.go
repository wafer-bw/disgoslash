package main

import (
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/examples/vercel/api"
)

func main() {
	syncer := &disgoslash.Syncer{
		Creds:           api.Credentials,
		SlashCommandMap: api.SlashCommandMap,
		GuildIDs:        api.GuildIDs,
	}
	syncer.Sync()
}
