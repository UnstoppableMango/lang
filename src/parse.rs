//! Source text to syntax tree.

use nom::{
    branch::alt,
    bytes::complete::{is_not, tag},
    character::complete::{alpha1, alphanumeric1, char, digit1, multispace0, multispace1},
    combinator::{consumed, cut, map, map_res, recognize, value, verify},
    error::{context, ContextError, ErrorKind, FromExternalError, ParseError},
    multi::{fold_many0, many0, many0_count},
    sequence::{delimited, pair, preceded},
    IResult, Parser,
};

use crate::ast::{BinOp, Expr, Stmt};
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

impl<'a, E> FromExternalError<&'a str, E> for Expected<'a> {
    fn from_external_error(input: &'a str, _: ErrorKind, _: E) -> Self {
        Expected {
            input,
            what: Some("an integer that fits in 64 bits"),
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

/// A run of ordinary characters, or one escape sequence.
enum Fragment<'a> {
    Literal(&'a str),
    Escaped(char),
}

/// Everything up to the next quote or backslash. `is_not` accepts newlines, so
/// an unterminated string swallows the rest of the file and is reported at the
/// end of it.
fn literal_chunk(input: &str) -> PResult<'_, &str> {
    is_not("\"\\").parse(input)
}

/// Always fails, at the character after the backslash.
fn unknown_escape(input: &str) -> PResult<'_, char> {
    Err(nom::Err::Failure(Expected {
        input,
        what: Some("a known escape: \\n, \\t, \\\\ or \\\""),
    }))
}

fn escape(input: &str) -> PResult<'_, char> {
    let (rest, _) = char('\\').parse(input)?;
    if rest.is_empty() {
        return Err(nom::Err::Failure(Expected {
            input,
            what: Some("an escape character after `\\`"),
        }));
    }
    // Without the cut, an unknown escape would end the fragment loop quietly
    // and the failure would be blamed on the missing closing quote.
    cut(alt((
        value('\n', char('n')),
        value('\t', char('t')),
        value('\\', char('\\')),
        value('"', char('"')),
        unknown_escape,
    )))
    .parse(rest)
}

fn string_literal(input: &str) -> PResult<'_, String> {
    context(
        "a string literal",
        delimited(
            char('"'),
            fold_many0(
                alt((
                    map(literal_chunk, Fragment::Literal),
                    map(escape, Fragment::Escaped),
                )),
                String::new,
                |mut text, fragment| {
                    match fragment {
                        Fragment::Literal(chunk) => text.push_str(chunk),
                        Fragment::Escaped(c) => text.push(c),
                    }
                    text
                },
            ),
            cut(context("a closing `\"`", char('"'))),
        ),
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

fn int_literal(input: &str) -> PResult<'_, Expr<'_>> {
    map_res(digit1, |text: &str| text.parse().map(Expr::Int)).parse(input)
}

fn primary(input: &str) -> PResult<'_, Expr<'_>> {
    context(
        "an expression",
        alt((
            int_literal,
            map(string_literal, Expr::Str),
            map(identifier, Expr::Name),
            delimited(
                char('('),
                preceded(multispace0, expr),
                preceded(multispace0, cut(context("a closing `)`", char(')')))),
            ),
        )),
    )
    .parse(input)
}

/// One level of left-associative infix operators over `next`.
fn infix<'a>(
    input: &'a str,
    operator: impl Parser<&'a str, Output = BinOp, Error = Expected<'a>>,
    next: impl Fn(&'a str) -> PResult<'a, Expr<'a>>,
) -> PResult<'a, Expr<'a>> {
    let (input, first) = next(input)?;
    let (input, rest) = many0(pair(
        delimited(multispace0, consumed(operator), multispace0),
        &next,
    ))
    .parse(input)?;

    let folded = rest
        .into_iter()
        .fold(first, |lhs, ((at, op), rhs)| Expr::Binary {
            op,
            at,
            lhs: Box::new(lhs),
            rhs: Box::new(rhs),
        });

    Ok((input, folded))
}

fn multiplicative(input: &str) -> PResult<'_, Expr<'_>> {
    infix(
        input,
        alt((value(BinOp::Mul, char('*')), value(BinOp::Div, char('/')))),
        primary,
    )
}

// Whitespace around an operator is multispace0, so `print 1\n+ 2` is one
// statement. That follows from the parser counting no newlines anywhere, which
// is the same position the statement loop takes.
fn additive(input: &str) -> PResult<'_, Expr<'_>> {
    infix(
        input,
        alt((value(BinOp::Add, char('+')), value(BinOp::Sub, char('-')))),
        multiplicative,
    )
}

fn expr(input: &str) -> PResult<'_, Expr<'_>> {
    additive(input)
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
