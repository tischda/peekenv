package main

type mockRegistry struct {}

var mock = mockRegistry{}

func (r mockRegistry) GetString(path regPath, valueName string) (value string, err error) {
	return `C:\Program Files\ConEmu;C:\Program Files\ConEmu\ConEmu;C:\Windows\SYSTEM32;C:\Windows;C:\Windows\SYSTEM32\WBEM`, nil
}

