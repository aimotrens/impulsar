package main

import (
	"fmt"
	"os"

	"github.com/aimotrens/impulsar/cmd"
)

var (
	compileDate     string = "unknown"
	impulsarVersion string = "vX.X.X"
)

func main() {
	command := cmd.Dispatch(os.Args[1:], buildInfo)
	if command == nil {
		fmt.Println("Usage: impulsar " + cmd.ColorizeCmd("command") + " [flags]")
		fmt.Println("       impulsar [<flags of " + cmd.ColorizeCmd("run") + " cmd>] <job>...")

		fmt.Println("\nCommands:")
		fmt.Println("  " + cmd.ColorizeCmd("run") + " - Run impulsar jobs (can be omitted for compatibility)")
		fmt.Println("  " + cmd.ColorizeCmd("gen") + " - Generate completion script")
		fmt.Println("  " + cmd.ColorizeCmd("show-jobs") + " - Show impulsar jobs")
		fmt.Println("  " + cmd.ColorizeCmd("version") + " - Show impulsar version")

		os.Exit(1)
	}

	command()
}

func buildInfo() (string, string) {
	return compileDate, impulsarVersion
}
