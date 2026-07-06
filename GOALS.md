# mangolang Goals

## Vision

mangolang is a general-purpose language built for personal exploration of language design and compiler implementation. Primary targets: tooling/CLI, backend services, and systems programming.

The guiding thesis: take Go's simplicity and operational directness, combine it with the expressiveness and safety guarantees of ML-family languages (OCaml, Haskell, F#), and invest heavily in developer experience.

## Core Philosophy

- **Small, orthogonal feature set.** Every feature must justify its complexity cost.
- **Explicit over implicit, except types.** Types are inferred; effects are visible; errors are values.
- **DX is a first-class feature.** Compiler errors, LSP, formatter, and compile speed are not afterthoughts.
- **No null. No exceptions. No magic.**

## Language Properties

### Type System
- Hindley-Milner type inference; minimal annotations required
- Parametric polymorphism (generics) with typeclass/trait-style constraints
- Algebraic data types: sum types (enums with payloads) + product types (records/structs)
- Exhaustive pattern matching; compiler enforces completeness
- Structural interfaces; implicit satisfaction, no `implements` keyword
- Immutability by default; `mut` annotation for mutable bindings

### Memory Model
- Region/arena-based allocation; no garbage collector, no Rust-style ownership
- Allocation scoped to regions; deterministic cleanup at region boundary
- Goal: predictable performance without borrow checker complexity

### Concurrency
- Goroutine-like lightweight cooperative threads
- Channels as primary communication (CSP model)
- Immutability by default makes state sharing safe by construction
- No shared mutable state without explicit opt-in

### Error Handling
- No exceptions
- `Result<T, E>` as the standard error type
- Errors handled via pattern match; explicit, exhaustive, visible
- Error values compose naturally with ADTs

### Syntax Principles
- No semicolons; newline-aware grammar
- Indentation with tabs (enforced by formatter)
- Type annotations before name: `Int x`, `fn foo(Int x) -> String`
- Pipeline operator `|>` for readable data transformation chains
- Novel aesthetic; not an imitation of any existing language

### Intermediate Compilation Targets
- **Go source**: fast iteration, leverage Go toolchain, easy to bootstrap
- **C**: maximum portability, simple codegen, widely supported
- **LLVM IR**: production path, full optimization pipeline
- Multiple targets supported in parallel; Go path ships first

### Developer Experience (non-negotiable)
- Sub-second incremental compilation
- Elm/Rust-quality error messages: point to cause, suggest fix, never cryptic
- LSP server as a first-class deliverable alongside the compiler
- Integrated formatter; one canonical style, no configuration debates

## Non-Goals

- Runtime reflection
- Implicit type coercions or numeric promotion
- Null or nil (use `Option<T>`)
- Exceptions
- Global mutable state
- Macros (in v1)
- Higher-kinded types (deferred; evaluate after core type system matures)
