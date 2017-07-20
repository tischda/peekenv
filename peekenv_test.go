package main

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

var sut peekenv

func init() {
	log.SetOutput(ioutil.Discard)
}

var testSections = []Section{
	{title: `PATH`, lines: []string{`C:\Program Files\ConEmu`, `C:\Program Files\ConEmu\ConEmu`, `C:\Windows\SYSTEM32`, `C:\Windows`}},
	{title: `TEMP`, lines: []string{`%USERPROFILE%\AppData\Local\Temp`}},
	{title: `TMP`, lines: []string{`c:\temp`}},
}

func TestGetSections(t *testing.T) {
	sut = peekenv{registry: mock}
	expected := testSections
	sut.populateSectionsFrom(REG_KEY_USER)
	if !reflect.DeepEqual(expected, sut.sections) {
		t.Errorf("Expected: %q, was: %q", expected, sut.sections)
	}
}

func TestStringer(t *testing.T) {
	sut = peekenv{
		sections: testSections,
		registry: mock,
	}
	expected := `[PATH]
C:\Program Files\ConEmu
C:\Program Files\ConEmu\ConEmu
C:\Windows\SYSTEM32
C:\Windows

[TEMP]
%USERPROFILE%\AppData\Local\Temp

[TMP]
c:\temp
`
	assertEquals(t, expected, sut.String())
}

func assertEquals(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}
}
