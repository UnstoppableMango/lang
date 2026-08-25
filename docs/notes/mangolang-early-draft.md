# mangolang: an early draft, not a decision

This note preserves an early brainstorm from a prior incarnation of this project, back when it went by the name mangolang on a branch called `fresh`, with a GOALS.md and SPEC.md instead of the current docs/wishlist.md and docs/design/ workflow.
None of it was ever run through the stage 0-4 process this repo now uses.
It predates that process.
Nothing here has any standing, per \[[docs/workflow.md]\]: no foundational decision (paradigm, compilation model, memory strategy, implementation language) may be presumed until a decision record reaches stage 3, and none of these ever did.
Recorded here as raw material for future notes sessions to riff on, argue with, or ignore, not as a starting point that carries any weight over any other idea.

## What the draft said

A pet language project for personal learning and growth, explicitly not optimizing for adoption.

- Name: mangolang, file extension `.um`
- Memory: region/arena allocation, no GC, no Rust-style ownership
- Type system: Hindley-Milner inference, parametric polymorphism, ADTs, structural interfaces
- Mutability: immutable by default, `mut` explicit
- Concurrency: goroutine-like green threads plus channels (CSP), leaning on immutable-by-default for safe sharing
- Error handling: `Result<T,E>` plus pattern matching, no exceptions, no `?` operator
- Syntax: novel, no semicolons, tabs for indentation, type-before-name (`Int x`), `|>` pipeline operator, `camelCase` for values, `PascalCase` for types
- IR targets: Go source first (bootstrap), then C and LLVM IR, multiple targets in parallel
- DX priorities: fast compile times, excellent error messages, LSP-first, an integrated formatter

Open questions the draft explicitly deferred and never answered: typeclass syntax, how arena annotations would look in source, the concrete concurrency primitives, string interpolation, FFI, operator overloading.

## Why it's worth keeping around

Some of this rhymes with notes already written under the current process.
No-exceptions-explicit-errors (\[[no-exceptions-explicit-errors]\]) matches the draft's `Result<T,E>` instinct.
Green threads (\[[green-threads-threading-model]\]) matches the draft's goroutine-like concurrency.
Functional implicit return (\[[functional-implicit-return]\]) and no-null (\[[no-null-type-or-representation]\]) are both plausibly downstream of the same taste that produced this draft, immutable-by-default, ADT-and-pattern-match-shaped, allergic to null and exceptions both.

Other parts don't obviously rhyme with anything written since, and are worth treating as genuinely open rather than quietly re-adopted: the region/arena memory strategy specifically (as opposed to GC or ownership, a foundational choice with its own decision record still to be drafted), the Go/C/LLVM multi-target IR plan, and the `Int x` type-before-name syntax choice.
Whether that's because they were explicitly reconsidered, or just because nobody has gotten to them yet under the new process, is not something this note can answer.
Worth treating each one as a fresh question, not a foregone conclusion, precisely because the draft's own claim to have "decided" them never went through any gate.
