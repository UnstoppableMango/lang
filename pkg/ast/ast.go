package ast

type Node interface{}

type Expr interface {
	Node
	exprNode()
}

type NumExpr struct {
	Value float64
}

func (*NumExpr) exprNode() {}

type VarExpr struct {
	Name string
}

func (*VarExpr) exprNode() {}

type BinExpr struct {
	Op       rune
	LHS, RHS Expr
}

func (*BinExpr) exprNode() {}

type CallExpr struct {
	Callee string
	Args   []Expr
}

func (*CallExpr) exprNode() {}

type Proto struct {
	Name string
	Args []string
}

type Func struct {
	Proto *Proto
	Body  Expr
}
