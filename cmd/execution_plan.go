package cmd

import (
	"fmt"
	"strings"

	"github.com/aimotrens/impulsar/cout"
	"github.com/aimotrens/impulsar/engine"
	"github.com/aimotrens/impulsar/model"
)

func buildExecutionPlan(e *engine.Engine, requestedJobs []string) (planString string, ok bool) {
	jobPlanSeparator := fmt.Sprintf("%s\n", strings.Repeat("─", 60))
	canExecute := true

	plan := &strings.Builder{}
	plan.WriteString(jobPlanSeparator)

	for _, jobName := range requestedJobs {
		jobPlan, jobOk := getJobPlanRecursive(e, jobName, "", 0, make(map[string]bool))
		plan.WriteString(jobPlan)
		plan.WriteString(jobPlanSeparator)
		canExecute = canExecute && jobOk
	}

	return plan.String(), canExecute
}

func getJobPlanRecursive(e *engine.Engine, jobName string, callSuffix string, level int, visited map[string]bool) (plan string, ok bool) {
	ok = true

	if _, ok := visited[jobName]; ok {
		return indent(level) + "- " + jobName + callSuffix + cout.Red(" (err circular)") + "\n", false
	}

	visited[jobName] = true

	// Shortcuts
	g := cout.Green
	b := cout.Bold
	gb := cout.Multi(g, b)
	gi := cout.Multi(g, cout.Italic)
	// ---

	job, found := e.GetJob(jobName)
	if !found {
		return indent(level) + "- " + jobName + callSuffix + cout.Red(" (err not found)") + "\n", false
	}

	jobRuns := genJobRuns(job)
	jobRunsSuffix := ""
	if len(jobRuns) > 0 {
		jobRunsSuffix = cout.Blue(fmt.Sprintf(" with %s", strings.Join(jobRuns, ", ")))
	}

	for _, subJob := range job.JobsPre {
		subPlan, subOk := getJobPlanRecursive(e, subJob, g(" via ")+gb(jobName)+g("/jobs:")+gi("pre"), level+1, visited)
		plan += subPlan
		ok = ok && subOk
	}

	plan += indent(level) + "- " + job.Name + callSuffix + jobRunsSuffix + "\n"

	for _, subJob := range job.JobsPost {
		subPlan, subOk := getJobPlanRecursive(e, subJob, g(" via ")+gb(jobName)+g("/jobs:")+gi("post"), level+1, visited)
		plan += subPlan
		ok = ok && subOk
	}

	return
}

func indent(l int) string {
	return strings.Repeat("  ", l)
}

func genJobRuns(j *model.Job) []string {
	jobRuns := []string{}
	if len(j.Foreach) > 0 {
		jobRuns = append(jobRuns, fmt.Sprintf("foreach(%d)", len(j.Foreach)))
	}

	if len(j.ScriptPre) > 0 {
		jobRuns = append(jobRuns, fmt.Sprintf("pre(%d)", len(j.ScriptPre)))
	}

	if len(j.Script) > 0 {
		jobRuns = append(jobRuns, fmt.Sprintf("script(%d)", len(j.Script)))
	}

	if len(j.ScriptPost) > 0 {
		jobRuns = append(jobRuns, fmt.Sprintf("post(%d)", len(j.ScriptPre)))
	}

	if len(j.ScriptFinalize) > 0 {
		jobRuns = append(jobRuns, fmt.Sprintf("finalize(%d)", len(j.ScriptPre)))
	}

	return jobRuns
}
