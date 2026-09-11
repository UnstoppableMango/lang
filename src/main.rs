// Minimal compiler front end: parses a sequence of `print "..."` statements and
// emits LLVM IR text.

mod diag;

use std::env;
use std::fs;
use std::process::exit;

use inkwell::context::Context;
use inkwell::module::Linkage;
use inkwell::AddressSpace;
use nom::{
    bytes::complete::{tag, take_until},
    combinator::map,
    error::{context, ContextError, ErrorKind, ParseError},
    sequence::{delimited, preceded},
    IResult, Parser,
};

use diag::Diagnostic;

struct Print {
    text: String,
}

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
        map(preceded(tag("print "), string_literal), |text: &str| {
            Print {
                text: text.to_string(),
            }
        }),
    )
    .parse(input)
}

// Whitespace between statements is skipped rather than counted, so a newline is
// no more of a separator than a space. Nothing here presumes a line-oriented
// grammar; that question belongs to a design doc, not to this parser.
fn program(source: &str) -> Result<Vec<Print>, Diagnostic> {
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

fn die(msg: &str) -> ! {
    eprintln!("unmangc: {msg}");
    exit(1);
}

fn main() {
    let path = match env::args().nth(1) {
        Some(p) => p,
        None => die("usage: unmangc <file>"),
    };

    let source = fs::read_to_string(&path).unwrap_or_else(|_| die("cannot open source file"));
    let stmts = program(&source).unwrap_or_else(|d| {
        eprint!("{}", d.render(&path, &source));
        exit(1);
    });

    let context = Context::create();
    let module = context.create_module("unmang");
    let builder = context.create_builder();

    let i8_type = context.i8_type();
    let i32_type = context.i32_type();
    let ptr_type = context.ptr_type(AddressSpace::default());

    let puts_type = i32_type.fn_type(&[ptr_type.into()], false);
    let puts_fn = module.add_function("puts", puts_type, None);

    let main_type = i32_type.fn_type(&[], false);
    let main_fn = module.add_function("main", main_type, None);
    let entry = context.append_basic_block(main_fn, "entry");
    builder.position_at_end(entry);

    let zero = i32_type.const_int(0, false);

    for stmt in stmts {
        let mut bytes = stmt.text.into_bytes();
        bytes.push(0);
        let str_type = i8_type.array_type(bytes.len() as u32);
        let global = module.add_global(str_type, None, ".str");
        global.set_linkage(Linkage::Private);
        global.set_unnamed_addr(true);
        global.set_constant(true);
        global.set_initializer(&context.const_string(&bytes, false));

        let str_ptr = unsafe {
            builder
                .build_gep(str_type, global.as_pointer_value(), &[zero, zero], "strptr")
                .unwrap()
        };

        builder
            .build_call(puts_fn, &[str_ptr.into()], "call")
            .unwrap();
    }

    builder
        .build_return(Some(&i32_type.const_int(0, false)))
        .unwrap();

    print!("{}", module.print_to_string().to_string());
}
