package cli

import (
	"flag"
	"fmt"

	"github.com/aimotrens/impulsar/pkg/tui"
)

type (
	// BuildInfoProvider is a function that provides impulsar build information
	BuildInfoProvider func() (compileDate string, impulsarVersion string)

	// flagLoader is a function that loads flags for a command
	flagLoader func(flagDefiner) *flag.FlagSet

	// flagDefiner is a function that defines flags for a command
	flagDefiner func(*flag.FlagSet)
)

// Dispatch returns the actual command function to run based on the arguments
func Dispatch(args []string, buildInfo BuildInfoProvider) func() {
	// if we have no args, we can't do anything
	if len(args) < 1 {
		return nil
	}

	fl := flagLoader(func(setFlags flagDefiner) *flag.FlagSet {
		flags := flag.NewFlagSet(args[0], flag.ExitOnError)
		setFlags(flags)
		flags.Parse(args[1:])
		return flags
	})

	// command switch
	switch args[0] {
	case "run":
		return func() { run(fl, buildInfo) }
	case "gen":
		return func() { genCompletionScript(fl) }
	case "show-jobs":
		return func() { showJobs(fl) }
	case "version":
		return func() { version(fl, buildInfo) }
	default:
		return func() {
			// prepend "run" to args
			// a bit ugly, but the new args variable is captured by the closure (fl) above from this switch
			args = append([]string{"run"}, args...)
			run(fl, buildInfo)
		}
	}
}

// noFlagsUsage returns a function that prints a message that a command has no flags
func noFlagsUsage(cmdName string) func() {
	return func() { fmt.Println("command " + tui.Blue(tui.Italic(cmdName)) + " has no flags") }
}

// ColorizeCmd returns a string with the command name blue colored and italicized
func ColorizeCmd(cmdName string) string {
	return tui.Blue(tui.Italic(cmdName))
}
