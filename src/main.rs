// The `unmangc` binary: read a file, compile it, write LLVM IR to stdout.

use std::env;
use std::fs;
use std::process::exit;

use unmangc::compile;

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

    match compile(&source) {
        Ok(ir) => print!("{ir}"),
        Err(d) => {
            eprint!("{}", d.render(&path, &source));
            exit(1);
        }
    }
}
