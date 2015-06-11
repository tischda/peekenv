// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version string = "1.1.0"

func main() {
	hkcu := flag.String("hkcu", "REQUIRED", "process input file into HKEY_CURRENT_USER environment")
	hklm := flag.String("hklm", "REQUIRED", "process input file into HKEY_LOCAL_MACHINE environment")
	showVersion := flag.Bool("version", false, "print version and exit")

	// configure logging
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-checkpaths] [-hkcu|-hklm] infile\n  infile: the input file\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("pokenv version", version)
		return
	}
	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	p := peekenv{
		registry: realRegistry{},
	}

	if *hkcu != "REQUIRED" {
		p.getEnv(PATH_USER, *hkcu)
	}
	if *hklm != "REQUIRED" {
		p.getEnv(PATH_MACHINE, *hklm)
	}
}
