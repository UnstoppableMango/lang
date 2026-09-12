//! Syntax tree to LLVM IR text.

use std::collections::HashMap;

use inkwell::builder::Builder;
use inkwell::context::Context;
use inkwell::module::{Linkage, Module};
use inkwell::types::ArrayType;
use inkwell::values::{FunctionValue, GlobalValue, IntValue};
use inkwell::AddressSpace;

use crate::ast::{BinOp, Expr, Stmt};
use crate::diag::Diagnostic;

/// A string constant in read-only data, kept with its type because a GEP needs
/// both.
#[derive(Clone, Copy)]
struct StrConst<'ctx> {
    ty: ArrayType<'ctx>,
    global: GlobalValue<'ctx>,
}

/// What an expression evaluates to. Strings print through `puts`, integers
/// through `printf`; there is no conversion between them in either direction.
#[derive(Clone, Copy)]
enum Value<'ctx> {
    Str(StrConst<'ctx>),
    Int(IntValue<'ctx>),
}

struct Codegen<'ctx, 'src> {
    source: &'src str,
    context: &'ctx Context,
    module: Module<'ctx>,
    builder: Builder<'ctx>,
    puts: FunctionValue<'ctx>,
    printf: FunctionValue<'ctx>,
    /// The `%lld` format string, created on first use so a program that prints
    /// no integers carries no unused constant.
    int_format: Option<StrConst<'ctx>>,
    env: HashMap<&'src str, Value<'ctx>>,
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
        let printf_type = context.i32_type().fn_type(&[ptr_type.into()], true);
        let printf = module.add_function("printf", printf_type, None);

        Codegen {
            source,
            context,
            module,
            builder: context.create_builder(),
            puts,
            printf,
            int_format: None,
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
                Stmt::Print(expr) => match self.value(expr)? {
                    Value::Str(konst) => {
                        let ptr = self.str_ptr(konst);
                        self.builder
                            .build_call(self.puts, &[ptr.into()], "call")
                            .unwrap();
                    }
                    Value::Int(int) => {
                        let format = self.int_format();
                        let ptr = self.str_ptr(format);
                        self.builder
                            .build_call(self.printf, &[ptr.into(), int.into()], "call")
                            .unwrap();
                    }
                },
            }
        }

        self.builder
            .build_return(Some(&i32_type.const_int(0, false)))
            .unwrap();

        Ok(())
    }

    fn value(&self, expr: &Expr<'src>) -> Result<Value<'ctx>, Diagnostic> {
        match expr {
            Expr::Str(text) => Ok(Value::Str(self.const_str(text))),
            Expr::Int(int) => Ok(Value::Int(
                self.context.i64_type().const_int(*int as u64, true),
            )),
            Expr::Name(name) => {
                self.env.get(name).copied().ok_or_else(|| {
                    Diagnostic::at(format!("unknown name `{name}`"), self.source, name)
                })
            }
            // These are the real LLVM instruction builders. LLVM folds
            // constant operands itself, and every operand the language can
            // currently express is a constant, so the emitted IR shows the
            // answer rather than an `add`.
            Expr::Binary { op, at, lhs, rhs } => {
                let lhs = self.int_operand(lhs, *op, at)?;
                let rhs = self.int_operand(rhs, *op, at)?;
                let result = match op {
                    BinOp::Add => self.builder.build_int_add(lhs, rhs, "add"),
                    BinOp::Sub => self.builder.build_int_sub(lhs, rhs, "sub"),
                    BinOp::Mul => self.builder.build_int_mul(lhs, rhs, "mul"),
                    BinOp::Div => {
                        self.check_division(lhs, rhs, at)?;
                        self.builder.build_int_signed_div(lhs, rhs, "div")
                    }
                };
                Ok(Value::Int(result.unwrap()))
            }
        }
    }

    /// `sdiv` is undefined for these operand pairs. Every operand is a
    /// constant today, so they are rejected at compile time; runtime operands
    /// will need a check in the emitted code instead.
    fn check_division(
        &self,
        lhs: IntValue<'ctx>,
        rhs: IntValue<'ctx>,
        at: &'src str,
    ) -> Result<(), Diagnostic> {
        let message = match (
            lhs.get_sign_extended_constant(),
            rhs.get_sign_extended_constant(),
        ) {
            (_, Some(0)) => "division by zero",
            (Some(i64::MIN), Some(-1)) => "division result does not fit in 64 bits",
            _ => return Ok(()),
        };
        Err(Diagnostic::at(message.to_string(), self.source, at))
    }

    fn int_operand(
        &self,
        expr: &Expr<'src>,
        op: BinOp,
        at: &'src str,
    ) -> Result<IntValue<'ctx>, Diagnostic> {
        match self.value(expr)? {
            Value::Int(int) => Ok(int),
            Value::Str(_) => Err(Diagnostic::at(
                format!("cannot apply `{op}` to a string"),
                self.source,
                at,
            )),
        }
    }

    fn str_ptr(&self, konst: StrConst<'ctx>) -> inkwell::values::PointerValue<'ctx> {
        let zero = self.context.i32_type().const_int(0, false);
        unsafe {
            self.builder
                .build_gep(
                    konst.ty,
                    konst.global.as_pointer_value(),
                    &[zero, zero],
                    "strptr",
                )
                .unwrap()
        }
    }

    fn int_format(&mut self) -> StrConst<'ctx> {
        *self
            .int_format
            .get_or_insert_with(|| const_str(self.context, &self.module, "%lld\n"))
    }

    fn const_str(&self, text: &str) -> StrConst<'ctx> {
        const_str(self.context, &self.module, text)
    }
}

fn const_str<'ctx>(context: &'ctx Context, module: &Module<'ctx>, text: &str) -> StrConst<'ctx> {
    let mut bytes = text.as_bytes().to_vec();
    bytes.push(0);

    let ty = context.i8_type().array_type(bytes.len() as u32);
    let global = module.add_global(ty, None, ".str");
    global.set_linkage(Linkage::Private);
    global.set_unnamed_addr(true);
    global.set_constant(true);
    global.set_initializer(&context.const_string(&bytes, false));

    StrConst { ty, global }
}
