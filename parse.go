package inifile

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	iniSectionRegex = regexp.MustCompile(`^\[([^\]\r\n]*)\](?:[ \t]*[;#][^\r\n]*)?$`)
	iniAssignRegex  = regexp.MustCompile(`^[ \t]*([^\r\n=]+)[ \t]*=[ \t]*([^\r\n]*)[ \t]*$`)
	iniUTF8BOM      = []byte("\uFEFF")
)

func scanLinesWithLineNumbers(scanLineNum, tokenLineNum *int) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if token != nil {
			*tokenLineNum = *scanLineNum
			if *tokenLineNum == 1 {
				token = bytes.TrimPrefix(token, iniUTF8BOM)
			}
		}
		if advance > 0 {
			*scanLineNum += bytes.Count(data[:advance], []byte{'\n'})
		}
		return
	}
}

// Parse reads INI data from an io.Reader and returns a new File.
//
// Unquoted values are trimmed of whitespace.
// Quoted values preserve whitespace inside the quotes.
// Parsed values are always valid UTF-8.
// Input must use LF (\n) or CRLF (\r\n) line endings; CR-only (\r) is not supported.
// Lines are limited by bufio.Scanner's default token size (~64 KiB).
//
// If dupKeysJoin is zero, a duplicate key will replace the previous value.
// If dupKeysJoin is nonzero, a duplicate key will append it's value to
// the preexisting key's value using dupKeysJoin as a separator.
//
// You are allowed to pass in a nil io.Reader, which results in
// a nil File and no error.
func Parse(r io.Reader, dupKeysJoin rune) (inif File, err error) {
	if r != nil {
		scanLineNum := 1
		tokenLineNum := 1
		section := ""
		line := ""
		inif = make(File)
		scanner := bufio.NewScanner(r)
		scanner.Split(scanLinesWithLineNumbers(&scanLineNum, &tokenLineNum))
		for err == nil && scanner.Scan() {
			line = strings.TrimSpace(scanner.Text())
			if len(line) > 0 && line[0] != ';' && line[0] != '#' {
				if groups := iniAssignRegex.FindStringSubmatch(line); groups != nil {
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
			if err = scanner.Err(); err != nil {
				inif = nil
			}
		} else {
			inif = nil
			err = SyntaxError{
				Line:   tokenLineNum,
				Source: line,
				Err:    err,
			}
		}
	}
	return
}
