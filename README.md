[![Build Status](https://github.com/tischda/peekenv/actions/workflows/build.yml/badge.svg)](https://github.com/tischda/peekenv/actions/workflows/build.yml)
[![Test Status](https://github.com/tischda/peekenv/actions/workflows/test.yml/badge.svg)](https://github.com/tischda/peekenv/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/tischda/peekenv/badge.svg)](https://coveralls.io/r/tischda/peekenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/tischda/peekenv/v3)](https://goreportcard.com/report/github.com/tischda/peekenv/v3)

# peekenv

Retrieves environment variables from the Windows registry.

### Install

~~~
go install github.com/tischda/peekenv/v3@latest
~~~

### Usage

~~~
Usage: peekenv [OPTIONS] [variables...]

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
~~~

### Examples

~~~
❯ peekenv psmodulepath
[PSModulePath]
%ProgramFiles%\WindowsPowerShell\Modules
%SystemRoot%\system32\WindowsPowerShell\v1.0\Modules
~~~

~~~
❯ peekenv -expand psmodulepath
[PSModulePath]
C:\Program Files\WindowsPowerShell\Modules
C:\WINDOWS\system32\WindowsPowerShell\v1.0\Modules                  
~~~

Note that path values are converted to multiples lines within the section.
This is the input format used by [pokenv](https://github.com/tischda/pokenv). 

### Alternatives

Built-in, see: `reg query /?`
