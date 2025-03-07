package parser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unstoppablemango/lang/pkg/kaleidoscope/ast"
	"github.com/unstoppablemango/lang/pkg/kaleidoscope/parser"
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
				Proto: &ast.Proto{Name: "foo", Args: []string{"x", "y"}},
				Body: &ast.BinExpr{
					Op:  '+',
					LHS: &ast.VarExpr{Name: "x"},
					RHS: &ast.CallExpr{
						Callee: "foo",
						Args: []ast.Expr{
							&ast.VarExpr{Name: "y"},
							&ast.NumExpr{Value: 4},
						},
					},
				},
			}},
		}),
		Entry(nil, "def foo(x y) x+y y;", &ast.File{Decls: []ast.Node{
			&ast.Func{
				Proto: &ast.Proto{Name: "foo", Args: []string{"x", "y"}},
				Body: &ast.BinExpr{
					Op:  '+',
					LHS: &ast.VarExpr{Name: "x"},
					RHS: &ast.VarExpr{Name: "y"},
				},
			},
			&ast.Func{
				Proto: &ast.Proto{Name: "__anon_expr", Args: []string{}},
				Body:  &ast.VarExpr{Name: "y"},
			},
		}}),
		Entry(nil, "extern sin(a);", &ast.File{Decls: []ast.Node{
			&ast.Proto{Name: "sin", Args: []string{"a"}},
		}}),
	)
})
