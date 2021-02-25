// Package disgoslash provides an HTTP handler which can be used as a serverless endpoint for a Discord Application
// and a syncer implementation for keeping your slash commands registered with Discord up to date.
//
// Handler
//  func Handler(w http.ResponseWriter, r *http.Request) {
//  	creds := &discord.Credentials{} // Use your Discord application & bot credentials
//  	commands := SlashCommandMap{}   // Use your slash commands
//  	handler := &Handler{
//  		SlashCommandMap: commands,
//  		Creds:           creds,
//  	}
//  	handler.Handle(w, r)
//  }
//
// Syncer
//  creds := &discord.Credentials{}  // Use your Discord application & bot credentials
//  commands := SlashCommandMap{}    // Use your slash commands
//  guildIDs := []string{}           // List your guild (server) IDs
//  syncer := &Syncer{
//  	SlashCommandMap: commands,
//  	GuildIDs: guildIDs,
//  	Creds: creds,
//  }
//  syncer.Sync()
package disgoslash
