package inifile

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
)

type unexpectedEOFReader struct {
	data string
	pos  int
}

func (r *unexpectedEOFReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.ErrUnexpectedEOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	if r.pos >= len(r.data) {
		return n, io.ErrUnexpectedEOF
	}
	return n, nil
}

func TestParseScannerError(t *testing.T) {
	_, err := Parse(strings.NewReader("k="+strings.Repeat("a", 70*1024)+"\n"), 0)
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("Parse() error = %v, want %v", err, bufio.ErrTooLong)
	}
	var got SyntaxError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, *SyntaxError) = false, want true")
	}
	if got.Line != 1 || got.Source != "" {
		t.Fatalf("Parse() syntax error = %#v, want line=1 source=\"\"", got)
	}
}

func TestParseScannerErrorLineNumberAfterValidLine(t *testing.T) {
	_, err := Parse(strings.NewReader("ok=1\nk="+strings.Repeat("a", 70*1024)+"\n"), 0)
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("Parse() error = %v, want %v", err, bufio.ErrTooLong)
	}
	var got SyntaxError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, *SyntaxError) = false, want true")
	}
	if got.Line != 2 || got.Source != "" {
		t.Fatalf("Parse() syntax error = %#v, want line=2 source=\"\"", got)
	}
}

func TestParseScannerUnexpectedEOFLineNumberAfterSingleLine(t *testing.T) {
	_, err := Parse(&unexpectedEOFReader{data: "k=v"}, 0)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("Parse() error = %v, want %v", err, io.ErrUnexpectedEOF)
	}
	var got SyntaxError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, *SyntaxError) = false, want true")
	}
	if got.Line != 1 || got.Source != "" {
		t.Fatalf("Parse() syntax error = %#v, want line=1 source=\"\"", got)
	}
}

func TestParseScannerUnexpectedEOFLineNumberAfterTwoLines(t *testing.T) {
	_, err := Parse(&unexpectedEOFReader{data: "a=1\nb=2"}, 0)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("Parse() error = %v, want %v", err, io.ErrUnexpectedEOF)
	}
	var got SyntaxError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, *SyntaxError) = false, want true")
	}
	if got.Line != 2 || got.Source != "" {
		t.Fatalf("Parse() syntax error = %#v, want line=2 source=\"\"", got)
	}
}

func TestParseNilReader(t *testing.T) {
	inif, err := Parse(nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if inif != nil {
		t.Fatal("expected nil inif")
	}
}

func TestParseUTF8BOM(t *testing.T) {
	inif, err := Parse(strings.NewReader("\uFEFF[s]\nk=v\n"), 0)
	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}
	got, ok := inif.Get("s", "k")
	if !ok || got != "v" {
		t.Fatalf("Parse() value = %q, %v; want %q, true", got, ok, "v")
	}
}

func TestParseValueWhitespaceRules(t *testing.T) {
	inif, err := Parse(strings.NewReader(`
u =   abc   
d = "  abc  "
s = '  abc  '
`), 0)
	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if got, _ := inif.Get("", "u"); got != "abc" {
		t.Fatalf("Parse() unquoted value = %q, want %q", got, "abc")
	}
	if got, _ := inif.Get("", "d"); got != "  abc  " {
		t.Fatalf("Parse() double-quoted value = %q, want %q", got, "  abc  ")
	}
	if got, _ := inif.Get("", "s"); got != "  abc  " {
		t.Fatalf("Parse() single-quoted value = %q, want %q", got, "  abc  ")
	}
}

func TestParseSectionTrailingComments(t *testing.T) {
	inif, err := Parse(strings.NewReader(`
[s1] # comment
k1=v1
[s2];comment
k2=v2
`), 0)
	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if got, ok := inif.Get("s1", "k1"); !ok || got != "v1" {
		t.Fatalf("Parse() value s1/k1 = %q, %v; want %q, true", got, ok, "v1")
	}
	if got, ok := inif.Get("s2", "k2"); !ok || got != "v2" {
		t.Fatalf("Parse() value s2/k2 = %q, %v; want %q, true", got, ok, "v2")
	}
}

func TestParseRejectsMalformedSectionHeader(t *testing.T) {
	_, err := Parse(strings.NewReader("[a][b]\nk=v\n"), 0)
	want := SyntaxError{Line: 1, Source: "[a][b]"}
	if !errors.Is(err, ErrSyntax) {
		t.Fatalf("errors.Is(err, ErrSyntax) = false, want true")
	}

	var got SyntaxError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, *SyntaxError) = false, want true")
	}
	if got.Line != want.Line || got.Source != want.Source {
		t.Fatalf("Parse() syntax error = %#v, want %#v", got, want)
	}
	if err.Error() != want.Error() {
		t.Fatalf("Parse() error string = %q, want %q", err.Error(), want.Error())
	}
}

func TestParseSyntaxErrorWrapsSourceError(t *testing.T) {
	_, err := Parse(strings.NewReader("k='broken\n"), 0)
	if err == nil {
		t.Fatalf("Parse() error = nil, want non-nil")
	}
	if !errors.Is(err, ErrSyntax) {
		t.Fatalf("errors.Is(err, ErrSyntax) = false, want true")
	}
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("errors.Is(err, strconv.ErrSyntax) = false, want true")
	}

	var got SyntaxError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(err, *SyntaxError) = false, want true")
	}
	if got.Line != 1 || got.Source != "'broken" {
		t.Fatalf("Parse() syntax error = %#v", got)
	}
	if !errors.Is(got.Err, strconv.ErrSyntax) {
		t.Fatalf("errors.Is(got.Err, strconv.ErrSyntax) = false, want true")
	}
}
