// Minimal compiler front end: parses a sequence of `print "..."` statements and
// emits LLVM IR text.

use std::env;
use std::fs;
use std::process::exit;

use inkwell::context::Context;
use inkwell::module::Linkage;
use inkwell::AddressSpace;
use nom::{
    bytes::complete::{tag, take_until},
    character::complete::multispace0,
    combinator::eof,
    multi::many1,
    sequence::{delimited, preceded, terminated},
    IResult, Parser,
};

struct Print {
    text: String,
}

fn parse_print(input: &str) -> IResult<&str, Print> {
    let (input, _) = tag("print ")(input)?;
    let (input, text) = delimited(tag("\""), take_until("\""), tag("\"")).parse(input)?;
    Ok((
        input,
        Print {
            text: text.to_string(),
        },
    ))
}

// Whitespace between statements is skipped rather than counted, so a newline is
// no more of a separator than a space. Nothing here presumes a line-oriented
// grammar; that question belongs to a design doc, not to this parser.
fn parse_program(input: &str) -> IResult<&str, Vec<Print>> {
    terminated(
        many1(preceded(multispace0, parse_print)),
        preceded(multispace0, eof),
    )
    .parse(input)
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
    if source.trim().is_empty() {
        die("empty source file");
    }
    let (_, stmts) = parse_program(&source).unwrap_or_else(|_| die("expected: print \"...\""));

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
