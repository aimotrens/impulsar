package model_test

import (
	"testing"

	"github.com/aimotrens/impulsar/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_Conditional_Unmarshal_OK(t *testing.T) {
	yamlString := `
if:
  - key: value
overwrite:
  script:
    - echo
`
	var c model.Conditional
	err := yaml.Unmarshal([]byte(yamlString), &c)
	assert.Nil(t, err)
	assert.Equal(t, "value", c.If[0]["key"])
	assert.Equal(t, "echo", c.Overwrite.Script[0])
}
