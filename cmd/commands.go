package cmd

import (
	"flag"
	"fmt"

	"github.com/aimotrens/impulsar/cout"
)

type flagLoader func(flagDefiner) *flag.FlagSet
type flagDefiner func(*flag.FlagSet)

func Dispatch(args []string) func() {
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
	case "version":
		return func() { version(fl) }
	case "run":
		return func() { run(fl) }
	case "gen":
		return func() { genCompletionScript(fl) }
	case "show-jobs":
		return func() { showJobs(fl) }
	default:
		return func() {
			args = append([]string{"run"}, args...)
			run(fl)
		}
	}
}

func noFlagsUsage(cmdName string) func() {
	return func() { fmt.Println("command " + cout.Blue(cout.Italic(cmdName)) + " has no flags") }
}

func ColorizeCmd(cmdName string) string {
	return cout.Blue(cout.Italic(cmdName))
}
