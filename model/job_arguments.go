package model

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type ArgumentMap map[string]*ArgumentDefinition

type ArgumentDefinition struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
}

type decoderMultiType struct {
	ArgumentDefinition
	string
}

func (a *ArgumentMap) UnmarshalYAML(v *yaml.Node) error {
	*a = make(ArgumentMap)

	tmp := make(map[string]*decoderMultiType)
	err := v.Decode(&tmp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal arguments")
	}

	for key, value := range tmp {
		(*a)[key] = &value.ArgumentDefinition
	}

	return nil
}

func (d *decoderMultiType) UnmarshalYAML(v *yaml.Node) error {
	err := v.Decode(&d.string)
	if err == nil {
		return nil
	}

	err = v.Decode(&d.ArgumentDefinition)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to unmarshal argument")
}
