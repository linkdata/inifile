package inifile

import (
	"testing"
	"unicode/utf8"
)

func TestParseValueSingleQuotedAlwaysValidUTF8(t *testing.T) {
	input := "'" + string([]byte{0xff}) + "'"

	value, err := ParseValue(input)
	if err != nil {
		t.Fatalf("ParseValue() error = %v, want nil", err)
	}
	if !utf8.ValidString(value) {
		t.Fatalf("ParseValue() produced invalid UTF-8: % x", []byte(value))
	}
	if value != "\uFFFD" {
		t.Fatalf("ParseValue() value = %q, want %q", value, "\uFFFD")
	}
}

func TestParseValueUnquotedAlwaysValidUTF8(t *testing.T) {
	input := string([]byte{0xff, 'a'})

	value, err := ParseValue(input)
	if err != nil {
		t.Fatalf("ParseValue() error = %v, want nil", err)
	}
	if !utf8.ValidString(value) {
		t.Fatalf("ParseValue() produced invalid UTF-8: % x", []byte(value))
	}
	if value != "\uFFFDa" {
		t.Fatalf("ParseValue() value = %q, want %q", value, "\uFFFDa")
	}
}

func TestParseValueUnquotedCommentAlwaysValidUTF8(t *testing.T) {
	input := string([]byte{0xff}) + "#comment"

	value, err := ParseValue(input)
	if err != nil {
		t.Fatalf("ParseValue() error = %v, want nil", err)
	}
	if !utf8.ValidString(value) {
		t.Fatalf("ParseValue() produced invalid UTF-8: % x", []byte(value))
	}
	if value != "\uFFFD" {
		t.Fatalf("ParseValue() value = %q, want %q", value, "\uFFFD")
	}
}
