package scanner_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unstoppablemango/lang/pkg/scanner"
	"github.com/unstoppablemango/lang/pkg/token"
)

var _ = Describe("Scanner", func() {
	DescribeTable("should scan",
		func(input string, pos int, toks []token.Token, lits []string) {
			buf := bytes.NewBufferString(input)
			s := scanner.NewScanner(buf)

			var (
				lastPos  token.Pos
				tokens   []token.Token
				literals []string
			)

			for pos, tok, lit := s.Scan(); tok != token.EOF; {
				lastPos = pos
				tokens = append(tokens, tok)

				if lit != "" {
					literals = append(literals, lit)
				}
			}

			Expect(tokens).To(Equal(toks))
			Expect(literals).To(Equal(lits))
			Expect(lastPos).To(Equal(token.Pos(pos)))
		},
		Entry(nil, "a", 1, []token.Token{token.IDENT}, []string{"a"}),
	)
})
