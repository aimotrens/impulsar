package main

import (
	"fmt"
	"os"

	"github.com/aimotrens/impulsar/internal/cli"
)

var (
	compileDate     string = "unknown"
	impulsarVersion string = "vX.X.X"
)

func main() {
	command := cli.Dispatch(os.Args[1:], buildInfo)
	if command == nil {
		fmt.Println("Usage: impulsar " + cli.ColorizeCmd("command") + " [flags]")
		fmt.Println("       impulsar [<flags of " + cli.ColorizeCmd("run") + " cmd>] <job>...")

		fmt.Println("\nCommands:")
		fmt.Println("  " + cli.ColorizeCmd("run") + " - Run impulsar jobs (can be omitted for compatibility)")
		fmt.Println("  " + cli.ColorizeCmd("gen") + " - Generate completion script")
		fmt.Println("  " + cli.ColorizeCmd("show-jobs") + " - Show impulsar jobs")
		fmt.Println("  " + cli.ColorizeCmd("version") + " - Show impulsar version")

		os.Exit(0)
	}

	command()
}

func buildInfo() (string, string) {
	return compileDate, impulsarVersion
}
