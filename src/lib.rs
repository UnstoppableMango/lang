//! The compiler, as a library, so tests can compile source text without
//! shelling out to the binary.

pub mod ast;
pub mod codegen;
pub mod diag;
pub mod parse;

pub use diag::Diagnostic;

/// Compile source text to LLVM IR text.
pub fn compile(source: &str) -> Result<String, Diagnostic> {
    let stmts = parse::program(source)?;
    codegen::emit(source, &stmts)
}
