package model_test

import (
	"testing"

	"github.com/aimotrens/impulsar/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_Job_SetDefaults(t *testing.T) {
	j := &model.Job{}
	j.SetDefaults()

	assert.Equal(t, ".", j.WorkDir)
	assert.NotNil(t, j.Shell)
	assert.NotNil(t, j.Variables)
}

func Test_Job_Overwrite_OK(t *testing.T) {
	j := &model.Job{
		Name: "job",
		Shell: &model.Shell{
			Type: "bash",
		},
		Variables: model.VariableMap{
			"key": "value",
		},
	}

	overwrite := &model.Job{
		Shell: &model.Shell{
			Type: "powershell",
		},
		Variables: model.VariableMap{
			"key": "value2",
		},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, "job", j.Name)
	assert.Equal(t, "powershell", j.Shell.Type)
	assert.Equal(t, "value2", j.Variables["key"])
}

func Test_Job_Overwrite_Name_NOK(t *testing.T) {
	j := &model.Job{
		Name: "job",
	}

	overwrite := &model.Job{
		Name: "job2",
	}

	err := j.Overwrite(overwrite)
	assert.NotNil(t, err)
}

func Test_Job_Overwrite_Conditional_NOK(t *testing.T) {
	j := &model.Job{
		Conditional: []*model.Conditional{
			{},
		},
	}

	overwrite := &model.Job{
		Conditional: []*model.Conditional{
			{},
		},
	}

	err := j.Overwrite(overwrite)
	assert.NotNil(t, err)
}

func Test_Job_Overwrite_Shell_OK(t *testing.T) {
	j := &model.Job{
		Shell: &model.Shell{
			Type: "bash",
		},
	}

	overwrite := &model.Job{
		Shell: &model.Shell{
			Type: "powershell",
		},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, "powershell", j.Shell.Type)
}

