//! The syntax tree.
//!
//! An identifier is stored as a slice of the source text, so it carries its own
//! position: `offset_of(source, name)` recovers the caret for a diagnostic.
//!
//! A binding here means no more than "a name refers to a string constant known
//! at compile time". It is deliberately not evidence about mutability, the
//! memory model, scoping, typing, or evaluation order, all of which are open
//! decisions in `docs/design/`. There is no assignment syntax, rebinding is
//! rejected rather than given a meaning, the value lives in read-only data and
//! is never allocated or copied, and with no blocks in the grammar the question
//! of lexical versus dynamic scope is not asked.

pub enum Stmt<'a> {
    Let { name: &'a str, value: Expr<'a> },
    Print(Expr<'a>),
}

pub enum Expr<'a> {
    Str(String),
    Name(&'a str),
}
