package model_test

import (
	"runtime"
	"testing"

	"github.com/aimotrens/impulsar/internal/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_Shell_SetDefaults(t *testing.T) {
	s := &model.Shell{}
	s.SetDefaults()

	if runtime.GOOS == "windows" {
		assert.Equal(t, "powershell", s.Type)
		assert.Equal(t, []string{"powershell", "-Command"}, s.BootCommand)
	} else if runtime.GOOS == "linux" {
		assert.Equal(t, "bash", s.Type)
		assert.Equal(t, []string{"bash", "-c"}, s.BootCommand)
	} else {
		assert.Equal(t, "bash", s.Type)
		assert.Equal(t, []string{"bash", "-c"}, s.BootCommand)
	}
}

func Test_Shell_Unmarshal_OK(t *testing.T) {
	yamlString := `
type: bash
image: image
uidGid: uidGid
bootCommand:
  - bootCommand
  - arg1
  - arg2
server: server
`

	var s model.Shell
	err := yaml.Unmarshal([]byte(yamlString), &s)
	assert.Nil(t, err)
	assert.Equal(t, "bash", s.Type)
	assert.Equal(t, "image", s.Image)
	assert.Equal(t, "uidGid", s.UidGid)
	assert.Equal(t, []string{"bootCommand", "arg1", "arg2"}, s.BootCommand)
	assert.Equal(t, "server", s.Server)
}
