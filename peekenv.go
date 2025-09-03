package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

var (
	REG_KEY_USER    = regKey{HKEY_CURRENT_USER, `Environment`}
	REG_KEY_MACHINE = regKey{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}
)

// Section represents a named group of environment variables
type Section struct {
	title string
	lines []string
}

// peekenv handles the reading and formatting of environment variables
type peekenv struct {
	sections []Section
	filters  []string
	registry Registry
}

// ExportEnv reads environment variables from the registry and writes them to the provided writer
func (p *peekenv) exportEnv(reg regKey, w io.Writer, printHeader bool) error {
	if err := p.populateSectionsFrom(reg); err != nil {
		return fmt.Errorf("populating sections: %w", err)
	}

	if len(p.sections) == 0 {
		return fmt.Errorf("no environment variables found")
	}

	if printHeader {
		if err := printInfoHeader(reg, w); err != nil {
			return fmt.Errorf("printing header: %w", err)
		}
	}

	_, err := io.WriteString(w, p.String())
	return err
}

func (p *peekenv) populateSectionsFrom(reg regKey) error {
	values := p.registry.EnumValues(reg)

	sort.Strings(values)
	for _, sectionTitle := range values {
		if len(p.filters) > 0 && !containsIgnoreCase(p.filters, sectionTitle) {
			continue
		}

		data, err := p.registry.GetString(reg, sectionTitle)
		if err != nil {
			return fmt.Errorf("getting registry value %q: %w", sectionTitle, err)
		}

		section := Section{
			title: sectionTitle,
			lines: strings.Split(data, ";"),
		}
		p.sections = append(p.sections, section)
	}
	return nil
}

func (p *peekenv) String() string {
	var sb strings.Builder
	for i, section := range p.sections {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		fmt.Fprintf(&sb, "[%s]\n%s", section.title, strings.Join(section.lines, "\n"))
	}
	sb.WriteString("\n")
	return sb.String()
}

func printInfoHeader(reg regKey, w io.Writer) error {
	now := time.Now().Format(time.RFC3339)
	var header string
	switch reg.hKeyIdx {
	case HKEY_CURRENT_USER:
		header = fmt.Sprintf("# HKEY_CURRENT_USER\\%s - Exported %s\n\n", reg.lpSubKey, now)
	case HKEY_LOCAL_MACHINE:
		header = fmt.Sprintf("# HKEY_LOCAL_MACHINE\\%s - Exported %s\n\n", reg.lpSubKey, now)
	default:
		return fmt.Errorf("unknown registry key type: %v", reg.hKeyIdx)
	}
	_, err := io.WriteString(w, header)
	return err
}

func containsIgnoreCase(slice []string, str string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, str) {
			return true
		}
	}
	return false
}
