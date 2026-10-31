# Influences

Traits of other languages the author likes or dislikes, kept as a soft guide for design.
The README's principles bind every design choice; an influence only informs one, and a principle wins any conflict.
Nothing here is a commitment to build anything, and nothing here is a decision.

A decision may cite an influence, and its candidates may say which influences they draw on, but no candidate has to honour one.
Likes may pull in different directions: the author can draw on one in one place and another elsewhere, and two such Likes are not a tension.
A trait missing from this list is not a dislike.

Each entry is one line: Like or Dislike, the language, the thing, and why.
A topic appears once it has an entry.

## Concurrency

- Like, Go: goroutines and channels, where any function can block without being marked, so no async type spreads to its callers.

## Types

- Like, Rust and Kotlin: optionality as a real type (`Option<T>`, `T?`), so "nothing" is a value the type admits rather than a bit pattern every type carries.
- Like, Nix: path literals as their own type, resolved against the file that wrote them, so a path means the same thing wherever the program runs.
- Like, Haskell, Koka, and Unison: purity tracked in the type (`IO`, effect systems), so what a function may do is part of what it says.
- Like, Go and TypeScript: structural interfaces satisfied with no declaration and checked statically, so a type conforms by its shape, with no association to keep in sync.
- Like, Rust, Haskell, and Swift: everyday monadic types (`Option`, `Result`, comprehensions) as friendly tools rather than abstractions to fear.

## Paradigm and syntax

- Like, Rust: immutable by default, with `mut` required to mutate a binding, so the common case of mutation is visible where the binding is made.
- Like, Rust and D: a method call is sugar for a function call taking the value first (Rust's `x.foo()` is `Type::foo(x)`, D's uniform function call syntax), so there is one mechanism with two spellings.
- Like, OCaml and Rust: a function body's trailing expression is its return value, with no `return` needed for the common case.
- Like, F#, Elixir, and OCaml: the `|>` pipeline operator, which reads a chain of calls in the order it runs.
- Dislike: semicolons, a terminator the parser could do without.

## Pattern matching

- Like, OCaml and Rust: exhaustiveness checking, nested destructuring, guards, and or-patterns, so a match takes a value apart in one place and the compiler says which case is missing.

## Errors

- Like, Rust and Go: recoverable errors as values, kept apart from panics, so a function that can fail says so in its signature.
- Dislike, Java, C#, and Python: exceptions, control flow that jumps across stack frames unseen.

## Modules and packaging

- Like, Go: a directory is a package, with no header line to keep in sync with the file's location.
- Like, Go and ML: explicit dependencies, with import cycles rejected, so the dependency graph stays a graph a reader can follow.
- Like, OCaml: applicative functors and the ML module system, worth studying for how modules can be parameterised.
- Like, Deno: a program gets only the network, filesystem, and other access it is granted, so what code can do is visible before it runs.

## Tooling

- Like, Go and Zig: one canonical formatter shipped with the language, with no configuration, so there is nothing to argue about.
- Like, Go: the whole toolchain under one command (build, format, test, vet, docs, dependencies).
- Like, Go: its simplicity philosophy, where a feature has to earn its place.
- Like, Go: fast compilation, kept as a goal rather than left to chance.
- Like, Go and Rust: running one file with no visible build step (`go run`, `cargo script`), so the user never sees a build or its output.
- Like, Go and Rust: the language's syntax tree as a public, importable library (`go/ast`, `syn`), so a tool needs no parser of its own.
- Like, C and Scheme: a written specification is the authority, and the compiler is one implementation of it.

## Compilation and metaprogramming

- Like, Zig: `comptime`, ordinary code the compiler runs while building, so build-time logic needs no second language.
- Like, Rust and Zig: C interop through an explicit boundary, so crossing into C is visible where it happens.
