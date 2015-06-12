// +build windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMainTimer(t *testing.T) {
	args := []string{"-version"}
	os.Args = append(os.Args, args...)

	expected := fmt.Sprintf("peekenv version %s\n", version)

	// this can be done only once or test framework will panic
	actual := captureOutput(main)

	if expected != actual {
		t.Errorf("Expected: %s, but was: %s", expected, actual)
	}
}

// captures Stdout and returns output of function f()
func captureOutput(f func()) string {
	// redirect output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	// reset output again
	w.Close()
	os.Stdout = old

	captured, _ := ioutil.ReadAll(r)
	return string(captured)
}
