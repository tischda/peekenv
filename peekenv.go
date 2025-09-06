//go:build windows

package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// peekenv handles the reading and formatting of environment variables.
// It maintains a map of environment variables, and optional variable name filters.
type peekenv struct {
	envMap    map[string]string
	variables []string
}

// getUserAndSystemEnv retrieves all environment variables from the registry,
// merging USER and SYSTEM variables, with SYSTEM taking precedence for "Path".
func (p *peekenv) getUserAndSystemEnv() error {

	// Read SYSTEM environment vars
	if err := p.getSystemVariables(); err != nil {
		return err
	}

	// Read USER environment vars
	if err := p.getUserVariables(true); err != nil {
		return err
	}
	return nil
}

// getSystemVariables reads system environment variables from the registry
func (p *peekenv) getSystemVariables() error {
	sysReg, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err == nil {
		defer sysReg.Close()
		err = p.getVariables(sysReg, false)
	}
	return err
}

// getUserVariables reads user environment variables from the registry
func (p *peekenv) getUserVariables(mergePaths bool) error {
	userReg, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.READ)
	if err == nil {
		defer userReg.Close()
		err = p.getVariables(userReg, mergePaths)
	}
	return err
}

// getVariables reads environment variables from the provided registry key
//
// The mergePaths flag indicates if "Path" variables should be merged
// with existing values in p.envMap.
// This presuposes p.envMap is already initialized with SYSTEM variables.
func (p *peekenv) getVariables(reg registry.Key, mergePaths bool) error {
	env, err := reg.ReadValueNames(0)
	for _, variable := range env {
		if len(p.variables) > 0 && !containsIgnoreCase(p.variables, variable) {
			continue
		}
		val, _, _ := reg.GetStringValue(variable)
		if mergePaths && (variable == "Path" || variable == "PsModulePath") {
			// Append USER Path to SYSTEM Path (system first, then user)
			p.envMap[variable] = p.envMap[variable] + ";" + val
		} else {
			p.envMap[variable] = val
		}
	}
	return err
}

// ExportEnv reads environment variables from the registry and writes them to the provided writer.
// Parameters:
//   - reg: the registry key to read environment variables from
//   - w: the writer to output the formatted environment variables to
//   - printHeader: if true, includes a header with registry path and timestamp
//
// Returns an error if reading from registry fails or if no environment variables are found.
func (p *peekenv) exportEnv(reg RegistryMode, w io.Writer, printHeader bool) error {
	header := ""
	now := time.Now().Format(time.RFC3339)

	switch reg {
	case USER:
		if err := p.getUserVariables(false); err != nil {
			return fmt.Errorf("reading user environment variables: %w", err)
		}
		header = fmt.Sprintf("# HKEY_CURRENT_USER\\Environment - Exported on %s\n\n", now)
	case MACHINE:
		if err := p.getSystemVariables(); err != nil {
			return fmt.Errorf("reading system environment variables: %w", err)
		}
		header = fmt.Sprintf("# HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment - Exported on %s\n\n", now)
	default:
		if err := p.getUserAndSystemEnv(); err != nil {
			return fmt.Errorf("reading environment variables: %w", err)
		}
		header = "# HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment\n" + "# HKEY_CURRENT_USER\\Environment"
		header += fmt.Sprintf("Exported on %s\n\n", now)
	}

	if len(p.envMap) == 0 {
		return fmt.Errorf("no environment variables found")
	}

	if printHeader {
		_, err := io.WriteString(w, header)
		if err != nil {
			return fmt.Errorf("writing header: %w", err)
		}
	}

	_, err := io.WriteString(w, p.String())
	return err
}

// String returns a formatted string representation of all variables.
//
// [M2_HOME]
// c:\usr\bin\maven

// [Path]
// c:\Windows\system32
// c:\Windows
//
// Path type variables with multiple values separated by semicolons will be printed separated by
// newlines for better readability. This is also the format expected when importing with pokenv.
func (p *peekenv) String() string {
	var sb strings.Builder

	// Get sorted keys for consistent alphabetical output
	keys := make([]string, 0, len(p.envMap))
	for k := range p.envMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString("[" + k + "]\n")
		sb.WriteString(strings.ReplaceAll(p.envMap[k], ";", "\n"))
	}
	sb.WriteString("\n")
	return sb.String()
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
