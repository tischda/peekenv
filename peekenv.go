//go:build windows

package main

import (
	"fmt"
	"io"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"
)

var modkernel32 = syscall.NewLazyDLL("kernel32.dll")
var procExpandEnvironmentStringsW = modkernel32.NewProc("ExpandEnvironmentStringsW")

// peekenv handles the reading and formatting of environment variables.
// It maintains sections of environment variables, optional variable name filters,
// and a registry interface for data access.
type peekenv struct {
	envMap    map[string]string
	variables []string
}

// getUserAndSystemEnv retrieves the current environment
func (p *peekenv) getUserAndSystemEnv() error {
	p.envMap = make(map[string]string)

	// Read SYSTEM environment vars
	sysReg, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err == nil {
		defer sysReg.Close()
		sysEnv, _ := sysReg.ReadValueNames(0)
		for _, name := range sysEnv {
			val, _, _ := sysReg.GetStringValue(name)
			p.envMap[name] = val
		}
	}

	// Read USER environment vars
	userReg, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.READ)
	if err == nil {
		defer userReg.Close()
		userEnv, _ := userReg.ReadValueNames(0)
		for _, name := range userEnv {
			val, _, _ := userReg.GetStringValue(name)
			if name == "Path" {
				// Append USER Path to SYSTEM Path (system first, then user)
				p.envMap[name] = p.envMap[name] + ";" + val
			} else {
				p.envMap[name] = val
			}
		}
	}

	//TODO: keep only filtered variables

	// for _, sectionTitle := range values {
	// 	if len(p.variables) > 0 && !containsIgnoreCase(p.variables, sectionTitle) {
	// 		continue
	// 	}

	return nil
}

// ExportEnv reads environment variables from the registry and writes them to the provided writer.
// Parameters:
//   - reg: the registry key to read environment variables from
//   - w: the writer to output the formatted environment variables to
//   - printHeader: if true, includes a header with registry path and timestamp
//
// Returns an error if reading from registry fails or if no environment variables are found.
func (p *peekenv) exportEnv(reg RegistryMode, w io.Writer, printHeader bool) error {
	if err := p.getUserAndSystemEnv(); err != nil {
		return fmt.Errorf("retrieving variables: %w", err)
	}

	if len(p.envMap) == 0 {
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

// String returns a formatted string representation of all sections.
// Each section is formatted as [title] followed by its lines,
// with sections separated by double newlines.
func (p *peekenv) String() string {
	var sb strings.Builder
	for k, v := range p.envMap {
		if true {
			sb.WriteString("\n\n")
		}
		// TODO: only for Path, replace ; with newlines
		fmt.Fprintf(&sb, "[%s]\n%s", k, v)
	}
	sb.WriteString("\n")
	return sb.String()
}

// printInfoHeader writes a header comment to the writer containing the registry path and export timestamp.
// Parameters:
//   - reg: the registry key to include in the header
//   - w: the writer to output the header to
//
// Returns an error if the registry key type is unknown or if writing fails.
func printInfoHeader(reg RegistryMode, w io.Writer) error {
	now := time.Now().Format(time.RFC3339)
	var header string
	switch reg {
	case USER:
		header = fmt.Sprintf("# HKEY_CURRENT_USER\\TODO - Exported %s\n\n", now)
	case MACHINE:
		header = fmt.Sprintf("# HKEY_LOCAL_MACHINE\\TODO - Exported %s\n\n", now)
	default:
		// TODO
	}
	_, err := io.WriteString(w, header)
	return err
}

// containsIgnoreCase checks if a string slice contains a target string using case-insensitive comparison.
// Parameters:
//   - slice: the string slice to search in
//   - str: the target string to search for
//
// Returns true if the target string is found (case-insensitive), false otherwise.
func containsIgnoreCase(slice []string, str string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, str) {
			return true
		}
	}
	return false
}
