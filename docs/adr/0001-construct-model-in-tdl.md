---
status: proposed
---

# Describe the language's constructs in TDL before writing them

The shapes of the language's constructs are described in TDL, as the construct model, before the compiler is written against them, so a construct is designed as data apart from the Rust that builds it.
The model starts as today's syntax tree and grows into a resolved representation, one milestone at a time; each milestone's spec lists its changes to the model.
It holds only what is built or accepted under the presume rule, and says nothing about what a construct means; a decision issue may carry its candidates as fenced `tdl` sketches instead.
Rust details such as lifetimes, `Box`, and source spans stay out of the model; spans come back through a Rust directive.

The model becomes normative through a Rust backend, `tdl-gen-rust`, built in tdl on a Rust plugin SDK beside its .NET one, because a backend is written in its own language where the tooling permits.
Once it exists, lang commits the generated Rust and checks it with `tdl gen --verify`; until then the model is descriptive, CI runs only `tdl check` over it, and drift from `src/ast.rs` is caught in review.
Where TDL cannot state a shape the language needs, that is a `needs-triage` issue in tdl, not a reason to bend the model.

## Considered Options

- Hand-written Rust as the only description: nothing stands between a decision and its implementation, which is the gap this fills.
- A Rust backend written in Go inside tdl: rejected, since tdl's shared emit helpers are moving into the host and a backend is meant to be written in its target language.
- Generating at build time from `build.rs`: rejected, so the generated diff shows in review and a cargo build needs no tdl on `PATH`.
