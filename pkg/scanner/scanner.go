package scanner

import (
	"bufio"
	"io"

	"github.com/unstoppablemango/lang/pkg/token"
)

type Scanner interface {
	Scan() (token.Pos, token.Token, string)
}

type scanner struct {
	s      *bufio.Scanner
	offset int
	err    error
}

func (s *scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	return
}

func NewScanner(r io.Reader) Scanner {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	return &scanner{s: s}
}
