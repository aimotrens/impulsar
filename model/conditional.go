package model

type Conditional struct {
	If        []string `yaml:"if"`
	Overwrite *Job     `yaml:"overwrite"`
}
