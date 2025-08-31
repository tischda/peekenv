[![Build Status](https://github.com/tischda/peekenv/actions/workflows/build.yml/badge.svg)](https://github.com/tischda/peekenv/actions/workflows/build.yml)
[![Test Status](https://github.com/tischda/peekenv/actions/workflows/test.yml/badge.svg)](https://github.com/tischda/peekenv/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/tischda/peekenv/badge.svg)](https://coveralls.io/r/tischda/peekenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/tischda/peekenv)](https://goreportcard.com/report/github.com/tischda/peekenv)

# peekenv

Retrieves environment variables from the Windows registry.

### Install

There are no dependencies.

~~~
go install github.com/tischda/peekenv@latest
~~~

### Usage

~~~
Usage: peekenv [-h] [-m] [-f outfile] [variables...]

OPTIONS:
  -f string
        file to dump the variables from the Windows environment (default "REQUIRED")
  -h
  -help
        displays this help message
  -i
  -info
        print info header
  -m
  -machine
        specifies that the variables should be read system wide (HKEY_LOCAL_MACHINE)
  -v
  -version
        print version and exit
~~~

Example:

~~~
# peekenv.exe -m pathext    
[PATHEXT]                   
.COM                        
.EXE                        
.BAT                        
.CMD                        
.VBS                        
.VBE                        
.JS                         
.JSE                        
.WSF                        
.WSH                        
.MSC                        
~~~

Note that the value `.COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC` is converted to multiples lines.
This is the input format used by [pokenv](https://github.com/tischda/pokenv). 

### Alternatives

Built-in, see: `reg query /?`
