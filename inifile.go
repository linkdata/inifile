package inifile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	iniSectionRegex = regexp.MustCompile(`^\[(.*)\]$`)
	iniAssignRegex  = regexp.MustCompile(`^([^=]+)=(.*)$`)
)

// ErrIniSyntax is returned when there is a syntax error in an INI file.
type ErrIniSyntax struct {
	Line   int
	Source string // The contents of the erroneous line, without leading or trailing whitespace
}

func (e ErrIniSyntax) Error() string {
	return fmt.Sprintf("invalid INI syntax on line %d: %s", e.Line, e.Source)
}

// A IniFile represents a parsed INI file.
// The sections keys are all lowercase and trimmed of whitespace.
// Values found outside of a named section are in the "" section.
type IniFile map[string]IniSection

// A IniSection represents a single section of an INI file.
// The value keys are all lowercase and trimmed of whitespace.
type IniSection map[string]string

// Section returns a named Section. A Section will be created if one does not already exist for the given name.
// The section name will be trimmed of whitespace and lowercased.
func (inif IniFile) Section(name string) (section IniSection) {
	name = inifKey(name)
	if section = inif[name]; section == nil {
		section = make(IniSection)
		inif[name] = section
	}
	return section
}

// Looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
func (inif IniFile) Get(section, key string) (value string, ok bool) {
	if s := inif[inifKey(section)]; s != nil {
		value, ok = s[inifKey(key)]
	}
	return
}

// Parse reads INI data from an io.Reader and stores the data in the IniFile.
// All section and key names will be lowercased and trimmed of whitespace.
// All values will be trimmed of whitespace.
// If dupKeysJoin is nonzero, a duplicate key will append it's value to
// the preexisting key's value using dupKeysJoin as a separator.
func (inif *IniFile) Parse(r io.Reader, dupKeysJoin rune) (err error) {
	*inif = make(IniFile)
	scanner := bufio.NewScanner(r)
	section := ""
	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++
		if len(line) > 0 && line[0] != ';' && line[0] != '#' {
			if groups := iniAssignRegex.FindStringSubmatch(line); groups != nil {
				key := inifKey(groups[1])
				val := strings.TrimSpace(groups[2])
				sec := inif.Section(section)
				if dupKeysJoin != 0 {
					if oldVal := sec[key]; oldVal != "" {
						val = oldVal + string(dupKeysJoin) + val
					}
				}
				sec[key] = val
			} else if groups := iniSectionRegex.FindStringSubmatch(line); groups != nil {
				section = inifKey(groups[1])
				inif.Section(section)
			} else {
				return ErrIniSyntax{lineNum, line}
			}
		}
	}
	return
}

// Load opens the given file and calls Parse.
func (inif *IniFile) Load(fn string, dupKeysJoin rune) (err error) {
	var f *os.File
	if f, err = os.Open(fn); err == nil {
		defer f.Close()
		err = inif.Parse(f, dupKeysJoin)
	}
	return
}

func inifKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
