package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

var PATH_USER = regPath{HKEY_CURRENT_USER, `Environment`}
var PATH_MACHINE = regPath{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type peekenv struct {
	sections []Section
	registry Registry
}

type Section struct {
	title string
	lines []string
}

func (p *peekenv) getEnv(path regPath, fileName string) {
	file, err := os.Create(fileName)
	checkFatal(err)
	defer file.Close()

	p.populateSectionsFrom(path)
	io.WriteString(file, p.String())
}

func (p *peekenv) populateSectionsFrom(path regPath) {
	values := p.registry.EnumValues(path)
	if len(values) > 0 {
		sort.Strings(values)
	}
	for _, sectionTitle := range values {
		data, err := p.registry.GetString(path, sectionTitle)
		checkFatal(err)
		section := Section{title: sectionTitle, lines: strings.Split(data, ";")}
		p.sections = append(p.sections, section)
	}
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
