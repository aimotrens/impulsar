package model

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ArgumentMap map[string]*ArgumentDefinition

type ArgumentDefinition struct {
	Description string   `yaml:"description"`
	Default     string   `yaml:"default"`
	Allowed     []string `yaml:"allowed"`
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
		argDef := &value.ArgumentDefinition
		(*a)[key] = argDef

		if argDef.Default != "" && len(argDef.Allowed) > 0 {
			defaultIsAllowedValue := false
			for _, allowed := range argDef.Allowed {
				if argDef.Default == allowed {
					defaultIsAllowedValue = true
					break
				}
			}

			if !defaultIsAllowedValue {
				return fmt.Errorf("default value is not in allowed values (line %v)", v.Line)
			}
		}
	}

	return nil
}

func (d *decoderMultiType) UnmarshalYAML(v *yaml.Node) error {
	err := v.Decode(&d.string)
	if err == nil {
		d.Description = d.string
		return nil
	}

	err = v.Decode(&d.ArgumentDefinition)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to unmarshal argument")
}
