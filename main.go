package main

import (
	"fmt"
	"os"

	"github.com/aimotrens/impulsar/cmd"
)

func main() {
	command := cmd.Dispatch(os.Args[1:])
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
