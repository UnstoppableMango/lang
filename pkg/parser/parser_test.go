package parser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unstoppablemango/lang/pkg/ast"
	"github.com/unstoppablemango/lang/pkg/parser"
)

var _ = Describe("Parser", func() {
	DescribeTable("should parse",
		func(input string, expected *ast.File) {
			p := parser.NewParser([]byte(input))

			f := p.Parse()

			Expect(f).To(Equal(expected))
		},
		Entry(nil, "test", &ast.File{
			Decls: []ast.Node{&ast.Func{
				Proto: &ast.Proto{
					Name: "__anon_expr",
					Args: []string{},
				},
				Body: &ast.VarExpr{Name: "test"},
			}},
		}),
		Entry(nil, "test()", &ast.File{
			Decls: []ast.Node{&ast.Func{
				Proto: &ast.Proto{
					Name: "__anon_expr",
					Args: []string{},
				},
				Body: &ast.CallExpr{Callee: "test"},
			}},
		}),
		Entry(nil, "def test() test2", &ast.File{
			Decls: []ast.Node{&ast.Func{
				Proto: &ast.Proto{Name: "test"},
				Body:  &ast.VarExpr{Name: "test2"},
			}},
		}),
		Entry(nil, "def foo(x y) x+foo(y, 4.0)", &ast.File{
			Decls: []ast.Node{&ast.Func{
				Proto: &ast.Proto{Name: "test"},
				Body:  &ast.VarExpr{Name: "test2"},
			}},
		}),
		Entry(nil, "def foo(x y) x+y y;", &ast.File{
			Decls: []ast.Node{&ast.Func{
				Proto: &ast.Proto{Name: "test"},
				Body:  &ast.VarExpr{Name: "test2"},
			}},
		}),
		Entry(nil, "extern sin(a);", &ast.File{
			Decls: []ast.Node{&ast.Func{
				Proto: &ast.Proto{Name: "test"},
				Body:  &ast.VarExpr{Name: "test2"},
			}},
		}),
	)
})
