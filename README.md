[![build](https://github.com/linkdata/inifile/actions/workflows/go.yml/badge.svg)](https://github.com/linkdata/inifile/actions/workflows/go.yml)
[![coverage](https://coveralls.io/repos/github/linkdata/inifile/badge.svg?branch=main)](https://coveralls.io/github/linkdata/inifile?branch=main)
[![goreport](https://goreportcard.com/badge/github.com/linkdata/inifile)](https://goreportcard.com/report/github.com/linkdata/inifile)
[![Docs](https://godoc.org/github.com/linkdata/inifile?status.svg)](https://godoc.org/github.com/linkdata/inifile)

# inifile

Simple INI file reader.

Section and key names are case-insensitive and ignore leading and trailing whitespace.

Supports line comments, trailing comments and quoted values.

Keys appearing more than once will either replace previous values or append to them with a separator.

```go
func ExampleParse() {
	const initext = `
# comments start with a hash
; or a semicolon

# global keys are in the unnamed (empty string) section
Username = " user name " # values can be quoted preserve whitespace or embed quotes

# section names are case insensitive and ignores leading and trailing whitespace
[ HTTP ] 
 port = 8080 # keys and values are stripped of leading and trailing whitespace
port=8081 # keys appearing more than once either append or replace values
`

	inif, err := inifile.Parse(strings.NewReader(initext), ',')
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	username, _ := inif.Get("", "username")
	ports, _ := inif.Get("http", "port")

	fmt.Printf("%q\n%q\n", username, ports)
	// Output:
	// " user name "
	// "8080,8081"
}
```