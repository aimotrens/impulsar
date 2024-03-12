package model_test

import (
	"testing"

	"github.com/aimotrens/impulsar/model"
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
