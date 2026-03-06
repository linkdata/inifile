package inifile

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	iniSectionRegex = regexp.MustCompile(`^\[([^\]]*)\](?:\s*[;#].*)?$`)
	iniAssignRegex  = regexp.MustCompile(`^\s*([^=]+)\s*=\s*(.*)\s*$`)
)

// Parse reads INI data from an io.Reader and returns a new File.
//
// Unquoted values are trimmed of whitespace.
// Quoted values preserve whitespace inside the quotes.
// Parsed values are always valid UTF-8.
//
// If dupKeysJoin is zero, a duplicate key will replace the previous value.
// If dupKeysJoin is nonzero, a duplicate key will append it's value to
// the preexisting key's value using dupKeysJoin as a separator.
//
// You are allowed to pass in a nil io.Reader, which results in
// a nil File and error.
func Parse(r io.Reader, dupKeysJoin rune) (inif File, err error) {
	if r != nil {
		inif = make(File)
		scanner := bufio.NewScanner(r)
		section := ""
		lineNum := 0
		lineSrc := ""
		for err == nil && scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())
			if lineNum == 1 {
				// Trim UTF-8 BOM if present
				line = strings.TrimPrefix(line, "\uFEFF")
			}
			lineSrc = line
			if len(line) > 0 && line[0] != ';' && line[0] != '#' {
				if groups := iniAssignRegex.FindStringSubmatch(line); groups != nil {
					lineSrc = groups[2]
					var value string
					if value, err = ParseValue(groups[2]); err == nil {
						inif.Section(section).Set(groups[1], value, dupKeysJoin)
					}
				} else if groups := iniSectionRegex.FindStringSubmatch(line); groups != nil {
					section = Key(groups[1])
				} else {
					err = strconv.ErrSyntax
				}
			}
		}
		if err == nil {
			lineSrc = ""
			if err = scanner.Err(); err != nil {
				lineNum++
			}
		}
		if err != nil {
			inif = nil
			err = SyntaxError{
				Line:   lineNum,
				Source: lineSrc,
				Err:    err,
			}
		}
	}
	return
}
