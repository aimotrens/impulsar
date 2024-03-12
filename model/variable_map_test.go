package model_test

import (
	"testing"

	"github.com/aimotrens/impulsar/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_VariableMap_Unmarshal_OK(t *testing.T) {
	yamlString := `
key1: value1
key2: value2
`

	var v model.VariableMap
	err := yaml.Unmarshal([]byte(yamlString), &v)
	assert.Nil(t, err)
	assert.Equal(t, "value1", v["key1"])
	assert.Equal(t, "value2", v["key2"])
}
