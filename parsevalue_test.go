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
