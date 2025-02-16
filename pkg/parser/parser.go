package parser

import (
	"fmt"
	"strconv"

	"github.com/unstoppablemango/lang/pkg/ast"
	"github.com/unstoppablemango/lang/pkg/scanner"
	"github.com/unstoppablemango/lang/pkg/token"
)

type Parser interface {
	Parse() ast.Node
}

var binOpPrec = map[token.Token]int{
	token.LT:  10, // <
	token.ADD: 20, // +
	token.SUB: 30, // -
	token.MUL: 40, // *
}

func tokPrec(tok token.Token) int {
	if prec, ok := binOpPrec[tok]; ok {
		return prec
	} else {
		return -1
	}
}

type parser struct {
	s scanner.Scanner

	pos token.Pos
	tok token.Token
	lit string
}

func NewParser(src []byte) Parser {
	s := scanner.NewScanner(src)
	return &parser{s: s}
}

func (p *parser) next() {
	p.pos, p.tok, p.lit = p.s.Scan()
}

func (p *parser) pexpr() ast.Expr {
	if lhs := p.pprimary(); lhs != nil {
		return p.pbinopRhs(0, lhs)
	} else {
		return nil
	}
}

func (p *parser) pnum() *ast.NumExpr {
	num, err := strconv.ParseFloat(p.lit, 64)
	if err != nil {
		panic(err)
	}

	p.next()
	return &ast.NumExpr{Value: num}
}

func (p *parser) pparen() ast.Expr {
	p.next() // eat '('
	e := p.pexpr()
	if e == nil {
		return nil
	}

	if p.tok != token.RPAREN {
		panic("expected ')' found '" + p.tok.String() + "'")
	}

	p.next() // eat ')'
	return e
}

func (p *parser) pident() ast.Expr {
	name := p.lit
	p.next()

	if p.tok != token.LPAREN {
		return &ast.VarExpr{Name: name}
	}

	p.next() // eat '('
	var args []ast.Expr
	if p.tok != token.RPAREN {
		for {
			if arg := p.pexpr(); arg != nil {
				args = append(args, arg)
			} else {
				return nil
			}

			if p.tok == token.RPAREN {
				break
			}

			if p.tok != token.COMMA {
				panic("expected ')' or ',' in argument list, found '" + p.tok.String() + "' (" + p.lit + ")")
			}
			p.next()
		}
	}

	p.next() // eat ')'
	return &ast.CallExpr{
		Callee: name,
		Args:   args,
	}
}

func (p *parser) pprimary() ast.Expr {
	switch p.tok {
	case token.IDENT:
		return p.pident()
	case token.NUM:
		return p.pnum()
	case token.LPAREN:
		return p.pparen()
	default:
		panic("unuspported state: " + p.String())
	}
}

func (p *parser) pbinopRhs(exprPrec int, lhs ast.Expr) ast.Expr {
	for {
		prec := tokPrec(p.tok)
		if prec < exprPrec {
			return lhs
		}

		op := rune(p.tok.String()[0])
		p.next()

		rhs := p.pprimary()
		if rhs == nil {
			return nil
		}

		nextPrec := tokPrec(p.tok)
		if prec < nextPrec {
			rhs = p.pbinopRhs(prec+1, rhs)
			if rhs == nil {
				return nil
			}
		}

		lhs = &ast.BinExpr{
			Op:  op,
			LHS: lhs,
			RHS: rhs,
		}
	}
}

func (p *parser) pproto() *ast.Proto {
	if p.tok != token.IDENT {
		panic("expected function name in prototype")
	}

	name := p.lit
	p.next()

	if p.tok != token.LPAREN {
		panic("expected '(' in prototype, found '" + p.tok.String() + "'")
	}
	p.next() // eat '('

	var args []string
	for ; p.tok == token.IDENT; p.next() {
		args = append(args, p.lit)
	}

	if p.tok != token.RPAREN {
		panic("expected ')' in prototype")
	}

	p.next() // eat ')'
	return &ast.Proto{Name: name, Args: args}
}

func (p *parser) pdef() *ast.Func {
	p.next() // eat 'def'
	proto := p.pproto()
	if proto == nil {
		return nil
	}

	if e := p.pexpr(); e != nil {
		return &ast.Func{
			Proto: proto,
			Body:  e,
		}
	} else {
		return nil
	}
}

func (p *parser) ptopexpr() *ast.Func {
	e := p.pexpr()
	if e == nil {
		return nil
	}

	proto := &ast.Proto{
		Name: "__anon_expr",
		Args: []string{},
	}

	return &ast.Func{
		Proto: proto,
		Body:  e,
	}
}

func (p *parser) pextern() *ast.Proto {
	p.next() // eat 'extern'
	return p.pproto()
}

func (p *parser) Parse() ast.Node {
	f := &ast.File{}
	p.next()

	for {
		switch p.tok {
		case token.EOF:
			return f
		case token.SEMI:
			p.next()
		case token.DEF:
			f.Decls = append(f.Decls, p.pdef())
		case token.EXTERN:
			f.Decls = append(f.Decls, p.pextern())
		default:
			f.Decls = append(f.Decls, p.ptopexpr())
		}
	}
}

func (p *parser) String() string {
	return fmt.Sprintf("pos(%d) tok(%s) lit(%s)", p.pos, p.tok, p.lit)
}
