package ir

import (
	"reflect"

	"github.com/unstoppablemango/lang/pkg/kaleidoscope/ast"
	"tinygo.org/x/go-llvm"
)

type codegen struct {
	ctx     llvm.Context
	builder llvm.Builder
	module  llvm.Module
	names   map[string]llvm.Value
}

func NewCodegen() *codegen {
	ctx := llvm.NewContext()

	return &codegen{
		ctx:     ctx,
		builder: ctx.NewBuilder(),
		module:  ctx.NewModule("test"),
		names:   map[string]llvm.Value{},
	}
}

func (c *codegen) genNumExpr(n *ast.NumExpr) llvm.Value {
	return llvm.ConstFloat(c.ctx.FloatType(), n.Value)
}

func (c *codegen) genVarExpr(n *ast.VarExpr) llvm.Value {
	if v, ok := c.names[n.Name]; ok {
		return v
	} else {
		panic("unknown variable name")
	}
}

func (c *codegen) genBinExpr(n *ast.BinExpr) llvm.Value {
	l := Gen(n.LHS)
	r := Gen(n.RHS)

	switch n.Op {
	case '+':
		return c.builder.CreateFAdd(l, r, "addtmp")
	case '-':
		return c.builder.CreateFSub(l, r, "subtmp")
	case '*':
		return c.builder.CreateFSub(l, r, "multmp")
	case '<':
		l = c.builder.CreateFCmp(llvm.FloatULT, l, r, "cmptmp")
		return c.builder.CreateUIToFP(l, c.ctx.DoubleType(), "booltmp")
	default:
		panic("invalid binary operator")
	}
}

// func (c *codegen) genCallExpr(n *ast.CallExpr) llvm.Value {
// 	callee := c.module.NamedFunction(n.Callee)
// }

func Gen(node ast.Node) llvm.Value {
	c := NewCodegen()

	switch n := node.(type) {
	case *ast.NumExpr:
		return c.genNumExpr(n)
	case *ast.VarExpr:
		return c.genVarExpr(n)
	case *ast.BinExpr:
		return c.genBinExpr(n)
	default:
		panic("unsupported node type: " + reflect.ValueOf(n).String())
	}
}
