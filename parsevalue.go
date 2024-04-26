package inifile

import (
	"strconv"
	"strings"
)

// ParseValue parses a value string from an INI-file, unquoting and removing trailing comments if needed.
func ParseValue(s string) (value string, err error) {
	value = strings.TrimSpace(s)
	if strings.IndexAny(value, `'"`) == 0 {
		if value, err = strconv.QuotedPrefix(value); err == nil {
			value, err = strconv.Unquote(value)
		}
	} else {
		if idx := strings.IndexAny(value, ";#"); idx >= 0 {
			value = strings.TrimSpace(value[:idx])
		}
	}
	return
}
