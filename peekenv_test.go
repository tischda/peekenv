package main

import (
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
