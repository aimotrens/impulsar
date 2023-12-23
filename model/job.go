package model

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Name        string            `yaml:"-"`
	Shell       *Shell            `yaml:"shell"`
	If          []string          `yaml:"if"`
	Conditional []*Conditional    `yaml:"conditional"`
	AllowFail   bool              `yaml:"allowFail"`
	WorkDir     string            `yaml:"workDir"`
	JobsPre     []string          `yaml:"jobs:pre"`
	JobsPost    []string          `yaml:"jobs:post"`
	Script      []string          `yaml:"script"`
	Variables   map[string]string `yaml:"variables"`
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

func (j *Job) Overwrite(overwrite *Job) error {
	if overwrite.Name != "" {
		return errors.New("cannot overwrite job name")
	}

	if overwrite.Conditional != nil {
		return errors.New("cannot overwrite job conditional")
	}

	if overwrite.Shell != nil {
		j.Shell = overwrite.Shell
	}

	if overwrite.If != nil {
		j.If = overwrite.If
	}

	if overwrite.AllowFail {
		j.AllowFail = overwrite.AllowFail
	}

	if overwrite.WorkDir != "" {
		j.WorkDir = overwrite.WorkDir
	}

	if overwrite.JobsPre != nil {
		j.JobsPre = overwrite.JobsPre
	}

	if overwrite.JobsPost != nil {
		j.JobsPost = overwrite.JobsPost
	}

	if overwrite.Script != nil {
		j.Script = overwrite.Script
	}

	if overwrite.Variables != nil {
		j.Variables = overwrite.Variables
	}

	overwrite.SetDefaults()

	return nil
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
