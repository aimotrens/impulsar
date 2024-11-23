package model_test

import (
	"testing"

	"github.com/aimotrens/impulsar/internal/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_Argument_Unmarshal_String_OK(t *testing.T) {
	yamlString := `
arg1: description
`
	var a model.ArgumentMap
	err := yaml.Unmarshal([]byte(yamlString), &a)
	assert.Nil(t, err)
	assert.Equal(t, "description", a["arg1"].Description)
	assert.Equal(t, "", a["arg1"].Default)
}

func Test_Argument_Unmarshal_Argument_OK(t *testing.T) {
	yamlString := `
arg1:
  description: "description"
  default: "default"
`
	var a model.ArgumentMap
	err := yaml.Unmarshal([]byte(yamlString), &a)
	assert.Nil(t, err)
	assert.Equal(t, "description", a["arg1"].Description)
	assert.Equal(t, "default", a["arg1"].Default)
}

func Test_Argument_Unmarshal_Argument2_OK(t *testing.T) {
	yamlString := `
arg1:
  description: "description"
  default: "default"
  allowed: ["default", "test"]
`
	var a model.ArgumentMap
	err := yaml.Unmarshal([]byte(yamlString), &a)
	assert.Nil(t, err)
	assert.Equal(t, "description", a["arg1"].Description)
	assert.Equal(t, "default", a["arg1"].Default)
	assert.Equal(t, []string{"default", "test"}, a["arg1"].Allowed)
}

func Test_Argument_Unmarshal_Argument_DefaultNotInAllowedValues_NOK(t *testing.T) {
	yamlString := `
arg1:
  description: "description"
  default: "default"
  allowed: ["allowed", "test"]
`
	var a model.ArgumentMap
	err := yaml.Unmarshal([]byte(yamlString), &a)
	assert.NotNil(t, err)
}
