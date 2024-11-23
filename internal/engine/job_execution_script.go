package engine

import (
	"fmt"
	"os"
	"strings"

	"github.com/aimotrens/impulsar/internal/model"
	"github.com/aimotrens/impulsar/pkg/tui"
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
			tui.Green(tui.Bold("["+j.Name+"] ")) +
				"(" + scriptPrint + ") " +
				tui.Blue(suffix) + "\n",
		)

		if err := e.execCommand(j, script); err != nil {
			return err
		}
	}

	return nil
}
