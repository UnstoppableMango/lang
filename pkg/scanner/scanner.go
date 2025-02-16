package scanner

import (
	"unicode"

	"github.com/unstoppablemango/lang/pkg/token"
)

type Scanner interface {
	Scan() (token.Pos, token.Token, string)
}

type scanner struct {
	src    []byte
	last   rune
	offset int
}

func (s *scanner) next() {
	if s.offset >= len(s.src) {
		s.last = -1
		return
	}

	r := s.src[s.offset]
	s.offset++
	s.last = rune(r)
}

func (s *scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	for unicode.IsSpace(s.last) {
		s.next()
	}

	pos = token.Pos(s.offset)

	switch s.last {
	case '(':
		tok = token.LPAREN
		s.next()
		return
	case ')':
		tok = token.RPAREN
		s.next()
		return
	case ',':
		tok = token.COMMA
		s.next()
		return
	}

	if unicode.IsLetter(s.last) {
		lit = lit + string(s.last)
		s.next()

		for alphanumeric(s.last) {
			lit = lit + string(s.last)
			s.next()
		}

		tok = token.Lookup(lit)
		return
	}

	if unicode.IsDigit(s.last) || s.last == '.' {
		for {
			lit = lit + string(s.last)
			s.next()

			if !unicode.IsDigit(s.last) && s.last != '.' {
				break
			}
		}

		tok = token.NUM
		return
	}

	if s.last == '#' {
		s.next()
		for s.last != -1 && s.last != '\n' && s.last != '\r' {
			s.next()
		}
		if s.last != -1 {
			return s.Scan()
		}
	}

	if s.last == -1 {
		tok = token.EOF
		return
	}

	return
}

func NewScanner(src []byte) Scanner {
	return &scanner{src, ' ', 0}
}

func alphanumeric(r rune) bool {
	return unicode.In(r, unicode.Letter, unicode.Digit)
}
