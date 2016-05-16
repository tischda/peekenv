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

