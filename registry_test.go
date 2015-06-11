// +build windows

package main

import (
	"testing"
)

func TestSetDeleteValue(t *testing.T) {

	var registry = realRegistry{}

	expected := `%USERPROFILE%\AppData\Local\Temp`

	// set value
	actual, err := registry.GetString(PATH_USER, "TEMP")
	if err != nil {
		t.Errorf("Error in SetString", err)
	}
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
