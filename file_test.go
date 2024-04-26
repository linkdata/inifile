package inifile

import "testing"

func TestFile(t *testing.T) {
	inif := make(File)
	if x := inif.GetDefault("", "key", "missing"); x != "missing" {
		t.Errorf("%q", x)
	}
	inif.Set("", "key", "val1", 0)
	if x := inif.GetDefault("", "key", "missing"); x != "val1" {
		t.Errorf("%q", x)
	}
	if x, ok := inif.Get("", "key"); !ok || x != "val1" {
		t.Errorf("%q %v", x, ok)
	}
}
