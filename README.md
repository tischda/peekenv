# peekenv [![Build status](https://ci.appveyor.com/api/projects/status/hrtwo6hrx10d7i88?svg=true)](https://ci.appveyor.com/project/tischda/peekenv)

Windows utility written in [Go](https://www.golang.org) to peek
environment variables from the registry.

### Compile

Tested with GO 1.4.2. There are no dependencies.

~~~
go build
~~~

### Usage

~~~
Usage: peekenv.exe [-var variable] [-hkcu|-hklm] outfile
  outfile: the output file (use '-' for stdout)
  -hkcu="REQUIRED": extract HKEY_CURRENT_USER environment to file
  -hklm="REQUIRED": extract HKEY_LOCAL_MACHINE environment to file
  -var="": extract only specified variable
  -version=false: print version and exit
~~~

Example:

~~~
u:\>peekenv.exe -var pathext -hklm -
# HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment - Exported 2015-06-12 10:34:30.7803973 +0200 CEST

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
.RB
.RBW
~~~

Note that the value `.COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC;.RB;.RBW` is converted to multiples lines.
This is the input format used by [pokenv](https://github.com/tischda/pokenv). 

### Other readers

Built-in, see: `reg query /?`
