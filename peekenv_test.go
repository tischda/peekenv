package main

import (
	"io/ioutil"
	"log"
	"testing"
)

var sut peekenv

func init() {
	sut = peekenv{
		registry:    mock,
	}
	log.SetOutput(ioutil.Discard)
}

func TestProcessLineValue(t *testing.T) {

}


func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}

