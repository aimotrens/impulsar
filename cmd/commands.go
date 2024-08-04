package cmd

import (
	"flag"
	"fmt"

	"github.com/aimotrens/impulsar/cout"
)

type BuildInfoProvider func() (compileDate string, impulsarVersion string)
type flagLoader func(flagDefiner) *flag.FlagSet
type flagDefiner func(*flag.FlagSet)

func Dispatch(args []string, buildInfo BuildInfoProvider) func() {
	if len(args) < 1 {
		return nil
	}

	fl := flagLoader(func(setFlags flagDefiner) *flag.FlagSet {
		flags := flag.NewFlagSet(args[0], flag.ExitOnError)
		setFlags(flags)
		flags.Parse(args[1:])
		return flags
	})

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

func noFlagsUsage(cmdName string) func() {
	return func() { fmt.Println("command " + cout.Blue(cout.Italic(cmdName)) + " has no flags") }
}

func ColorizeCmd(cmdName string) string {
	return cout.Blue(cout.Italic(cmdName))
}
