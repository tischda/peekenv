// +build windows

package main

import (
	"testing"
)

func TestGetValue(t *testing.T) {
	var registry = realRegistry{}
	expected := `%USERPROFILE%\AppData\Local\Temp`
	actual, err := registry.GetString(PATH_USER, "TEMP")
	if err != nil {
		t.Errorf("Error in SetString", err)
	}
	assertEquals(t, expected, actual)
}

func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
