// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version string = "1.0.0"

func main() {
	hkcu := flag.String("hkcu", "REQUIRED", "extract HKEY_CURRENT_USER environment to file")
	hklm := flag.String("hklm", "REQUIRED", "extract HKEY_LOCAL_MACHINE environment to file")
	showVersion := flag.Bool("version", false, "print version and exit")

	// configure logging
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-hkcu|-hklm] outfile\n  outfile: the output file\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("peekenv version", version)
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
