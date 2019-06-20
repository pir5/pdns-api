package pdns_api

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type GlobalFlags struct {
	ConfPath *string
	PidPath  *string
	LogPath  *string
}

// A Command is an implementation of a pdns-api command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmdFlags *GlobalFlags, args []string) error

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'pdns-api help' output.
	Short string

	// Long is the long message shown in the 'pdns-api help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}
