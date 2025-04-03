package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRootCmd(t *testing.T) {
	t.Run("Root command has expected subcommands", func(t *testing.T) {
		if rootCmd == nil {
			t.Fatal("rootCmd is nil")
		}

		// Verify basic properties
		if rootCmd.Use != "bt" {
			t.Errorf("Expected root command Use to be 'bt', got %s", rootCmd.Use)
		}

		if rootCmd.Version != Version {
			t.Errorf("Expected root command Version to be '%s', got %s", Version, rootCmd.Version)
		}

		// Verify expected subcommands exist
		expectedCmds := []string{"new", "list", "project"}
		for _, name := range expectedCmds {
			found := false
			for _, cmd := range rootCmd.Commands() {
				if cmd.Name() == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find subcommand '%s'", name)
			}
		}
	})
}

func TestListCmd(t *testing.T) {
	t.Run("List command output", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		listCmd.Run(listCmd, []string{})

		// Restore stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Basic check for expected content
		if len(output) == 0 {
			t.Error("Expected non-empty output")
		}

		expectedStrings := []string{
			"Available frameworks:",
			"---------------------",
		}

		for _, s := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(s)) {
				t.Errorf("Expected output to contain '%s'", s)
			}
		}
	})
}

func TestVersionTemplate(t *testing.T) {
	t.Run("Version template is set", func(t *testing.T) {
		vt := rootCmd.VersionTemplate()
		if vt == "" {
			t.Error("Expected non-empty version template")
		}
	})
}

func TestNewAndProjectCmdArgs(t *testing.T) {
	t.Run("New command requires 2 args", func(t *testing.T) {
		if newCmd.Args == nil {
			t.Fatal("newCmd.Args is nil")
		}

		// Test with wrong number of args
		err := newCmd.Args(newCmd, []string{"framework"})
		if err == nil {
			t.Error("Expected error for too few args, got nil")
		}

		// Test with correct number of args
		err = newCmd.Args(newCmd, []string{"framework", "project-name"})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Project command requires 1 arg", func(t *testing.T) {
		if projectCmd.Args == nil {
			t.Fatal("projectCmd.Args is nil")
		}

		// Test with no args
		err := projectCmd.Args(projectCmd, []string{})
		if err == nil {
			t.Error("Expected error for too few args, got nil")
		}

		// Test with correct number of args
		err = projectCmd.Args(projectCmd, []string{"project-name"})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
