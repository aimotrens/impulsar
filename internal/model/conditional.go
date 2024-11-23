package model

type Conditional struct {
	If        []VariableMap `yaml:"if"`
	Overwrite *Job          `yaml:"overwrite"`
}
