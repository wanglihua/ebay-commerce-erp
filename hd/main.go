package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/revel/revel"
)

const version = "1.0.0"

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) int

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short template.HTML

	// Long is the long message shown in the 'go help <this-command>' output.
	Long template.HTML

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
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
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Long)))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var commands = []*Command{
	cmdRun,
	cmdRoute,
	cmdPackage,
	cmdLic,
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		os.Exit(cmdRun.Run(cmdRun, args))
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Run(cmd, args[1:])
			return
		}
	}

	revel.ERROR.Println("unknow command: " + args[0])
}