func Test_Job_Overwrite_If_OK(t *testing.T) {
	j := &model.Job{
		If: []model.VariableMap{
			{},
		},
	}

	overwrite := &model.Job{
		If: []model.VariableMap{
			{},
		},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
}

func Test_Job_Overwrite_AllowFail_OK(t *testing.T) {
	j := &model.Job{
		AllowFail: false,
	}

	overwrite := &model.Job{
		AllowFail: true,
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.True(t, j.AllowFail)
}

func Test_Job_Overwrite_WorkDir_OK(t *testing.T) {
	j := &model.Job{
		WorkDir: ".",
	}

	overwrite := &model.Job{
		WorkDir: "tmp",
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, "tmp", j.WorkDir)
}

func Test_Job_Overwrite_JobsPre_OK(t *testing.T) {
	j := &model.Job{
		JobsPre: []string{"job1"},
	}

	overwrite := &model.Job{
		JobsPre: []string{"job2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"job2"}, j.JobsPre)
}

func Test_Job_Overwrite_JobsPost_OK(t *testing.T) {
	j := &model.Job{
		JobsPost: []string{"job1"},
	}

	overwrite := &model.Job{
		JobsPost: []string{"job2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"job2"}, j.JobsPost)
}

func Test_Job_Overwrite_JobsFinalize_OK(t *testing.T) {
	j := &model.Job{
		JobsFinalize: []string{"job1"},
	}

	overwrite := &model.Job{
		JobsFinalize: []string{"job2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"job2"}, j.JobsFinalize)
}

func Test_Job_Overwrite_ScriptPre_OK(t *testing.T) {
	j := &model.Job{
		ScriptPre: []string{"echo 1"},
	}

	overwrite := &model.Job{
		ScriptPre: []string{"echo 2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"echo 2"}, j.ScriptPre)
}

func Test_Job_Overwrite_Script_OK(t *testing.T) {
	j := &model.Job{
		Script: []string{"echo 1"},
	}

	overwrite := &model.Job{
		Script: []string{"echo 2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"echo 2"}, j.Script)
}

func Test_Job_Overwrite_ScriptPost_OK(t *testing.T) {
	j := &model.Job{
		ScriptPost: []string{"echo 1"},
	}

	overwrite := &model.Job{
		ScriptPost: []string{"echo 2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"echo 2"}, j.ScriptPost)
}

func Test_Job_Overwrite_ScriptFinalize_OK(t *testing.T) {
	j := &model.Job{
		ScriptFinalize: []string{"echo 1"},
	}

	overwrite := &model.Job{
		ScriptFinalize: []string{"echo 2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"echo 2"}, j.ScriptFinalize)
}

func Test_Job_Overwrite_Variables_OK(t *testing.T) {
	j := &model.Job{
		Variables: model.VariableMap{
			"key": "value",
		},
	}

	overwrite := &model.Job{
		Variables: model.VariableMap{
			"key": "value2",
		},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, "value2", j.Variables["key"])
}

func Test_Job_Overwrite_VariablesExcluded_OK(t *testing.T) {
	j := &model.Job{
		VariablesExcluded: []string{"key"},
	}

	overwrite := &model.Job{
		VariablesExcluded: []string{"key2"},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
	assert.Equal(t, []string{"key2"}, j.VariablesExcluded)
}

func Test_Job_Overwrite_Foreach_OK(t *testing.T) {
	j := &model.Job{
		Foreach: []model.VariableMap{
			{},
		},
	}

	overwrite := &model.Job{
		Foreach: []model.VariableMap{
			{},
		},
	}

	err := j.Overwrite(overwrite)
	assert.Nil(t, err)
}

func Test_Job_Overwrite_Arguments_NOK(t *testing.T) {
	j := &model.Job{
		Arguments: model.ArgumentMap{},
	}

	overwrite := &model.Job{
		Arguments: model.ArgumentMap{},
	}

	err := j.Overwrite(overwrite)
	assert.NotNil(t, err)
}

func Test_Job_Unmarshal_OK(t *testing.T) {
	yamlString := `
shell:
  type: bash
if:
  - key: value
allowFail: true
workDir: tmp
jobs:pre:
  - job1
jobs:post:
  - job2
jobs:finalize:
  - job3
script:pre:
  - echo 1
script:
  - echo 2
script:post:
  - echo 3
script:finalize:
  - echo 4
variables:
  key: value
foreach:
  - key: value
arguments:
  arg1: description
`

	var j model.Job
	err := yaml.Unmarshal([]byte(yamlString), &j)
	assert.Nil(t, err)
	assert.Equal(t, "bash", j.Shell.Type)
	assert.Equal(t, "value", j.If[0]["key"])
	assert.True(t, j.AllowFail)
	assert.Equal(t, "tmp", j.WorkDir)
	assert.Equal(t, []string{"job1"}, j.JobsPre)
	assert.Equal(t, []string{"job2"}, j.JobsPost)
	assert.Equal(t, []string{"job3"}, j.JobsFinalize)
	assert.Equal(t, []string{"echo 1"}, j.ScriptPre)
	assert.Equal(t, []string{"echo 2"}, j.Script)
	assert.Equal(t, []string{"echo 3"}, j.ScriptPost)
	assert.Equal(t, []string{"echo 4"}, j.ScriptFinalize)
	assert.Equal(t, "value", j.Variables["key"])
	assert.Equal(t, "value", j.Foreach[0]["key"])
	assert.Equal(t, "description", j.Arguments["arg1"].Description)
	assert.Equal(t, "", j.Arguments["arg1"].Default)
}

func Test_Job_Unmarshal_JobsAndJobsPreSet_NOK(t *testing.T) {
	yamlString := `
jobs:
  - job1
jobs:pre:
  - job2
`

	var j model.Job
	err := yaml.Unmarshal([]byte(yamlString), &j)
	assert.NotNil(t, err)
}
