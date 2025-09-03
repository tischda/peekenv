//go:build windows

// Package main provides a utility to dump Windows environment variables
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "" // version is set during build

type Config struct {
	help     bool
	machine  bool
	info     bool
	version  bool
	filename string
}

func initFlags() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.help, "help", false, "displays this help message")
	flag.BoolVar(&cfg.help, "h", false, "")
	flag.BoolVar(&cfg.machine, "machine", false, "specifies that the variables should be read system wide (HKEY_LOCAL_MACHINE)")
	flag.BoolVar(&cfg.machine, "m", false, "")
	flag.BoolVar(&cfg.info, "info", false, "print info header")
	flag.BoolVar(&cfg.info, "i", false, "")
	flag.BoolVar(&cfg.version, "version", false, "print version and exit")
	flag.BoolVar(&cfg.version, "v", false, "")
	flag.StringVar(&cfg.filename, "f", "REQUIRED", "file to dump the variables from the Windows environment")
	return cfg
}

func main() {
	log.SetFlags(0)
	cfg := initFlags()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-m] [-f outfile] [variables...]\n\nOPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *Config) error {
	if cfg.version {
		fmt.Println("peekenv version", version)
		return nil
	}
	if cfg.help {
		flag.Usage()
		os.Exit(1)
	}
	return process(cfg)
}

func process(cfg *Config) error {
	var file *os.File
	var err error

	if cfg.filename == "REQUIRED" {
		file = os.Stdout
	} else {
		file, err = os.Create(cfg.filename)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
	}
	defer file.Close()

	peekenv := peekenv{
		filters:  flag.Args(),
		registry: realRegistry{},
	}

	if cfg.machine {
		return peekenv.exportEnv(REG_KEY_MACHINE, file, cfg.info)
	}
	return peekenv.exportEnv(REG_KEY_USER, file, cfg.info)
}
