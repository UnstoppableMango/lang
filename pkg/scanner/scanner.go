package scanner

import (
	"io"

	"github.com/unstoppablemango/lang/pkg/token"
)

type Scanner interface {
	Scan() (token.Pos, token.Token, string)
}

type scanner struct {
	r      io.Reader
	offset int
}

func (s *scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	return
}

func NewScanner(r io.Reader) Scanner {
	return &scanner{r: r}
}
