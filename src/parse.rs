//! Source text to syntax tree.

use nom::{
    branch::alt,
    bytes::complete::{tag, take_until},
    character::complete::{alpha1, alphanumeric1, char, multispace0, multispace1},
    combinator::{map, recognize, verify},
    error::{context, ContextError, ErrorKind, ParseError},
    multi::many0_count,
    sequence::{delimited, pair, preceded},
    IResult, Parser,
};

use crate::ast::{Expr, Stmt};
use crate::diag::Diagnostic;

/// Words a name may not be. The list is ad hoc and states no reserved-word
/// policy; it exists so `let let = "x"` does not parse.
const KEYWORDS: [&str; 2] = ["let", "print"];

/// What the parser wanted, and where.
///
/// `what` stays `None` until an enclosing `context` names it, so the deepest
/// position wins while the innermost description wins.
#[derive(Debug)]
struct Expected<'a> {
    input: &'a str,
    what: Option<&'static str>,
}

type PResult<'a, T> = IResult<&'a str, T, Expected<'a>>;

impl<'a> ParseError<&'a str> for Expected<'a> {
    fn from_error_kind(input: &'a str, _: ErrorKind) -> Self {
        Expected { input, what: None }
    }

    fn append(_: &'a str, _: ErrorKind, other: Self) -> Self {
        other
    }

    fn or(self, other: Self) -> Self {
        // Of two alternatives, the one that consumed the most describes the
        // failure best. A tie goes to the one tried first.
        if other.input.len() < self.input.len() {
            other
        } else {
            self
        }
    }
}

impl<'a> ContextError<&'a str> for Expected<'a> {
    fn add_context(_: &'a str, ctx: &'static str, other: Self) -> Self {
        Expected {
            input: other.input,
            what: other.what.or(Some(ctx)),
        }
    }
}

fn string_literal(input: &str) -> PResult<'_, &str> {
    context(
        "a string literal",
        delimited(tag("\""), take_until("\""), tag("\"")),
    )
    .parse(input)
}

fn identifier(input: &str) -> PResult<'_, &str> {
    context(
        "an identifier",
        verify(
            recognize(pair(
                alt((alpha1, tag("_"))),
                many0_count(alt((alphanumeric1, tag("_")))),
            )),
            |name: &str| !KEYWORDS.contains(&name),
        ),
    )
    .parse(input)
}

fn expr(input: &str) -> PResult<'_, Expr<'_>> {
    alt((
        map(string_literal, |text: &str| Expr::Str(text.to_string())),
        map(identifier, Expr::Name),
    ))
    .parse(input)
}

fn let_stmt(input: &str) -> PResult<'_, Stmt<'_>> {
    map(
        (
            tag("let"),
            multispace1,
            identifier,
            multispace0,
            char('='),
            multispace0,
            expr,
        ),
        |(_, _, name, _, _, _, value)| Stmt::Let { name, value },
    )
    .parse(input)
}

fn print_stmt(input: &str) -> PResult<'_, Stmt<'_>> {
    map(preceded((tag("print"), multispace1), expr), Stmt::Print).parse(input)
}

fn statement(input: &str) -> PResult<'_, Stmt<'_>> {
    context("a statement", alt((let_stmt, print_stmt))).parse(input)
}

// Whitespace between statements is skipped rather than counted, so a newline is
// no more of a separator than a space. Nothing here presumes a line-oriented
// grammar; that question belongs to a design doc, not to this parser.
pub fn program(source: &str) -> Result<Vec<Stmt<'_>>, Diagnostic> {
    let mut input = source;
    let mut stmts = Vec::new();

    loop {
        input = input.trim_start_matches([' ', '\t', '\r', '\n']);
        if input.is_empty() {
            break;
        }

        let (rest, stmt) = statement(input).map_err(|e| to_diagnostic(source, e))?;
        stmts.push(stmt);
        input = rest;
    }

    if stmts.is_empty() {
        return Err(Diagnostic::new("empty source file", 0));
    }

    Ok(stmts)
}

fn to_diagnostic(source: &str, err: nom::Err<Expected<'_>>) -> Diagnostic {
    match err {
        // Failure rather than Error means a parser committed before failing;
        // both describe the same position, so both render the same way.
        nom::Err::Error(e) | nom::Err::Failure(e) => Diagnostic::at(
            format!("expected {}", e.what.unwrap_or("a statement")),
            source,
            e.input,
        ),
        nom::Err::Incomplete(_) => unreachable!("every parser used here is complete"),
    }
}
