//go:build windows

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
	user    bool
	machine bool
	header  bool
	expand  bool
	output  string
	help    bool
	version bool
}

func initFlags() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.user, "u", false, "")
	flag.BoolVar(&cfg.user, "user", false, "read only user variables (HKEY_CURRENT_USER)")
	flag.BoolVar(&cfg.machine, "m", false, "")
	flag.BoolVar(&cfg.machine, "machine", false, "read only system variables (HKEY_LOCAL_MACHINE)")
	flag.BoolVar(&cfg.header, "h", false, "")
	flag.BoolVar(&cfg.header, "header", false, "print info header")
	flag.BoolVar(&cfg.expand, "x", false, "")
	flag.BoolVar(&cfg.expand, "expand", false, "expand environment variables to values (eg. %APPDATA%)")
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
          read user variables (HKEY_CURRENT_USER)"
  -m, --machine"
          read system variables (HKEY_LOCAL_MACHINE)"
  -h, --header
          print info header
  -x, --expand
          expand environment variables to values (eg. %APPDATA%)
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

	// Process the environment variables
	peekenv := peekenv{
		envMap:    make(map[string]string),
		variables: flag.Args(),
	}

	if err := peekenv.exportEnv(cfg); err != nil {
		log.Fatalln(err)
	}
}
