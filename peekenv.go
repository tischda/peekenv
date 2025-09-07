//go:build windows

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// RegistryMode represents which registry keys to read from
type RegistryMode int

const (
	MACHINE RegistryMode = iota // Read only from HKEY_LOCAL_MACHINE
	USER                        // Read only from HKEY_CURRENT_USER
	BOTH                        // Read from both registries, user takes precedence
)

var (
	// Header strings for different registry modes
	headerStrings = map[RegistryMode]string{
		USER:    "# HKEY_CURRENT_USER\\Environment\n",
		MACHINE: "# HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment\n",
		BOTH:    "# HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment\n# HKEY_CURRENT_USER\\Environment\n",
	}
)

// peekenv handles the reading and formatting of environment variables.
// It maintains a map of environment variables, and which variables to export (if specified).
type peekenv struct {
	envMap    map[string]string
	variables []string
}

// exportEnv reads environment variables from the registry and writes them to the output.
//
// Parameters:
//   - cfg: the runtime configuration specifying registry mode, output options, etc.
//
// Returns an error if reading from registry fails or no environment variables are found.
func (p *peekenv) exportEnv(cfg *Config) error {

	mode := BOTH
	if cfg.machine && cfg.user {
		mode = BOTH
	} else if cfg.machine {
		mode = MACHINE
	} else if cfg.user {
		mode = USER
	}

	if err := p.readRegistry(mode); err != nil {
		return err
	}

	// Expand variables if requested
	if cfg.expand {
		for k, v := range p.envMap {
			p.envMap[k] = expandVariable(v)
		}
	}
	return p.writeOutput(cfg, mode)
}

// readRegistry reads environment variables from the Windows registry based on the specified mode.
//
// Parameters:
//   - mode: specifies which registry keys to read from (USER, MACHINE, or BOTH)
//
// Returns an error if registry access fails or no environment variables are found.
func (p *peekenv) readRegistry(mode RegistryMode) error {
	switch mode {
	case USER:
		if err := p.getUserVariables(false); err != nil {
			return fmt.Errorf("reading user environment variables: %w", err)
		}
	case MACHINE:
		if err := p.getSystemVariables(); err != nil {
			return fmt.Errorf("reading system environment variables: %w", err)
		}
	default:
		// order matters, first system, then user (so user can override)
		if err := p.getSystemVariables(); err != nil {
			return fmt.Errorf("reading system environment variables: %w", err)
		}
		if err := p.getUserVariables(true); err != nil {
			return fmt.Errorf("reading user environment variables: %w", err)
		}
	}

	if len(p.envMap) == 0 {
		return fmt.Errorf("no environment variables found")
	}
	return nil
}

// writeOutput writes the formatted environment variables to the specified output.
//
// Parameters:
//   - cfg: the runtime configuration containing output file path and header options
//   - mode: the registry mode used for header generation
//
// Returns an error if file creation, header writing, or variable writing fails.
func (p *peekenv) writeOutput(cfg *Config, mode RegistryMode) error {

	// Open output file or use stdout
	var err error
	var file *os.File
	if cfg.output == "stdout" {
		file = os.Stdout
	} else {
		file, err = os.Create(cfg.output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
	}
	defer file.Close()

	// Print header if requested
	if cfg.header {
		header := headerStrings[mode]
		now := time.Now().Format("2006-01-02 15:04:05 -0700 MST")
		header += fmt.Sprintf("# Exported on %s\n\n", now)
		_, err := io.WriteString(file, header)
		if err != nil {
			return fmt.Errorf("writing header: %w", err)
		}
	}

	// Print variables in proper format
	_, err = io.WriteString(file, p.String())
	return err
}

// getSystemVariables opens the HKEY_LOCAL_MACHINE registry key and populates
// p.envMap with system variables.
//
// Returns an error if the registry cannot be accessed or read.
func (p *peekenv) getSystemVariables() error {
	sysReg, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err == nil {
		defer sysReg.Close()
		err = p.getVariables(sysReg, false)
	}
	return err
}

// getUserVariables opens the HKEY_CURRENT_USER registry key and populates
// p.envMap with user variables
//
// Parameters:
//   - mergePaths: if true, merges "Path" and "PsModulePath" with existing values in p.envMap
//
// Returns an error if the registry cannot be accessed or read.
func (p *peekenv) getUserVariables(mergePaths bool) error {
	userReg, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.READ)
	if err == nil {
		defer userReg.Close()
		err = p.getVariables(userReg, mergePaths)
	}
	return err
}

// getVariables reads environment variables from the provided registry key.
//
// Parameters:
//   - reg: the registry key to read variables from
//   - mergePaths: if true, merges "Path" and "PsModulePath" with existing values in p.envMap
//
// Merging presupposes that p.envMap has already been initialized with SYSTEM variables.
// Therefore, call getSystemVariables() before calling this with mergePaths=true.
//
// Returns an error if the registry values cannot be read.
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

// String returns string representation of all variables, formatted like this:
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

	// Sort keys for consistent alphabetical output (case-insensitive)
	keyMap := make(map[string]string)
	keys := make([]string, 0, len(p.envMap))
	for k := range p.envMap {
		keys = append(keys, strings.ToLower(k))
		keyMap[strings.ToLower(k)] = k
	}
	sort.Strings(keys)

	// Format output with a section header for each variable
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		originalKey := keyMap[k]
		sb.WriteString("[" + originalKey + "]\n")
		sb.WriteString(strings.ReplaceAll(p.envMap[originalKey], ";", "\n"))
	}
	sb.WriteString("\n")
	return sb.String()
}

// containsIgnoreCase checks if a string slice contains a target string using
// case-insensitive comparison.
//
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
