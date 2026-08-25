# Feature wishlist

The stage 0 idea pool for the language, per the [feature design workflow](workflow.md).
Anything here is an idea worth remembering, nothing more.
Entries carry no commitment, and contradictory entries may coexist.

Entry format:

- One sentence per feature.
- Optionally a minimal example, inline.
- Optionally prior art: other languages, books, videos, papers.

Every foundational decision (paradigm, compilation model, memory strategy, implementation language) is currently open, so entries must not be read as presuming any of them.

> The entries below are seeds that demonstrate the format.
> Replace or remove them freely as real ideas accumulate.

## Syntax and surface

- Pattern matching with compiler-checked exhaustiveness.
  - Prior art: OCaml, Rust, Erlang.

  ```
  match x { Some(n) => n, None => 0 }
  ```
- Significant or insignificant whitespace chosen deliberately rather than inherited by habit.
  - Prior art: Python, Haskell layout rule, Go's semicolon insertion.

## Types

- Type inference strong enough that annotations are for communication, not for the compiler.
  - Prior art: Hindley-Milner, F#, Elm.
- Nullability absent from the language, with optionality expressed in the type system.
  - Prior art: Rust, Kotlin, Tony Hoare's "billion dollar mistake" talk.

  ```
  Option<T> or T? as a real type
  ```

## Semantics and evaluation

- Expressions over statements wherever possible, so most constructs yield a value.
  - Prior art: Rust, ML family.

  ```
  let x = if c { a } else { b }
  ```
- A clearly defined and small set of implicit conversions, possibly empty.
  - Prior art: Go (none), C (many, regretted).

## Memory and resources

- Deterministic resource cleanup tied to scope rather than to a finalizer that may never run.
  - Prior art: RAII in C++/Rust, Python `with`, C# `using`.
- A memory strategy chosen as an explicit decision record, not an accident of the first implementation.
  - Prior art: tracing GC (Go, Java), ownership (Rust), reference counting (Swift).

## Errors and failure

- Recoverable errors as values distinct from unrecoverable panics.
  - Prior art: Rust, Go, Erlang's "let it crash".

  ```
  Result<T, E> plus panic
  ```
- Error messages designed as a first-class feature with their own quality bar.
  - Prior art: Elm, Rust diagnostics.

## Concurrency

- A concurrency story where data races are hard or impossible to express.
  - Prior art: Rust `Send`/`Sync`, Erlang isolated processes, Pony.
- Structured concurrency, where a child task cannot outlive its parent scope.
  - Prior art: Trio, Kotlin coroutines, "Notes on structured concurrency" (Nathaniel J. Smith).

## Modules and packaging

- A module system where dependencies are explicit and cycles are rejected.
  - Prior art: Go packages, ML modules.
- Reproducible builds as a design constraint from day one.
  - Prior art: Nix, Go modules, cargo lockfiles.

## Interop and embedding

- A C-compatible foreign function interface as the lowest common denominator.
  - Prior art: every serious language eventually.
- Embeddability of the language or its runtime into a host program.
  - Prior art: Lua, Wasm components.

## Metaprogramming

- Compile-time code execution rather than a separate textual macro language.
  - Prior art: Zig comptime, Rust const eval, Lisp macros as the ancestor.

  ```
  const fn / comptime
  ```
- Derivable boilerplate (equality, formatting, serialization) generated from type definitions.
  - Prior art: Rust `#[derive]`, Haskell `deriving`.

## Tooling and developer experience

- One canonical formatter with no configuration, shipped with the language.
  - Prior art: gofmt, zig fmt.
- A language server planned alongside the compiler, not bolted on afterward.
  - Prior art: rust-analyzer lessons, the LSP paper trail from Microsoft.
- Fast compile times treated as a feature with a budget, measured continuously.
  - Prior art: Go's compilation speed goals, Jonathan Blow's Jai demos.

## Rejected

Nothing yet.
Rejected entries move here with a one-sentence reason, per the workflow.
