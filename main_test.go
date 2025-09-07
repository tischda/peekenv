//go:build windows

package main

import (
	"flag"
	"os"
	"testing"
)

func TestInitFlags(t *testing.T) {
	// Save original command line and reset flags
	originalArgs := os.Args
	originalCommandLine := flag.CommandLine

	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalCommandLine
	}()

	// Create a new flag set for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test initFlags() function
	cfg := initFlags()

	// Test default values
	if cfg.user != false {
		t.Errorf("Expected user default to be false, got %v", cfg.user)
	}
	if cfg.machine != false {
		t.Errorf("Expected machine default to be false, got %v", cfg.machine)
	}
	if cfg.header != false {
		t.Errorf("Expected header default to be false, got %v", cfg.header)
	}
	if cfg.expand != false {
		t.Errorf("Expected expand default to be false, got %v", cfg.expand)
	}
	if cfg.output != "stdout" {
		t.Errorf("Expected output default to be 'stdout', got %v", cfg.output)
	}
	if cfg.help != false {
		t.Errorf("Expected help default to be false, got %v", cfg.help)
	}
	if cfg.version != false {
		t.Errorf("Expected version default to be false, got %v", cfg.version)
	}

	// Test that flags can be parsed
	testArgs := []string{
		"peekenv",
		"-u",
		"-m",
		"-h",
		"-x",
		"-o", "test.txt",
		"-v",
	}

	// Reset flag set and reinitialize
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg = initFlags()

	// Parse test arguments
	err := flag.CommandLine.Parse(testArgs[1:])
	if err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	// Verify flags were set correctly
	if !cfg.user {
		t.Error("Expected user flag to be true")
	}
	if !cfg.machine {
		t.Error("Expected machine flag to be true")
	}
	if !cfg.header {
		t.Error("Expected header flag to be true")
	}
	if !cfg.expand {
		t.Error("Expected expand flag to be true")
	}
	if cfg.output != "test.txt" {
		t.Errorf("Expected output to be 'test.txt', got %v", cfg.output)
	}
	if !cfg.version {
		t.Error("Expected version flag to be true")
	}
}
