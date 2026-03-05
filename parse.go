package inifile

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

var (
	iniSectionRegex = regexp.MustCompile(`^\[([^\]]*)\](?:\s*[;#].*)?$`)
	iniAssignRegex  = regexp.MustCompile(`^([^=]+)=(.*)$`)
	ErrNilReader    = errors.New("nil reader")
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
// r must be a non-nil reader. Passing a nil reader returns ErrNilReader.
func Parse(r io.Reader, dupKeysJoin rune) (File, error) {
	if r == nil {
		return nil, ErrNilReader
	}
	inif := make(File)
	scanner := bufio.NewScanner(r)
	section := ""
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if lineNum == 1 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}
		if len(line) > 0 && line[0] != ';' && line[0] != '#' {
			if groups := iniAssignRegex.FindStringSubmatch(line); groups != nil {
				value, err := ParseValue(groups[2])
				if err != nil {
					return nil, SyntaxError{
						Line:   lineNum,
						Source: strings.TrimSpace(groups[2]),
						Err:    err,
					}
				}
				inif.Section(section).Set(groups[1], value, dupKeysJoin)
			} else if groups := iniSectionRegex.FindStringSubmatch(line); groups != nil {
				section = Key(groups[1])
			} else {
				return nil, SyntaxError{
					Line:   lineNum,
					Source: line,
				}
			}
		}
	}
	err := scanner.Err()
	if err != nil {
		inif = nil
	}
	return inif, err
}
