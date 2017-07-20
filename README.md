# peekenv [![Build status](https://ci.appveyor.com/api/projects/status/4but7lwfch3n65h0?svg=true)](https://ci.appveyor.com/project/tischda/peekenv)

Windows utility written in [Go](https://www.golang.org) to peek
environment variables from the registry.

### Install

There are no dependencies.

~~~
go get github.com/tischda/peekenv
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

### Other readers

Built-in, see: `reg query /?`
