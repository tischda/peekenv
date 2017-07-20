// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var version string

var file_name string
var flag_help = flag.Bool("help", false, "displays this help message")
var flag_machine = flag.Bool("machine", false, "specifies that the variables should be read system wide (HKEY_LOCAL_MACHINE)")
var flag_info = flag.Bool("info", false, "print info header")
var flag_version = flag.Bool("version", false, "print version and exit")

func init() {
	flag.BoolVar(flag_help, "h", false, "")
	flag.BoolVar(flag_info, "i", false, "")
	flag.BoolVar(flag_machine, "m", false, "")
	flag.BoolVar(flag_version, "v", false, "")
	flag.StringVar(&file_name, "f", "REQUIRED", "file to dump the variables from the Windows environment")
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-m] [-f outfile] [variables...]\n\nOPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flag_version {
		fmt.Println("peekenv version", version)
		return
	}
	if *flag_help {
		flag.Usage()
		os.Exit(1)
	}
	process()
}

func process() {
	var file *os.File

	if file_name == "REQUIRED" {
		file = os.Stdout
	} else {
		var err error
		file, err = os.Create(file_name)
		checkFatal(err)
	}
	defer file.Close()

	peekenv := peekenv{
		filters:  flag.Args(),
		registry: realRegistry{},
	}

	if *flag_machine {
		peekenv.exportEnv(REG_KEY_MACHINE, file)
	} else {
		peekenv.exportEnv(REG_KEY_USER, file)
	}
}

func checkFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
