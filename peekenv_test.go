package main

import (
	"os"
	"strings"
	"testing"
)

func TestPeekenv_String(t *testing.T) {
	tests := []struct {
		name     string
		envMap   map[string]string
		expected string
	}{
		{
			name: "single variable",
			envMap: map[string]string{
				"TEMP": "C:\\Temp",
			},
			expected: "[TEMP]\nC:\\Temp\n",
		},
		{
			name: "multiple variables",
			envMap: map[string]string{
				"TEMP": "C:\\Temp",
				"USER": "johndoe",
			},
			expected: "[TEMP]\nC:\\Temp\n\n[USER]\njohndoe\n",
		},
		{
			name: "Path variable with semicolons",
			envMap: map[string]string{
				"Path": "C:\\Windows\\System32;C:\\Windows;C:\\Program Files\\Git\\bin",
			},
			expected: "[Path]\nC:\\Windows\\System32\nC:\\Windows\nC:\\Program Files\\Git\\bin\n",
		},
		{
			name: "Path and other variables",
			envMap: map[string]string{
				"Path": "C:\\Windows;C:\\Program Files",
				"TEMP": "C:\\Temp",
			},
			expected: "[Path]\nC:\\Windows\nC:\\Program Files\n\n[TEMP]\nC:\\Temp\n",
		},
		{
			name:     "empty map",
			envMap:   map[string]string{},
			expected: "\n",
		},
		{
			name: "variable with empty value",
			envMap: map[string]string{
				"EMPTY": "",
			},
			expected: "[EMPTY]\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &peekenv{
				envMap: tt.envMap,
			}
			result := p.String()

			// For tests with multiple variables, we need to handle the fact that
			// map iteration order is not guaranteed in Go
			if len(tt.envMap) > 1 && !strings.Contains(tt.name, "Path") {
				// Check that all expected sections are present
				for key, value := range tt.envMap {
					expectedSection := "[" + key + "]\n" + value
					if !strings.Contains(result, expectedSection) {
						t.Errorf("String() = %q, missing section for %s", result, key)
					}
				}
				// Check that result ends with newline
				if !strings.HasSuffix(result, "\n") {
					t.Errorf("String() = %q, should end with newline", result)
				}
			} else {
				// For single variable or Path tests, we can do exact comparison
				if result != tt.expected {
					t.Errorf("String() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

func TestPeekenv_ExportEnv_Both(t *testing.T) {
	// Create a temporary file for output
	tmpFile, err := os.CreateTemp("", "peekenv_test_both_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Will run after Close() due to LIFO
	defer tmpFile.Close()           // Will run first, ensuring file is closed before removal

	// Create a peekenv instance that will read from real registry
	p := &peekenv{
		envMap:    make(map[string]string),
		variables: []string{}, // No filters, read all variables
	}

	cfg := &Config{
		output: tmpFile.Name(),
		header: true,
		expand: false,
	}

	// Execute the test with real registry reading
	err = p.exportEnv(cfg)

	if err != nil {
		t.Fatalf("exportEnv() error = %v", err)
	}

	// Read the output from the temporary file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}
	output := string(content)

	// Check for expected header content
	expectedHeaders := []string{
		"# HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment",
		"# HKEY_CURRENT_USER\\Environment",
		"# Exported on",
	}

	for _, expected := range expectedHeaders {
		if !strings.Contains(output, expected) {
			t.Errorf("Output should contain header %q, but got:\n%s", expected, output)
		}
	}

	// Verify that PATH variable exists and contains expected Windows system paths
	if !strings.Contains(output, "[Path]") {
		t.Error("Output should contain [Path] section")
	}

	// Check for typical Windows system paths that should be in PATH
	expectedSystemPaths := []string{
		"Windows\\System32",
		"Windows",
	}

	for _, expectedPath := range expectedSystemPaths {
		if !strings.Contains(output, expectedPath) {
			t.Errorf("PATH should contain system path %q", expectedPath)
		}
	}

	// Check for typical user PATH entry (WindowsApps is commonly in user PATH)
	if !strings.Contains(output, "WindowsApps") {
		t.Log("WindowsApps not found in PATH - this may be normal depending on system configuration")
	}

	// Verify that OS variable exists and contains Windows_NT
	if !strings.Contains(output, "[OS]") {
		t.Error("Output should contain [OS] section")
	}

	if !strings.Contains(output, "Windows_NT") {
		t.Error("OS variable should contain Windows_NT")
	}

	// Verify output format - should have sections with proper formatting
	lines := strings.Split(output, "\n")
	foundOSSection := false
	foundPathSection := false

	for _, line := range lines {
		if line == "[OS]" {
			foundOSSection = true
		}
		if line == "[Path]" {
			foundPathSection = true
		}
	}

	if !foundOSSection {
		t.Error("Should have properly formatted [OS] section")
	}

	if !foundPathSection {
		t.Error("Should have properly formatted [Path] section")
	}
}

func TestPeekenv_ExportEnv_Machine_Expand_Windir(t *testing.T) {
	// Create a temporary file for output
	tmpFile, err := os.CreateTemp("", "peekenv_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Will run after Close() due to LIFO
	defer tmpFile.Close()           // Will run first, ensuring file is closed before removal

	// Create a peekenv instance that will read from real registry
	p := &peekenv{
		envMap:    make(map[string]string),
		variables: []string{"windir"}, // Filter for only windir variable
	}

	cfg := &Config{
		machine: true, // Read only machine variables
		output:  tmpFile.Name(),
		header:  false, // No header for cleaner output
		expand:  true,  // Expand environment variables
	}

	// Execute the test with machine registry reading only
	err = p.exportEnv(cfg)

	if err != nil {
		t.Fatalf("exportEnv() error = %v", err)
	}

	// Read the output from the temporary file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temporary file: %v", err)
	}
	output := string(content)

	// Verify that windir variable exists and is properly formatted
	if !strings.Contains(output, "[windir]") {
		t.Error("Output should contain [windir] section")
	}

	// Check that the expanded value contains Windows directory path
	expectedPaths := []string{
		"C:\\WINDOWS",
		"C:\\Windows", // Alternative casing
	}

	foundExpectedPath := false
	for _, expectedPath := range expectedPaths {
		if strings.Contains(output, expectedPath) {
			foundExpectedPath = true
			break
		}
	}

	if !foundExpectedPath {
		t.Errorf("windir should contain Windows directory path, got output:\n%s", output)
	}

	// Verify the output format matches expected structure
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		t.Errorf("Output should have at least 2 lines (section header + value), got %d lines", len(lines))
	}

	// First line should be the section header
	if lines[0] != "[windir]" {
		t.Errorf("First line should be [windir], got %q", lines[0])
	}

	// Second line should contain the Windows path
	if !strings.Contains(strings.ToUpper(lines[1]), "WINDOWS") {
		t.Errorf("Second line should contain Windows path, got %q", lines[1])
	}
}
