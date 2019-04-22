// +build windows

package main

import (
	"testing"
)

func TestGetValue(t *testing.T) {
	var registry = realRegistry{}
	expected := `AMD64`
	actual, err := registry.GetString(REG_KEY_MACHINE, "PROCESSOR_ARCHITECTURE")
	if err != nil {
		t.Errorf("Error in SetString: %q", err)
	}
	assertEquals(t, expected, actual)
}
