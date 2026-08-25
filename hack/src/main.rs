// Throwaway hello-world compiler: parses `print "..."` and emits LLVM IR text.
// Carries no design commitment; see the repo workflow docs.

use std::env;
use std::fs;
use std::process::exit;

use inkwell::context::Context;
use inkwell::module::Linkage;
use inkwell::AddressSpace;
use nom::{
    bytes::complete::{tag, take_until},
    sequence::delimited,
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
    let line = source.lines().next().unwrap_or_else(|| die("empty source file"));
    let (_, stmt) = parse_print(line).unwrap_or_else(|_| die("expected: print \"...\""));

    let context = Context::create();
    let module = context.create_module("unmang");
    let builder = context.create_builder();

    let i8_type = context.i8_type();
    let i32_type = context.i32_type();

    let mut bytes = stmt.text.into_bytes();
    bytes.push(0);
    let str_type = i8_type.array_type(bytes.len() as u32);
    let global = module.add_global(str_type, None, ".str");
    global.set_linkage(Linkage::Private);
    global.set_unnamed_addr(true);
    global.set_constant(true);
    global.set_initializer(&context.const_string(&bytes, false));

    let ptr_type = context.ptr_type(AddressSpace::default());
    let puts_type = i32_type.fn_type(&[ptr_type.into()], false);
    let puts_fn = module.add_function("puts", puts_type, None);

    let main_type = i32_type.fn_type(&[], false);
    let main_fn = module.add_function("main", main_type, None);
    let entry = context.append_basic_block(main_fn, "entry");
    builder.position_at_end(entry);

    let zero = i32_type.const_int(0, false);
    let str_ptr = unsafe {
        builder
            .build_gep(str_type, global.as_pointer_value(), &[zero, zero], "strptr")
            .unwrap()
    };

    builder
        .build_call(puts_fn, &[str_ptr.into()], "call")
        .unwrap();
    builder
        .build_return(Some(&i32_type.const_int(0, false)))
        .unwrap();

    print!("{}", module.print_to_string().to_string());
}
