package brechars

import (
	"fmt"
	"strings"

	"src.userspace.com.au/felix/lexer"
)

const (
	lower   = "abcdefghijklmnopqrstuvwxyz"
	upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numeric = "0123456789"
	space   = " \t\n\r\f\v"
	punct   = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]"
)

// Generator will generate strings of characters
// that match the provided POSIX bracket expression.
type Generator struct {
	l       *lexer.Lexer
	tok     *lexer.Token
	maxRune *rune
	minRune *rune
}

// Option functions allows configuration of the generator.
type Option func(*Generator) error

// New creates a new generator.
func New(opts ...Option) (*Generator, error) {
	out := new(Generator)
	for _, o := range opts {
		if err := o(out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// MaxRune sets the maximum rune for any generated sequences.
func MaxRune(r rune) Option {
	return func(g *Generator) error {
		g.maxRune = &r
		return nil
	}
}

// MinRune sets the minimum rune for any generated sequences.
func MinRune(r rune) Option {
	return func(g *Generator) error {
		g.minRune = &r
		return nil
	}
}

func ensureRangeLimits(g *Generator) {
	if g.maxRune == nil {
		maxRune := '\u007F'
		g.maxRune = &maxRune
	}
	if g.minRune == nil {
		minRune := '\u0000'
		g.minRune = &minRune
	}
}

func (g *Generator) next() bool {
	var ok bool
	g.tok, ok = g.l.NextToken()
	return g.tok != nil && !ok
}

// Generate will return the string from the POSIX bracket expression.
func (g *Generator) Generate(be string) (string, error) {
	ensureRangeLimits(g)
	g.l = lexer.New(be, startState)
	g.l.Start()

	g.next()
	if g.tok.Type != tBREStart {
		return "", fmt.Errorf("missing opening '['")
	}
	return g.buildSequence()
}

func (g *Generator) buildSequence() (string, error) {
	var out strings.Builder
	for g.next() {
		//fmt.Println(g.tok.Value)
		switch g.tok.Type {
		case tCharacter:
			out.WriteString(g.tok.Value)
		case tClass:
			s, err := g.getClass(g.tok.Value)
			if err != nil {
				return "", err
			}
			out.WriteString(g.filter(s, ""))
		case tRangeStart:
			start := g.tok.Value
			if g.next(); g.tok.Type != tRangeDash {
				// Impossible situ?
				return "", fmt.Errorf("invalid range")
			}
			if g.next(); g.tok.Type != tRangeEnd {
				return "", fmt.Errorf("invalid range")
			}
			end := g.tok.Value
			s := g.getRange([]rune(start)[0], []rune(end)[0])
			out.WriteString(g.filter(s, ""))
		case tBREStart, tBREEnd:
			// No op
		case tNot:
			nots, err := g.buildSequence()
			if err != nil {
				return "", err
			}
			s := g.getRange(*g.minRune, *g.maxRune)
			out.WriteString(g.filter(s, nots))
		case lexer.ErrorToken:
			return "", fmt.Errorf("%s", g.tok.Value)
		default:
			panic("invalid token")
		}
	}
	return out.String(), nil
}

func (g Generator) getClass(c string) (string, error) {
	var out string
	switch c {
	case ":alnum:":
		out = numeric + upper + lower
	case ":cntrl:":
		out = g.getRange('\u0000', '\u001F') + "\u007F"
	case ":lower:":
		out = lower
	case ":space:":
		out = space
	case ":alpha:":
		out = upper + lower
	case ":digit:":
		out = numeric
	case ":print:":
		fallthrough
	case ":graph:":
		c, err := g.getClass(":cntrl:")
		if err != nil {
			return "", err
		}
		out = g.filter(g.getRange(*g.minRune, *g.maxRune), c)

	case ":upper:":
		out = upper
	case ":blank:":
		out = " \t"
	case ":word:":
		c, err := g.getClass(":alnum:")
		if err != nil {
			return "", err
		}
		out = c + "_"
	case ":punct:":
		out = punct
	case ":xdigit:":
		out = "abcdefABCDEF" + numeric
	default:
		return "", fmt.Errorf("invalid class '%s'", c)
	}
	return out, nil
}

func (g *Generator) getRange(start, end rune) string {
	// Swap?
	if start > end {
		tmp := start
		start = end
		end = tmp
	}
	var out strings.Builder
	for i := start; i <= end; i++ {
		out.WriteRune(i)
	}
	return out.String()
}

func (g *Generator) filter(in, exclude string) string {
	var out strings.Builder
	for _, r := range in {
		if r < *g.minRune {
			continue
		}
		if r > *g.maxRune {
			continue
		}
		if strings.ContainsRune(exclude, r) {
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}
