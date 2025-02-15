package scanner_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unstoppablemango/lang/pkg/scanner"
	"github.com/unstoppablemango/lang/pkg/token"
)

var _ = Describe("Scanner", func() {
	DescribeTable("should scan",
		func(input string, pos int, toks []token.Token, lits []string) {
			s := scanner.NewScanner([]byte(input))

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

				pos, tok, lit = s.Scan()
			}

			Expect(tokens).To(Equal(toks))
			Expect(literals).To(Equal(lits))
			Expect(lastPos).To(Equal(token.Pos(pos)))
		},
		Entry(nil, "a", 1, []token.Token{token.IDENT}, []string{"a"}),
		Entry(nil, "ab", 1, []token.Token{token.IDENT}, []string{"ab"}),
		Entry(nil, "abc", 1, []token.Token{token.IDENT}, []string{"abc"}),
		Entry(nil, "ab c", 4, []token.Token{token.IDENT, token.IDENT}, []string{"ab", "c"}),
		Entry(nil, "abc1", 1, []token.Token{token.IDENT}, []string{"abc1"}),
		Entry(nil, "abc 1", 5, []token.Token{token.IDENT, token.NUM}, []string{"abc", "1"}),
		Entry(nil, "1", 1, []token.Token{token.NUM}, []string{"1"}),
		Entry(nil, "12", 1, []token.Token{token.NUM}, []string{"12"}),
		Entry(nil, "12.3", 1, []token.Token{token.NUM}, []string{"12.3"}),
		Entry(nil, "12.3 # test", 1, []token.Token{token.NUM}, []string{"12.3"}),
	)
})
