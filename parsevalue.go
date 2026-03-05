package inifile

import (
	"strconv"
	"strings"
)

// ParseValue parses a value string from an INI-file, unquoting and removing trailing comments if needed.
func ParseValue(s string) (value string, err error) {
	if value = strings.TrimSpace(s); len(value) > 0 {
		switch value[0] {
		case '"':
			value, err = parseDoubleQuotedValue(value)
		case '\'':
			value, err = parseSingleQuotedValue(value)
		default:
			value = parseUnquotedValue(value)
		}
	}
	return
}

func validateTail(tail string) (err error) {
	tail = strings.TrimSpace(tail)
	if tail != "" && tail[0] != ';' && tail[0] != '#' {
		err = strconv.ErrSyntax
	}
	return
}

func parseDoubleQuotedValue(value string) (result string, err error) {
	var quoted string
	if quoted, err = strconv.QuotedPrefix(value); err == nil {
		if err = validateTail(value[len(quoted):]); err == nil {
			result, err = strconv.Unquote(quoted)
		}
	}
	return
}

func parseSingleQuotedValue(value string) (result string, err error) {
	var tail string
	if result, tail, err = parseSingleQuoted(value); err == nil {
		err = validateTail(tail)
	}
	return
}

func parseSingleQuoted(value string) (result, tail string, err error) {
	var b strings.Builder
	b.Grow(len(value))
	value = value[1:]
	for err == nil && len(value) > 0 {
		if value[0] == '\'' {
			result = b.String()
			tail = value[1:]
			return
		}
		var singlebyte rune
		var multibyte bool
		if singlebyte, multibyte, value, err = strconv.UnquoteChar(value, '\''); err == nil {
			if multibyte {
				b.WriteRune(singlebyte)
			} else {
				b.WriteByte(byte(singlebyte)) // #nosec G115
			}
		}
	}
	err = strconv.ErrSyntax
	return
}

func parseUnquotedValue(value string) (result string) {
	result = value
	if idx := strings.IndexAny(value, ";#"); idx >= 0 {
		result = strings.TrimSpace(value[:idx])
	}
	return
}
