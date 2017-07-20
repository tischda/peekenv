package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var REG_KEY_USER = regKey{HKEY_CURRENT_USER, `Environment`}
var REG_KEY_MACHINE = regKey{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type peekenv struct {
	sections []Section
	filters  []string
	registry Registry
}

type Section struct {
	title string
	lines []string
}

func (p *peekenv) exportEnv(reg regKey, file *os.File) {
	p.populateSectionsFrom(reg)
	if len(p.sections) == 0 {
		log.Fatalln("Fatal: Nothing found")
	}
	if *flag_info {
		printInfoHeader(reg, file)
	}
	io.WriteString(file, p.String())
}

func (p *peekenv) populateSectionsFrom(reg regKey) {
	values := p.registry.EnumValues(reg)
	sort.Strings(values)
	for _, sectionTitle := range values {

		// check if the current section is amongst the requested variables
		if len(p.filters) > 0 && !contains(p.filters, sectionTitle) {
			continue
		}
		data, err := p.registry.GetString(reg, sectionTitle)
		checkFatal(err)
		section := Section{title: sectionTitle, lines: strings.Split(data, ";")}
		p.sections = append(p.sections, section)
	}
}

func (p *peekenv) String() string {
	var result string
	for _, section := range p.sections {
		result += fmt.Sprintf("[%s]\n%s\n\n", section.title, strings.Join(section.lines, "\n"))
	}
	return strings.TrimSuffix(result, "\n")
}

func printInfoHeader(reg regKey, file *os.File) {
	now := fmt.Sprintf(" - Exported %s\n\n", time.Now())
	switch reg.hKeyIdx {
	case HKEY_CURRENT_USER:
		io.WriteString(file, `# HKEY_CURRENT_USER\`+reg.lpSubKey+now)
	case HKEY_LOCAL_MACHINE:
		io.WriteString(file, `# HKEY_LOCAL_MACHINE\`+reg.lpSubKey+now)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
