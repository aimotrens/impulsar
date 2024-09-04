package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aimotrens/impulsar/cout"
	"github.com/aimotrens/impulsar/model"
)

// Collects all arguments for a job recursively
func (e *Engine) CollectArgs(job string) {
	if j, ok := e.GetJob(job); ok {
		e.readArgsIntoJobVars(j)

		for _, pre := range j.JobsPre {
			e.CollectArgs(pre)
		}

		for _, post := range j.JobsPost {
			e.CollectArgs(post)
		}

		return
	}

	fmt.Printf("Job %s not found\n", job)
}

// Processes all arguments for a job
// If it exists as env var, it will be used
// If it does not exist, it will be asked
func (e *Engine) readArgsIntoJobVars(j *model.Job) {
	for name, argDef := range j.Arguments {
		if _, ok := j.Variables[name]; ok {
			continue
		}

		if val, ok := e.Variables[name]; ok {
			j.Variables[name] = val
			continue
		}

		if argDef.Allowed == nil {
			fmt.Printf("%s %s (%s) [%s]: ",
				cout.Green(cout.Bold("["+j.Name+"]")),
				cout.Gray(name),
				argDef.Description,
				cout.LightYellow(argDef.Default),
			)
		} else {
			fmt.Printf("%s %s (%s) [%s] {%s}: ",
				cout.Green(cout.Bold("["+j.Name+"]")),
				cout.Gray(name),
				argDef.Description,
				cout.LightYellow(argDef.Default),
				cout.FormatByText(strings.Join(argDef.Allowed, ", "), argDef.Default, cout.Italic),
			)
		}

		var value string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			value = scanner.Text()
		}

		if value == "" {
			value = argDef.Default
		}

		if argDef.Allowed != nil {
			isAllowedValue := false

			for _, a := range argDef.Allowed {
				if a == value {
					isAllowedValue = true
					break
				}
			}

			if !isAllowedValue {
				fmt.Printf(cout.Red("Value %s not allowed\n"), value)
				os.Exit(1)
			}
		}

		j.Variables[name] = value
	}
}
