package util

import (
	"os/exec"
)

// CommandExists checks if a command exists on the system
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
