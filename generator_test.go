package bechars

import (
	"testing"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"[abc]", "abc"},
		// First characters
		{"[]", ""},
		{"[-]", "-"},
		{"[]abc]", "]abc"},
		// Character
		{"[\u0e010-2]", "ก012"},
		// Not
		{"[^:cntrl::punct:]", " 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"},
		{"[^-:cntrl::digit:]", " !\"#$%&'()*+,./:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
		{"[^]:cntrl::digit:]", " !\"#$%&'()*+,-./:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\^_`abcdefghijklmnopqrstuvwxyz{|}~"},
		// Classes
		{"[:alnum:]", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"},
		{"[:alpha:]", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"},
		{"[:digit:]", "0123456789"},
		{"[:space:]", " \t\n\r\f\v"},
		{"[:blank:]", " \t"},
		{"[:word:]", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"},
		{"[:cntrl:]", "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\u007f"},
		{"[:lower:]", "abcdefghijklmnopqrstuvwxyz"},
		{"[:upper:]", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"[:digit::upper:]", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{"[:print:]", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
		{"[:graph:]", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
		{"[:punct:]", "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]"},
		{"[:xdigit:]", "abcdefABCDEF0123456789"},
		{"[:digit::punct::upper:]", "0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		// Ranges
		{"[a-d]", "abcd"},
		{"[\x20-\x7E]", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
		// Swapped
		{"[\x7E-\x20]", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
	}

	for _, tt := range tests {
		p, err := New()
		if err != nil {
			t.Fatalf("New failed with %q", err)
		}
		actual, err := p.Generate(tt.in)
		if err != nil {
			t.Errorf("%s => failed with %q", tt.in, err)
		}
		if actual != tt.expected {
			t.Errorf("%s => expected %q, got %q", tt.in, tt.expected, actual)
		}
	}
}

func TestOptions(t *testing.T) {
	tests := []struct {
		in       string
		opts     []Option
		expected string
	}{
		{"[a-z]", []Option{MinRune('a'), MaxRune('f')}, "abcdef"},
		{"[^:cntrl::punct:]", []Option{MinRune('a'), MaxRune('z')}, "abcdefghijklmnopqrstuvwxyz"},
		{"[:upper:]", []Option{MinRune('a'), MaxRune('z')}, ""},
	}

	for _, tt := range tests {
		p, err := New(tt.opts...)
		if err != nil {
			t.Fatalf("New failed with %q", err)
		}
		actual, err := p.Generate(tt.in)
		if err != nil {
			t.Errorf("%s => failed with %q", tt.in, err)
		}
		if actual != tt.expected {
			t.Errorf("%s => expected %q, got %q", tt.in, tt.expected, actual)
		}
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"abc", "missing opening '['"},
		{"[a-]", "invalid range"},
		{"[aa-]", "invalid range"},
		{"[a-b-]", "parse error, unexpected '-'"},
		{"[:digit]", "parse error, expecting ':'"},
		{"[ab", "parse error, unexpected EOF"},
		{"[:blah:]", "invalid class ':blah:'"},
	}

	for _, tt := range tests {
		p, err := New()
		if err != nil {
			t.Fatalf("New failed with %q", err)
		}
		_, err = p.Generate(tt.in)
		if err == nil {
			t.Errorf("%s => should have failed", tt.in)
		}
		if err.Error() != tt.expected {
			t.Errorf("%s => expected %q, got %q", tt.in, tt.expected, err.Error())
		}
	}
}
