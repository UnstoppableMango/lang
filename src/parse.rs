//! Source text to syntax tree.

use nom::{
    bytes::complete::{tag, take_until},
    combinator::map,
    error::{context, ContextError, ErrorKind, ParseError},
    sequence::{delimited, preceded},
    IResult, Parser,
};

use crate::ast::Print;
use crate::diag::Diagnostic;

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

fn statement(input: &str) -> PResult<'_, Print> {
    context(
        "a statement",
        map(preceded(tag("print "), string_literal), |text: &str| Print {
            text: text.to_string(),
        }),
    )
    .parse(input)
}

// Whitespace between statements is skipped rather than counted, so a newline is
// no more of a separator than a space. Nothing here presumes a line-oriented
// grammar; that question belongs to a design doc, not to this parser.
pub fn program(source: &str) -> Result<Vec<Print>, Diagnostic> {
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
