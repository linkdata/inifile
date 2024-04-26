package inifile

import (
	"os"
	"strings"
)

// Key returns s lowercased and with whitespace trimmed.
func Key(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Load reads INI data from the given file and calls Parse on it.
func Load(fn string, dupKeysJoin rune) (inif File, err error) {
	var f *os.File
	f, err = os.Open(fn) // #nosec G304
	if err == nil {
		defer f.Close()
		inif, err = Parse(f, dupKeysJoin)
	}
	return
}
