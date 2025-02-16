package token

import "strconv"

type Token int

const (
	ILLEGAL Token = iota

	EOF    // EOF
	DEF    // def
	EXTERN // extern
	IDENT  // IDENT
	NUM    // NUM

	LPAREN // (
	RPAREN // )
	COMMA  // ,
)

var tokens = [...]string{
	EOF:    "EOF",
	DEF:    "def",
	EXTERN: "extern",
	IDENT:  "IDENT",
	NUM:    "NUM",
}

var keywords = map[string]Token{
	tokens[DEF]:    DEF,
	tokens[EXTERN]: EXTERN,
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
	if k, ok := keywords[ident]; ok {
		return k
	} else {
		return IDENT
	}
}
