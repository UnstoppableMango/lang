# Feature wishlist

The stage 0 idea pool for the language, per the [feature design workflow](workflow.md).
Anything here is an idea worth remembering, nothing more.
Entries carry no commitment, and contradictory entries may coexist.

Entry format:

- One sentence per feature.
- Optionally a minimal example, inline.
- Optionally prior art: other languages, books, videos, papers.

Every foundational decision (paradigm, compilation model, memory strategy) is currently open, so entries must not be read as presuming any of them.

## Design philosophy

- Few ways to do one thing, favoring convention over configurable primitives.

  - Prior art: Go's design stance.

- Hard to do the wrong thing by accident; escape hatches are marked visibly (e.g. `unsafe`), never the quiet default.

  - Prior art: Rust's `unsafe` blocks.

## Syntax and surface

- Pattern matching with compiler-checked exhaustiveness, nested destructuring, guards, and or-patterns.

  - Prior art: OCaml, Rust, Erlang.

  ```
  match x { Some(n) => n, None => 0 }
  ```

- Whitespace significance chosen deliberately, not inherited by habit.

  - Prior art: Python, Haskell's layout rule, Go's semicolon insertion.

- Expression-oriented control flow: a function body's trailing expression is its return value.

  - Prior art: OCaml, Rust's tail expressions.

  ```
  fn add(a, b) = a + b
  ```

- Structural interface satisfaction, checked statically.

  - Prior art: Go interfaces, TypeScript.

## Types

- Type inference scoped locally within a function body, with explicit signatures at function and module boundaries.

  - Prior art: Rust, Swift, Kotlin.

- Nullability absent from the language; optionality expressed as a real type.

  - Prior art: Rust, Kotlin, Tony Hoare's "billion dollar mistake" talk.

  ```
  Option<T> or T? as a real type
  ```

- A distinct path/file type, source-relative and resolved against the file that referenced it.

  - Prior art: Nix path literals, Zig `@embedFile`, Go `//go:embed`.

- Purity (what a function may do, e.g. perform I/O) tracked in its type.

  - Prior art: Haskell's `IO` type, Koka/Unison effect systems.

## Semantics and evaluation

- Expressions over statements wherever possible, so most constructs yield a value.

  - Prior art: Rust, ML family.

  ```
  let x = if c { a } else { b }
  ```

- A small, clearly defined set of implicit conversions, possibly empty.

  - Prior art: Go (none), C (many, regretted).

- First-class short-circuiting syntax for failure- and absence-shaped types (`Option`, `Result`).

  - Prior art: Rust's `?`, Swift's optional chaining.

## Memory and resources

- Deterministic resource cleanup tied to scope.
  - Prior art: RAII in C++/Rust, Python `with`, C# `using`.
- The default memory regime is ownership and borrowing, checked statically with no runtime cost.
  - Prior art: Rust ownership and borrowing.
- A second, explicitly marked regime exists for values that are genuinely shared or graph-shaped, reclaimed by compiler-inserted reference counting rather than a tracing collector.
  - Prior art: Perceus-style compiler-inserted reference counting.
- Which regime a value belongs to is a property of its type's definition, not something restated in every function signature that passes the value along.
- A type opts into deterministic cleanup of a non-memory resource by declaring a release action, run wherever the language reclaims a value of that type.
  - Prior art: C++/Rust RAII (`Drop`).
- A value crossing the C FFI boundary pins to a stable address or transfers ownership across it.
- Values are movable by default, and only a value pinned at a boundary that requires one is guaranteed a stable address.
  - Prior art: Rust's move semantics and `Pin`.
- The reference-counted regime provides a weak-reference tool so a cycle can be broken manually, rather than being collected automatically.
  - Prior art: Swift's `weak`/`unowned`.

## Errors and failure

- Recoverable errors as values, distinct from unrecoverable panics.

  - Prior art: Rust, Go, Erlang's "let it crash".

  ```
  Result<T, E> plus panic
  ```

- Error messages designed as a first-class feature with their own quality bar.

  - Prior art: Elm, Rust diagnostics.

## Concurrency

- A concurrency story where data races are hard or impossible to express.
  - Prior art: Rust `Send`/`Sync`, Erlang isolated processes, Pony.
- Structured concurrency: a child task cannot outlive its parent scope.
  - Prior art: Trio, Kotlin coroutines, "Notes on structured concurrency" (Nathaniel J. Smith).
- Concurrency without function coloring: any function can block, no async-marked type infects its callers.
  - Prior art: Go goroutines, Erlang processes, Java virtual threads (Loom).
- Tasks are stackful and scheduled M:N onto a pool of OS threads, not compiled to stackless state machines.
  - Prior art: Go's scheduler.
