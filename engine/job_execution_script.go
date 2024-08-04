package engine

import (
	"fmt"
	"os"
	"strings"

	"github.com/aimotrens/impulsar/cout"
	"github.com/aimotrens/impulsar/model"
)

func (e *Engine) runScriptBlock(j *model.Job, scriptBlock []string, suffix string) error {
	for _, script := range scriptBlock {
		if script == "STOP" {
			fmt.Printf("Job %s failed, due to STOP command\n", j.Name)
			os.Exit(1)
		}

		scriptPrint := ""
		for _, s := range strings.Split(script, "\n") {
			scriptPrint += strings.Trim(s, " ") + "; "
		}

		scriptPrint = strings.Trim(scriptPrint, " ;")

		if len(scriptPrint) > 81 {
			scriptPrint = scriptPrint[0:50] + " . . . . . " + scriptPrint[len(scriptPrint)-20:]
		}

		fmt.Print(
			cout.Green(cout.Bold("["+j.Name+"] ")) +
				"(" + scriptPrint + ") " +
				cout.Blue(suffix) + "\n",
		)

		if err := e.execCommand(j, script); err != nil {
			return err
		}
	}

	return nil
}
