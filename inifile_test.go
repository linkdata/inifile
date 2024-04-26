package inifile

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		dupKeysJoin rune
		wantInif    File
		wantErr     error
	}{
		{
			name: "tests/quotedcomment.ini",
			wantInif: File{
				"": Section{
					"key": "quoted\"#",
				},
			},
		},
		{
			name:    "tests/brokenquote.ini",
			wantErr: SyntaxError{Line: 1, Source: "\"broken"},
		},
		{
			name:    "tests/nonexistant.ini",
			wantErr: os.ErrNotExist,
		},
		{
			name:     "tests/empty.ini",
			wantInif: File{},
		},
		{
			name: "tests/nojoin.ini",
			wantInif: File{
				"": Section{
					"global_1": "1",
					"global_2": "2",
				},
				"sect_1": Section{
					"sect_1_1": "1_1",
				},
			},
		},
		{
			name:        "tests/join.ini",
			dupKeysJoin: ',',
			wantInif: File{
				"": Section{
					"key": "1,2,3",
				},
			},
		},
		{
			name:    "tests/helloworld.ini",
			wantErr: SyntaxError{Line: 1, Source: "hello world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInif, err := Load(tt.name, tt.dupKeysJoin)
			if err == nil && tt.wantErr == nil {
				if !reflect.DeepEqual(gotInif, tt.wantInif) {
					t.Errorf("Load() = %v, want %v", gotInif, tt.wantInif)
				}
			} else {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if errors.Is(err, ErrSyntax) {
					errstr := err.Error()
					wantstr := tt.wantErr.Error()
					if errstr != wantstr {
						t.Errorf("Load() error = %q, wantErr %q", err, tt.wantErr)
						return
					}

				}

			}
		})
	}
}
