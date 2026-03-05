package inifile

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func TestSyntaxErrorIs(t *testing.T) {
	err := SyntaxError{
		Line:   7,
		Source: "'broken",
		Err:    strconv.ErrSyntax,
	}

	if !errors.Is(err, ErrSyntax) {
		t.Fatalf("errors.Is(err, ErrSyntax) = false, want true")
	}
	if !errors.Is(err, SyntaxError{}) {
		t.Fatalf("errors.Is(err, SyntaxError{}) = false, want true")
	}
	if errors.Is(err, SyntaxError{Line: 7, Source: "'broken"}) {
		t.Fatalf("errors.Is(err, SyntaxError{Line, Source}) = true, want false")
	}
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("errors.Is(err, strconv.ErrSyntax) = false, want true")
	}
}

func TestSyntaxErrorAs(t *testing.T) {
	wrapped := fmt.Errorf("wrapped: %w", SyntaxError{
		Line:   3,
		Source: "\"ok\"junk",
		Err:    strconv.ErrSyntax,
	})

	var gotValue SyntaxError
	if !errors.As(wrapped, &gotValue) {
		t.Fatalf("errors.As(wrapped, *SyntaxError) = false, want true")
	}
	if gotValue.Line != 3 || gotValue.Source != "\"ok\"junk" {
		t.Fatalf("got value = %#v", gotValue)
	}
	if !errors.Is(gotValue.Err, strconv.ErrSyntax) {
		t.Fatalf("errors.Is(gotValue.Err, strconv.ErrSyntax) = false, want true")
	}
}
