package inifile

import "fmt"

// ErrSyntax is returned when there is a syntax error in an INI file.
var ErrSyntax SyntaxError

type SyntaxError struct {
	Line   int    // Line number where the error was found
	Source string // The contents of the erroneous line, without leading or trailing whitespace
}

func (e SyntaxError) Is(err error) (ok bool) {
	_, ok = err.(SyntaxError)
	return
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("invalid INI syntax on line %d: %s", e.Line, e.Source)
}
