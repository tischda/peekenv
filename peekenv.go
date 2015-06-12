package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

var PATH_USER = regPath{HKEY_CURRENT_USER, `Environment`}
var PATH_MACHINE = regPath{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type peekenv struct {
	sections []Section
	variable string
	registry Registry
}

type Section struct {
	title string
	lines []string
}

func (p *peekenv) exportEnv(path regPath, filename string) {
	p.populateSectionsFrom(path)
	p.writeFile(path, filename)
}

func (p *peekenv) populateSectionsFrom(path regPath) {
	values := p.registry.EnumValues(path)
	if len(values) > 0 {
		sort.Strings(values)
	}
	for _, sectionTitle := range values {
		if p.variable != "" && !strings.EqualFold(p.variable, sectionTitle) {
			continue
		}
		data, err := p.registry.GetString(path, sectionTitle)
		checkFatal(err)
		section := Section{title: sectionTitle, lines: strings.Split(data, ";")}
		p.sections = append(p.sections, section)
	}
}

func (p *peekenv) writeFile(path regPath, filename string) {
	var file *os.File
	if filename == "-" {
		file = os.Stdout
	} else {
		f, err := os.Create(filename)
		checkFatal(err)
		file = f
		defer file.Close()
	}

	now := fmt.Sprintf(" - Exported %s\n\n", time.Now())
	switch path.hKeyIdx {
	case HKEY_CURRENT_USER:
		io.WriteString(file, `# HKEY_CURRENT_USER\`+path.lpSubKey+now)
	case HKEY_LOCAL_MACHINE:
		io.WriteString(file, `# HKEY_LOCAL_MACHINE\`+path.lpSubKey+now)
	}
	io.WriteString(file, p.String())
}

func (p *peekenv) String() string {
	var result string
	for _, section := range p.sections {
		result += fmt.Sprintf("[%s]\n", section.title)
		for _, line := range section.lines {
			result += fmt.Sprintf("%s\n", line)
		}
		result += "\n"
	}
	return result
}

func checkFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
