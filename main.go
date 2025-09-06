//go:build windows

// Package main provides a utility to dump Windows environment variables
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// https://goreleaser.com/cookbooks/using-main.version/
var (
	name    string
	version string
	date    string
	commit  string
)

// flags
type Config struct {
	help    bool
	user    bool
	machine bool
	header  bool
	version bool
	output  string
}

// RegistryMode represents which registry keys to read from
type RegistryMode int

const (
	MACHINE RegistryMode = iota // Read only from HKEY_LOCAL_MACHINE
	USER                        // Read only from HKEY_CURRENT_USER
	BOTH                        // Read from both registries, user takes precedence
)

func initFlags() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.user, "u", false, "")
	flag.BoolVar(&cfg.user, "user", false, "read only user variables (HKEY_CURRENT_USER)")
	flag.BoolVar(&cfg.machine, "m", false, "")
	flag.BoolVar(&cfg.machine, "machine", false, "read only system variables (HKEY_LOCAL_MACHINE)")
	flag.BoolVar(&cfg.header, "h", false, "")
	flag.BoolVar(&cfg.header, "header", false, "print info header")
	flag.StringVar(&cfg.output, "o", "stdout", "")
	flag.StringVar(&cfg.output, "output", "stdout", "file to dump the environment variables to")
	flag.BoolVar(&cfg.help, "?", false, "")
	flag.BoolVar(&cfg.help, "help", false, "displays this help message")
	flag.BoolVar(&cfg.version, "v", false, "")
	flag.BoolVar(&cfg.version, "version", false, "print version and exit")
	return cfg
}

func main() {
	log.SetFlags(0)
	cfg := initFlags()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+name+` [OPTIONS] [variables...]

Retrieves environment variables from the Windows registry. By default,
both system and user variables are read. You can filter using OPTIONS.

If no variables are specified, all environment variables are printed.

OPTIONS:

  -u, --user"
        read only user variables (HKEY_CURRENT_USER)"
  -m, --machine"
        read only system variables (HKEY_LOCAL_MACHINE)"
  -h, --header
		print info header
  -o, --output FILE
		file to dump the environment variables to (default: stdout)
  -?, --help
        display this help message
  -v, --version
        print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  $ "+name+` TEMP
      [TEMP]
      c:\temp`)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || cfg.version {
		fmt.Printf("%s %s, built on %s (commit: %s)\n", name, version, date, commit)
		return
	}

	if cfg.help {
		flag.Usage()
		return
	}

	if cfg.machine && cfg.user {
		fmt.Printf("cannot specify both -m/--machine and -u/--user")
		os.Exit(1)
	}

	if err := process(cfg); err != nil {
		log.Fatalln("Error:", err)
	}
}

func process(cfg *Config) error {
	var file *os.File
	var err error

	if cfg.output == "stdout" {
		file = os.Stdout
	} else {
		file, err = os.Create(cfg.output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
	}
	defer file.Close()

	peekenv := peekenv{
		variables: flag.Args(),
	}

	if cfg.machine {
		return peekenv.exportEnv(MACHINE, file, cfg.header)
	} else if cfg.user {
		return peekenv.exportEnv(USER, file, cfg.header)
	}
	return peekenv.exportEnv(BOTH, file, cfg.header)
}
