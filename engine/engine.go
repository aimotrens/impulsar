package engine

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/aimotrens/impulsar/cout"

	"github.com/aimotrens/impulsar/model"
)

type (
	Engine struct {
		jobList   map[string]*model.Job
		shellMap  map[string]Shell
		Variables model.VariableMap
	}

	ExecutorConstructor func(*Engine) Shell
)

var executorMap = make(map[string]ExecutorConstructor)

func RegisterExecutor(name string, constructor ExecutorConstructor) {
	executorMap[name] = constructor
}

func New(jobList map[string]*model.Job, additionalEnvVars model.VariableMap) *Engine {
	envVars := make(model.VariableMap)

	// Aggregate environment and additional variables
	{
		for _, v := range os.Environ() {
			var kv = strings.Split(v, "=")
			envVars[kv[0]] = kv[1]
		}

		for key, value := range additionalEnvVars {
			envVars[key] = value
		}
	}

	e := &Engine{
		jobList:   jobList,
		shellMap:  make(map[string]Shell),
		Variables: envVars,
	}

	for name, constructor := range executorMap {
		e.shellMap[name] = constructor(e)
	}

	return e
}

// Collects all arguments for a job recursively
func (e *Engine) CollectArgs(job string) {
	if j, ok := e.jobList[job]; ok {
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

func (e *Engine) RunJob(job string) {
	if scheduledJob, ok := e.jobList[job]; ok {
		e.executeJob(scheduledJob)
		return
	}

	fmt.Printf("Job %s not found\n", job)
}

func (e *Engine) executeJob(j *model.Job) {
	e.evaluateConditionalField(j)

	runScriptBlock := func(scriptBlock []string, suffix string) error {
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

	exitWithFinalizer := func() {
		_ = runScriptBlock(j.ScriptFinalize, "via script:finalize")

		for _, finalize := range j.JobsFinalize {
			e.RunJob(finalize)
		}

		os.Exit(1)
	}

	if e.evaluateIfCondition(j) {
		for _, pre := range j.JobsPre {
			e.RunJob(pre)
		}

		if j.ScriptPre != nil {
			if err := runScriptBlock(j.ScriptPre, "via script:pre"); err != nil {
				exitWithFinalizer()
			}
		}

		if j.Foreach != nil {
			for _, f := range j.Foreach {
				for k, v := range f {
					e.Variables[k] = v
				}

				if err := runScriptBlock(j.Script, "via foreach"); err != nil {
					exitWithFinalizer()
				}
			}
		} else {
			if err := runScriptBlock(j.Script, ""); err != nil {
				exitWithFinalizer()
			}
		}

		if j.ScriptPost != nil {
			if err := runScriptBlock(j.ScriptPost, "via script:post"); err != nil {
				exitWithFinalizer()
			}
		}

		if j.ScriptFinalize != nil {
			_ = runScriptBlock(j.ScriptFinalize, "view script:finalize")
		}

		for _, post := range j.JobsPost {
			e.RunJob(post)
		}

		for _, fin := range j.JobsFinalize {
			e.RunJob(fin)
		}
	} else {
		fmt.Printf("Job %s skipped, no condition matched\n", j.Name)
	}
}

func (e *Engine) execCommand(j *model.Job, script string) error {
	if shell, ok := e.shellMap[j.Shell.Type]; !ok {
		fmt.Printf("Shell type %s not supported\n", j.Shell.Type)
		os.Exit(1)
		return nil
	} else {
		return shell.Execute(j, script)
	}
}

func (e *Engine) LookupVarFunc(j *model.Job) func(string) string {
	return func(s string) string {
		if v, ok := j.Variables[s]; ok {
			return v
		}

		if v, ok := e.Variables[s]; ok {
			return v
		}

		return ""
	}
}

func (e *Engine) evaluateIfCondition(j *model.Job) bool {
	if j.If == nil {
		return true
	}

	envVars := e.aggregateEnvVars(j)

	for _, varSet := range j.If {
		var success = true

		for k, v := range varSet {
			success = success && (envVars[k] == v)
		}

		if success {
			return true
		}
	}

	return false
}

func (e *Engine) evaluateConditionalField(j *model.Job) {
	if j.Conditional == nil {
		return
	}

	var envVars = e.aggregateEnvVars(j)

	for _, v := range j.Conditional {
		for _, varSet := range v.If {
			var success = true

			for k, v := range varSet {
				success = success && (envVars[k] == v)
			}

			if success {
				if v.Overwrite != nil {
					j.Overwrite(v.Overwrite)
				}

				return
			}
		}
	}
}

func (e *Engine) aggregateEnvVars(j *model.Job) model.VariableMap {
	var envVars = make(model.VariableMap)

	for key, value := range e.Variables {
		envVars[key] = value
	}

	for key, value := range j.Variables {
		envVars[key] = value
	}

	envVars["OS"] = runtime.GOOS
	envVars["ARCH"] = runtime.GOARCH

	return envVars
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
