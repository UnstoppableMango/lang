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
	offset int
}

func (s *scanner) next() rune {
	if s.offset >= len(s.src) {
		return -1
	}

	r := s.src[s.offset]
	s.offset++

	return rune(r)
}

func (s *scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	last := ' '
	for unicode.IsSpace(last) {
		last = s.next()
	}

	pos = token.Pos(s.offset)

	if unicode.IsLetter(last) {
		lit = lit + string(last)
		last = s.next()

		for alphanumeric(last) {
			lit = lit + string(last)
			last = s.next()
		}

		tok = token.Lookup(lit)

		return
	}

	if unicode.IsDigit(last) || last == '.' {
		for {
			lit = lit + string(last)
			last = s.next()

			if !unicode.IsDigit(last) && last != '.' {
				break
			}
		}

		tok = token.NUM
		return
	}

	if last == '#' {
		last = s.next()
		for last != -1 && last != '\n' && last != '\r' {
			last = s.next()
		}
		if last != -1 {
			return s.Scan()
		}
	}

	if last == -1 {
		tok = token.EOF
		return
	}

	return
}

func NewScanner(src []byte) Scanner {
	return &scanner{src, 0}
}

func alphanumeric(r rune) bool {
	return unicode.In(r, unicode.Letter, unicode.Digit)
}
