package model

import (
	"os/exec"
	"runtime"
)

const (
	SHELL_TYPE_POWERSHELL = "powershell"
	SHELL_TYPE_PWSH       = "pwsh"
	SHELL_TYPE_BASH       = "bash"
	SHELL_TYPE_DOCKER     = "docker"
	SHELL_TYPE_SSH        = "ssh"
	SHELL_TYPE_CUSTOM     = "custom"
)

type Shell struct {
	Type        string   `yaml:"type"`
	Image       string   `yaml:"image"`
	UidGid      string   `yaml:"uidGid"`
	BootCommand []string `yaml:"bootCommand"`
	Server      string   `yaml:"server"`
}

func (e *Shell) SetDefaults() {
	if e.Type == "" {
		switch runtime.GOOS {
		case "windows":
			e.Type = "powershell"

			// if pwsh exists, prefer it over powershell
			if _, err := exec.LookPath("pwsh"); err == nil {
				e.Type = "pwsh"
			}
		case "linux":
			e.Type = "bash"
		default:
			e.Type = "bash"
		}
	}

	if e.BootCommand == nil {
		switch e.Type {
		case "powershell":
			e.BootCommand = []string{"powershell", "-noprofile", "-Command"}
		case "pwsh":
			e.BootCommand = []string{"pwsh", "-noprofile", "-Command"}
		case "bash":
			e.BootCommand = []string{"bash", "-c"}

			if runtime.GOOS == "windows" {
				e.BootCommand = append([]string{"wsl", "--exec"}, e.BootCommand...)
			}
		default:
			e.BootCommand = []string{"bash", "-c"}
		}
	}
}
