package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain_Version_Integration(t *testing.T) {
	// Build the test binary
	cmd := exec.Command("go", "build", "-o", "peekenv_test.exe", ".")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove("peekenv_test.exe")

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "version flag short",
			args:     []string{"-v"},
			expected: "built on",
		},
		{
			name:     "version flag long",
			args:     []string{"--version"},
			expected: "built on",
		},
		{
			name:     "version argument",
			args:     []string{"version"},
			expected: "built on",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the built binary with test arguments
			cmd := exec.Command("./peekenv_test.exe", tt.args...)
			output, err := cmd.Output()

			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			outputStr := strings.TrimSpace(string(output))

			// Check that version output contains expected text
			if !strings.Contains(outputStr, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, outputStr)
			}
		})
	}
}
