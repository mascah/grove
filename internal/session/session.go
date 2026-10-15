// Package session runs the harness, Claude Code, in the foreground of this
// terminal until it exits.
package session

import (
	"os"
	"os/exec"
)

// Start opens a new conversation in dir with prompt as its first message.
func Start(dir, prompt string) error { return attach(dir, prompt) }

// Resume reopens the most recent conversation in dir. Each workspace is its
// own directory, so that is the workspace's last session.
func Resume(dir string) error { return attach(dir, "--continue") }

func attach(dir string, args ...string) error {
	cmd := exec.Command("claude", args...)
	cmd.Dir = dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
