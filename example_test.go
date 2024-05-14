package inifile_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/linkdata/inifile"
)

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
