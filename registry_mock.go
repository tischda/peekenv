package main

type mockRegistry struct {
	env map[string]string
}

var mock = mockRegistry{}

func init() {
	mock.env = map[string]string{
		`TEMP`: `%USERPROFILE%\AppData\Local\Temp`,
		`TMP`:  `c:\temp`,
		`PATH`: `C:\Program Files\ConEmu;C:\Program Files\ConEmu\ConEmu;C:\Windows\SYSTEM32;C:\Windows`,
	}
}

func (r mockRegistry) GetString(path regKey, valueName string) (value string, err error) {
	return r.env[valueName], nil
}

func (r mockRegistry) EnumValues(path regKey) []string {
	keys := make([]string, 0, len(r.env))
	for k := range r.env {
		keys = append(keys, k)
	}
	return keys
}
