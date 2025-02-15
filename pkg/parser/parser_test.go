package parser_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/unstoppablemango/lang/pkg/parser"
)

var _ = Describe("Parser", func() {
	DescribeTable("should parse",
		func(input string) {
			p := parser.NewParser([]byte(input))

			p.Parse()
		},
		Entry(nil, "def test()"),
	)
})
