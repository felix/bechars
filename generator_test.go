package bechars

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		in       string
		expected string
		options  []Option
	}{
		{
			in:       "[abc]",
			expected: "abc",
		},
		// First characters
		{in: "[]", expected: ""},
		{in: "[-]", expected: "-"},
		{in: "[]abc]", expected: "]abc"},
		// Character
		{in: "[\u0e010-2]",
			expected: "ก012",
		},
		// Unicode
		{in: "[ก0-2]",
			expected: "ก012",
		},
		// Not
		{
			in:       "[^:cntrl::punct:]",
			expected: " 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		},
		{
			in:       "[^-:cntrl::digit:]",
			expected: " !\"#$%&'()*+,./:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		{
			in:       "[^]:cntrl::digit:]",
			expected: " !\"#$%&'()*+,-./:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		// Classes
		{
			in:       "[:alnum:]",
			expected: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		},
		{
			in:       "[:alpha:]",
			expected: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		},
		{
			in:       "[:digit:]",
			expected: "0123456789",
		},
		{
			in:       "[:space:]",
			expected: " \t\n\r\f\v",
		},
		{
			in:       "[:blank:]",
			expected: " \t",
		},
		{
			in:       "[:word:]",
			expected: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_",
		},
		{
			in:       "[:cntrl:]",
			expected: "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\u007f",
		},
		{
			in:       "[:lower:]",
			expected: "abcdefghijklmnopqrstuvwxyz",
		},
		{
			in:       "[:upper:]",
			expected: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			in:       "[:digit::upper:]",
			expected: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			in:       "[:print:]",
			expected: " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		{
			in:       "[:graph:]",
			expected: " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		{
			in:       "[:punct:]",
			expected: "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]",
		},
		{
			in:       "[:xdigit:]",
			expected: "abcdefABCDEF0123456789",
		},
		{
			in:       "[:digit::punct::upper:]",
			expected: "0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		// Ranges
		{
			in:       "[a-d]",
			expected: "abcd",
		},
		{
			in:       "[\x20-\x7E]",
			expected: " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		// Swapped
		{
			in:       "[\x7E-\x20]",
			expected: " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		},
		// MaxRune
		{
			in:       "[ก-ฎ]",
			expected: "กขฃคฅฆงจฉชซฌญฎ",
			options:  []Option{MaxRune('\uffff')},
		},
		// // Exclude graphic
		// {
		// 	in:       "[฿-ๆ็่้๊๋์]",
		// 	expected: "฿เแโใไๅๆ็่้๊๋์",
		// 	options:  []Option{MaxRune('\uffff'), OnlyGraphic(true)},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p, err := New(tt.options...)
			if err != nil {
				t.Fatal(err)
			}
			actual, err := p.Generate(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(actual, tt.expected); diff != "" {
				t.Error(diff)
			}
			if actual != tt.expected {
				t.Errorf("got %v, want %v", actual, tt.expected)
			}
		})
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
