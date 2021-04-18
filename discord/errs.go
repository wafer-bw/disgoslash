package discord

import "errors"

// ErrNotSubcommand is returned when an ApplicationCommandInteractionOptions
// slice length is not 1 but GetSubcommand was called on it or when
// the only option in the slice was nil.
//
// A subcommand will be the only option in an ApplicationCommandInteractionOptions
// slice if the corresponding ApplicationCommandOption's Type is
// ApplicationCommandOptionTypeSubCommand.
var ErrNotSubcommand error = errors.New("option's slice is not the correct length (1) for a subcommand")

// ErrNotSubcommandGroup is returned when an ApplicationCommandInteractionOptions
// slice length is not 1 but GetSubcommandGroup was called on it
// the only option in the slice was nil.
//
// A subcommand group will be the only option in an ApplicationCommandInteractionOptions
// slice if the corresponding ApplicationCommandOption's Type is
// ApplicationCommandOptionTypeSubCommandGroup.
var ErrNotSubcommandGroup error = errors.New("option's slice is not the correct length (1) for a subcommand")
