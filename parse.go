package inifile

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var (
	iniSectionRegex = regexp.MustCompile(`^\[(.*)\]$`)
	iniAssignRegex  = regexp.MustCompile(`^([^=]+)=(.*)$`)
)

// Parse reads INI data from an io.Reader and returns a new File.
//
// All values will be trimmed of whitespace.
//
// If dupKeysJoin is zero, a duplicate key will replace the previous value.
// If dupKeysJoin is nonzero, a duplicate key will append it's value to
// the preexisting key's value using dupKeysJoin as a separator.
//
// r must be a non-nil reader. Passing a nil reader causes a panic.
func Parse(r io.Reader, dupKeysJoin rune) (File, error) {
	inif := make(File)
	scanner := bufio.NewScanner(r)
	section := ""
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != ';' && line[0] != '#' {
			if groups := iniAssignRegex.FindStringSubmatch(line); groups != nil {
				value, err := ParseValue(groups[2])
				if err != nil {
					return nil, SyntaxError{lineNum, strings.TrimSpace(groups[2])}
				}
				inif.Section(section).Set(groups[1], value, dupKeysJoin)
			} else if groups := iniSectionRegex.FindStringSubmatch(line); groups != nil {
				section = Key(groups[1])
			} else {
				return nil, SyntaxError{lineNum, line}
			}
		}
	}
	err := scanner.Err()
	if err != nil {
		inif = nil
	}
	return inif, err
}
