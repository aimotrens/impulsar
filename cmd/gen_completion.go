package cmd

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/aimotrens/impulsar/cout"
)

var (
	//go:embed shell_completion_scripts/bash.sh
	bashCompletion string

	//go:embed shell_completion_scripts/zsh.sh
	zshCompletion string

	//go:embed shell_completion_scripts/powershell.ps1
	powershellCompletion string
)

func genCompletionScript(fl flagLoader) {
	var outputFile string

	flags := fl(func(fs *flag.FlagSet) {
		fs.Usage = func() {
			fmt.Println("Usage: impulsar " + ColorizeCmd("gen") + " [-o <output file>] bash|zsh|powershell")
			fs.PrintDefaults()
		}
		fs.StringVar(&outputFile, "o", "", "output file")
	})

	if flags.NArg() == 0 {
		flags.Usage()
		os.Exit(1)
	}

	output := os.Stdout
	if outputFile != "" {
		var err error
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Println(cout.Red("Error creating file:" + outputFile))
			os.Exit(1)
		}
	}

	switch flags.Arg(0) {
	case "bash":
		fmt.Fprint(output, bashCompletion)
	case "zsh":
		fmt.Fprint(output, zshCompletion)
	case "powershell":
		fmt.Fprint(output, powershellCompletion)
	default:
		fmt.Println(cout.Red("Unknown shell: " + flags.Arg(0)))
		os.Exit(1)
	}
}
