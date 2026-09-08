//! The syntax tree.
//!
//! An identifier or an operator is stored as a slice of the source text, so it
//! carries its own position: `offset_of(source, slice)` recovers the caret for
//! a diagnostic.
//!
//! A binding here means no more than "a name refers to a value known at compile
//! time". It is deliberately not evidence about mutability, the memory model,
//! scoping, typing, or evaluation order, all of which are open decisions in
//! `docs/design/`. There is no assignment syntax, rebinding is rejected rather
//! than given a meaning, a string lives in read-only data and is never
//! allocated or copied, and with no blocks in the grammar the question of
//! lexical versus dynamic scope is not asked.

use std::fmt;

pub enum Stmt<'a> {
    Let { name: &'a str, value: Expr<'a> },
    Print(Expr<'a>),
}

pub enum Expr<'a> {
    Str(String),
    Int(i64),
    Name(&'a str),
    Binary {
        op: BinOp,
        /// The operator's own slice, so a type error points at the operator.
        at: &'a str,
        lhs: Box<Expr<'a>>,
        rhs: Box<Expr<'a>>,
    },
}

#[derive(Clone, Copy)]
pub enum BinOp {
    Add,
    Sub,
    Mul,
    Div,
}

impl fmt::Display for BinOp {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(match self {
            BinOp::Add => "+",
            BinOp::Sub => "-",
            BinOp::Mul => "*",
            BinOp::Div => "/",
        })
    }
}
