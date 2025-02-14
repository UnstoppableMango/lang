package token

import "strconv"

type Token int

const (
	EOF Token = iota
	DEF
	EXTERN
	IDENT
	NUM
)

var tokens = [...]string{
	EOF:    "EOF",
	DEF:    "def",
	EXTERN: "extern",
	IDENT:  "IDENT",
	NUM:    "NUM",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

// Lookup maps an identifier to its keyword token or [IDENT] (if not a keyword).
func Lookup(ident string) Token {
	switch ident {
	case tokens[DEF]:
		return DEF
	case tokens[EXTERN]:
		return EXTERN
	default:
		return IDENT
	}
}
