package inifile

import "fmt"

// ErrSyntax is returned when there is a syntax error in an INI file.
var ErrSyntax SyntaxError

type SyntaxError struct {
	Line   int    // Line number where the error was found
	Source string // The contents of the erroneous line, without leading or trailing whitespace
	Err    error  // Optional source error
}

func (e SyntaxError) Is(target error) bool {
	return target == ErrSyntax
}

func (e SyntaxError) Unwrap() error {
	return e.Err
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("invalid INI syntax on line %d: %s", e.Line, e.Source)
}
