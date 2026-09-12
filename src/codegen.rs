//! Syntax tree to LLVM IR text.

use std::collections::HashMap;

use inkwell::builder::Builder;
use inkwell::context::Context;
use inkwell::module::{Linkage, Module};
use inkwell::types::ArrayType;
use inkwell::values::{FunctionValue, GlobalValue};
use inkwell::AddressSpace;

use crate::ast::{Expr, Stmt};
use crate::diag::Diagnostic;

/// A string constant in read-only data, kept with its type because a GEP needs
/// both.
#[derive(Clone, Copy)]
struct StrConst<'ctx> {
    ty: ArrayType<'ctx>,
    global: GlobalValue<'ctx>,
}

struct Codegen<'ctx, 'src> {
    source: &'src str,
    context: &'ctx Context,
    module: Module<'ctx>,
    builder: Builder<'ctx>,
    puts: FunctionValue<'ctx>,
    env: HashMap<&'src str, StrConst<'ctx>>,
}

pub fn emit<'src>(source: &'src str, stmts: &[Stmt<'src>]) -> Result<String, Diagnostic> {
    let context = Context::create();
    let mut codegen = Codegen::new(source, &context);
    codegen.run(stmts)?;
    Ok(codegen.module.print_to_string().to_string())
}

impl<'ctx, 'src> Codegen<'ctx, 'src> {
    fn new(source: &'src str, context: &'ctx Context) -> Self {
        let module = context.create_module("unmang");
        let ptr_type = context.ptr_type(AddressSpace::default());
        let puts_type = context.i32_type().fn_type(&[ptr_type.into()], false);
        let puts = module.add_function("puts", puts_type, None);

        Codegen {
            source,
            context,
            module,
            builder: context.create_builder(),
            puts,
            env: HashMap::new(),
        }
    }

    fn run(&mut self, stmts: &[Stmt<'src>]) -> Result<(), Diagnostic> {
        let i32_type = self.context.i32_type();
        let main_fn = self
            .module
            .add_function("main", i32_type.fn_type(&[], false), None);
        let entry = self.context.append_basic_block(main_fn, "entry");
        self.builder.position_at_end(entry);

        for stmt in stmts {
            match stmt {
                Stmt::Let { name, value } => {
                    if self.env.contains_key(name) {
                        return Err(Diagnostic::at(
                            format!("`{name}` is already bound"),
                            self.source,
                            name,
                        ));
                    }
                    let konst = self.value(value)?;
                    self.env.insert(name, konst);
                }
                Stmt::Print(expr) => {
                    let konst = self.value(expr)?;
                    let zero = i32_type.const_int(0, false);
                    let ptr = unsafe {
                        self.builder
                            .build_gep(
                                konst.ty,
                                konst.global.as_pointer_value(),
                                &[zero, zero],
                                "strptr",
                            )
                            .unwrap()
                    };
                    self.builder
                        .build_call(self.puts, &[ptr.into()], "call")
                        .unwrap();
                }
            }
        }

        self.builder
            .build_return(Some(&i32_type.const_int(0, false)))
            .unwrap();

        Ok(())
    }

    fn value(&self, expr: &Expr<'src>) -> Result<StrConst<'ctx>, Diagnostic> {
        match expr {
            Expr::Str(text) => Ok(self.const_str(text)),
            Expr::Name(name) => {
                self.env.get(name).copied().ok_or_else(|| {
                    Diagnostic::at(format!("unknown name `{name}`"), self.source, name)
                })
            }
        }
    }

    fn const_str(&self, text: &str) -> StrConst<'ctx> {
        let mut bytes = text.as_bytes().to_vec();
        bytes.push(0);

        let ty = self.context.i8_type().array_type(bytes.len() as u32);
        let global = self.module.add_global(ty, None, ".str");
        global.set_linkage(Linkage::Private);
        global.set_unnamed_addr(true);
        global.set_constant(true);
        global.set_initializer(&self.context.const_string(&bytes, false));

        StrConst { ty, global }
    }
}
