package model

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Name      string            `yaml:"-"`
	Shell     *Shell            `yaml:"shell"`
	AllowFail bool              `yaml:"allowFail"`
	WorkDir   string            `yaml:"workDir"`
	JobsPre   []string          `yaml:"jobs:pre"`
	JobsPost  []string          `yaml:"jobs:post"`
	Script    []string          `yaml:"script"`
	Variables map[string]string `yaml:"variables"`
}

func (j *Job) SetDefaults() {
	if j.Shell == nil {
		j.Shell = &Shell{}
	}
	j.Shell.SetDefaults()

	if j.WorkDir == "" {
		j.WorkDir, _ = os.Getwd()
	}
}

func (j *Job) UnmarshalYAML(v *yaml.Node) error {
	type OriginJob Job
	var tmpJob struct {
		OriginJob
		Jobs []string `yaml:"jobs"`
	}

	err := v.Decode(&tmpJob.OriginJob)
	if err != nil {
		return err
	}

	err = v.Decode(&tmpJob)
	if err != nil {
		return err
	}

	*j = Job(tmpJob.OriginJob)

	if tmpJob.Jobs != nil && tmpJob.JobsPre != nil {
		return errors.New("jobs and jobs:pre are mutually exclusive")
	}

	*j = Job(tmpJob.OriginJob)

	if tmpJob.Jobs != nil {
		j.JobsPre = tmpJob.Jobs
	}

	return nil
}
