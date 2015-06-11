package main

import (
	"regexp"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

var PATH_USER = regPath{HKEY_CURRENT_USER, `Environment`}
var PATH_MACHINE = regPath{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type peekenv struct {
	registry Registry
}

func (p *peekenv) getEnv(path regPath, fileName string) {
	p.getVars(path)
	p.saveFile(fileName)
}

func (p *peekenv) getVars(path regPath) {
}

func (p *peekenv) saveFile(fileName string) {

}

