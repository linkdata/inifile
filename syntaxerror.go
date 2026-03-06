package inifile

import (
	"fmt"
	"strconv"
)

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
	var linetxt string
	err := e.Err
	if err == nil {
		err = strconv.ErrSyntax
	}
	if e.Line > 0 {
		linetxt = fmt.Sprintf("line %d: ", e.Line)
	}
	return fmt.Sprintf("%s%v: %q", linetxt, err, e.Source)
}