- Moving or sharing a value across a task boundary is checked at compile time by the same ownership-and-borrowing regime the memory strategy already uses, not a separate capability system.
  - Prior art: Rust's `Send`/`Sync`.
- `scope` and `spawn` are language keywords, not a library convention added after the fact.
  - Prior art: Java's structured concurrency JEPs, Trio.
- Both lock-shaped primitives and channels exist as synchronization primitives.
  - Prior art: Rust's `Mutex<T>` alongside channels.
- The language treats very-high-fan-out I/O, tens or hundreds of thousands of concurrent in-flight tasks, as a target workload.
  - Prior art: Go, Erlang.

## Modules and packaging

- A module system where dependencies are explicit and cycles are rejected.
  - Prior art: Go packages, ML modules.
- Reproducible builds as a design constraint from day one.
  - Prior art: Nix, Go modules, cargo lockfiles.
- A directory is a module; no header line to keep in sync with the file's location.
  - Prior art: Go packages.
- Dependencies declare the capabilities they require (network, filesystem, etc.), checked by the compiler.
  - Prior art: Deno permissions, capability-secure languages (E, Monte).
- Build tooling reports the cost of each dependency (build time, size, transitive count) on every build.
- Taking a dependency forks and trims its source to what's reachable from your own call sites.

## Interop and embedding

- A C-compatible foreign function interface through an explicit boundary.
  - Prior art: Rust's `extern "C"`, Zig's C interop.
- Embeddability of the language or its runtime into a host program.
  - Prior art: Lua, Wasm components.

## Metaprogramming

- Compile-time code execution as the metaprogramming model.

  - Prior art: Zig comptime, Rust const eval, Lisp macros as the ancestor.

  ```
  const fn / comptime
  ```

- Derivable boilerplate (equality, formatting, serialization) generated from type definitions.

  - Prior art: Rust `#[derive]`, Haskell `deriving`.

- The compiler's AST/IR as a public, versioned library, imported directly by the formatter, linter, and language server.

  - Prior art: Go's `go/ast`, Rust's `syn`/`proc-macro2`.

## Tooling and developer experience

- One canonical formatter with no configuration, shipped with the language.
  - Prior art: gofmt, zig fmt.
- A language server planned alongside the compiler, exposing semantic queries (type-of, go-to-definition) as a library API.
  - Prior art: rust-analyzer, the LSP paper trail.
- Fast, incremental compilation: editing one function doesn't force re-checking a whole file or module.
  - Prior art: Go's compilation speed goals, rust-analyzer's query-based (salsa) architecture.
- A comprehensive toolchain under one command: build, format, test, vet, docs, dependency management.
  - Prior art: Go's `go` command.
- Build tooling built on existing content-addressed infrastructure (Nix-shaped).
- A canonical, single textual representation for any syntax tree, so reformatting alone never produces a diff.
- Stable per-definition identity that survives reformatting, enabling a semantic blame/diff tool.

## Foundational decisions

- A written specification is the authority; the reference compiler is one conformant implementation.
  - Prior art: C, Scheme.
- A no-ceremony way to run a single file, with no separate build step and no second scripting dialect.
  - Prior art: `go run`, `cargo script`.
- The compiler compiles ahead-of-time to native code via LLVM, with no runtime.
  - Prior art: Go, Rust.
- Comptime code runs on a compiler-internal interpreter, restricted to a subset of the language, operating over the same IR the AOT backend consumes.
  - Prior art: Zig `comptime` (restricted-subset compile-time execution).
- The compiler is architected as an incremental query system from the start, serving both ordinary builds and the language server, including caching comptime results in the same query graph.
  - Prior art: rust-analyzer's salsa-based incremental architecture, Julia's JIT latency and precompilation as a cautionary case for what an unmanaged comptime cache costs.
- Data with a closed, recursive shape, such as the compiler's own AST, is modeled as an algebraic data type and inspected by pattern matching.
  - Prior art: OCaml, Rust, Haskell.
- Values are immutable by default, and mutating a binding requires an explicit marker.
  - Prior art: Rust's `let mut`.
- A method call is sugar over a free function that takes the value as its first argument, rather than a second, independent way to attach behavior to data.
  - Prior art: Rust's method resolution, D's uniform function call syntax.
- A function counts as one of a type's methods by taking that type as its first parameter, and a value satisfies an interface purely by that method set matching, with no declared relationship anywhere between the type and the interface.
  - Prior art: Go interfaces.
- A function only counts toward a type's method set if it is defined in the same module as the type.
  - Prior art: Go's package-scoped method sets.
- The reference compiler is eventually self-hosted, rewritten in its own target language.
  - Prior art: Go (self-hosted since 1.5), Rust (self-hosted).

## Rejected

Nothing yet.
Rejected entries move here with a one-sentence reason, per the workflow.
