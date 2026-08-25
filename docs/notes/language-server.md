# Language server developed alongside the compiler

Expanding the wishlist entry ("a language server planned alongside the compiler, not bolted on afterward") past the one sentence.
\[[dev-tooling-philosophy]\] already named the mechanism this needs (shared AST/IR as a library) in passing.
This note tries to get concrete about what "alongside" means operationally, not just architecturally, and where it gets hard.

## What "bolted on afterward" actually looks like, so the alternative is legible

The failure mode being avoided has a specific shape, worth naming instead of gesturing at:

1. Compiler exists as a CLI binary that takes source, produces output, exits.
1. Years later, someone wants IDE support.
1. They cannot import the compiler as a library (it was never factored that way), so they write a second parser.
1. The second parser drifts from the first: new syntax lands in the compiler, the LSP chokes on it for months until someone ports the change.
1. The LSP settles for degraded, best-effort analysis (syntax highlighting, maybe scoping) because full type information requires reimplementing the type checker too, which nobody signs up for twice.

TypeScript, Go (pre-`go/packages`), and C++ (clangd took over a decade after clang existed) all lived through some version of this.
rust-analyzer is the counter-case people point to, but it's worth being honest that rust-analyzer is *also* a second implementation, a from-scratch analyzer that deliberately does not reuse rustc's internals, because rustc's internals were not shaped for incremental, query-based, error-tolerant use.
So "alongside the compiler" has at least two very different readings and rust-analyzer is evidence for the harder one, not the easy one.

## Two readings of "alongside," and they are not the same project

**Reading A: same library, two entry points.**
One compiler, factored as a library from day one.
`lang build` and `lang-server` both link against it.
Cheapest to build first, but only works if the compiler's internal architecture happens to fit what a language server needs (see next section) — and compilers are not usually written with that in mind, they're written to go from source to output as directly as possible.

**Reading B: same team, two implementations, kept in sync on purpose.**
This is what rust-analyzer actually is.
More expensive up front (two things to build), but sidesteps forcing the batch compiler into an architecture it doesn't want.
"Alongside" here means organizational, not literal code-sharing: the LSP team sits next to the compiler team, gets a seat at the syntax-design table, and the AST/IR *shapes* are shared (or at least kept isomorphic) even if the *code* is not.

The wishlist entry doesn't disambiguate these, and it matters a lot which one is meant, since they lead to opposite early decisions: Reading A says "design the compiler's internals for incrementality from the start." Reading B says "don't force that on the compiler, build a second analyzer, but make sure new syntax can't land without both being updated in the same PR."

## What a language server actually needs that a batch compiler doesn't

This is the crux, and it's not a small list:

- **Error tolerance.** A batch compiler can bail on the first syntax error. An LSP has to parse `let x = foo(` mid-keystroke and still produce a tree good enough to offer completions after the open paren. Error-tolerant parsing is a different parser design (see: Roslyn's red-green trees, rust-analyzer's rowan), not a mode flag on a normal recursive-descent parser.
- **Incrementality at the sub-file, sub-function grain.** Batch compilation is fine recompiling a whole file, or even a whole package, per build. An LSP needs "user typed one character, what's the minimal recomputation," which is a fundamentally different execution model (query-based / salsa-style incremental computation, not "run the pipeline top to bottom").
- **Concurrent, cancellable analysis.** The user keeps typing while the last request is still being answered. Batch compilers don't need to think about cancellation at all.
- **Position-addressable everything.** Every AST node needs a source span cheap enough to hold in memory for an entire open file, and every type/binding needs to be queryable *by position*, not just producible as a side effect of code generation. A compiler optimizing purely for throughput will happily discard this information as soon as it's not needed for codegen.
- **Multiple "versions" of the world at once.** Go-to-definition across an edit that hasn't been saved, of a file that's still open with unsaved changes, while the file on disk is what the batch compiler would see. The LSP's model of "what does the source look like" is a layered overlay, not a filesystem read.

None of these are exotic, individually, but none of them are things a batch compiler needs for its own job.
Building "alongside" cheaply requires either designing the compiler around these from day one (expensive, and a lot of that expense is paid before there's any user-visible payoff), or accepting Reading B and paying for two things.

## Doodle: what the shared boundary could look like, if Reading A is chosen

```
lang/compiler/          -- the library, no CLI awareness
  parse/                -- error-tolerant by construction, not error-tolerant as an afterthought
  ast/                  -- immutable, position-addressable nodes
  typecheck/             -- exposed as a query: type_of(node) -> Type, not just "does this program typecheck"
  ir/

lang/cmd/lang/          -- thin CLI: parse args, call compiler/, print or exit
lang/cmd/lang-server/   -- thin LSP: parse args, call compiler/, translate queries to LSP protocol messages
```

The load-bearing claim is `typecheck/` exposed as a query rather than a pass.
A batch compiler's type checker usually wants to say "typecheck this whole program, tell me pass/fail plus diagnostics."
An LSP wants to ask "what is the type of the expression under the cursor," possibly before the rest of the file even typechecks.
Those are different APIs over the same underlying algorithm, and retrofitting the query shape onto a pass-shaped checker later is close to a rewrite, not a refactor, which is the real cost of "bolted on afterward."

## Doodle: hover / go-to-def as the litmus test

A cheap way to tell, at any point in the compiler's development, whether "alongside" is actually happening or just being intended:

Can `lang-server` answer "what is the type of this expression" using *only* library calls into the compiler, with no source-text re-parsing of its own and no duplicated symbol table?

If yes, at whatever stage the compiler currently is (even pre-1.0, even with half the type system unbuilt), the architecture is on track.
If the honest answer requires the LSP to re-derive scoping or types itself because the compiler doesn't expose that as queryable state, that's the "second implementation" failure mode arriving early, quietly, one convenience shortcut at a time.
This could be a running check, not just a one-time design review, since it's the kind of property that erodes gradually rather than breaking all at once.

## Where this fights other wishlist/notes items

- \[[fast-incremental-compilation]\]: if incremental compilation is already a batch-build goal (not just an LSP goal), the query-based architecture reading A needs is *shared infrastructure*, not extra cost paid only for IDE support. This is the strongest argument for Reading A: if incrementality is wanted for `lang build` anyway (faster CI, faster local rebuilds), the LSP gets it nearly free. Worth checking whether that note and this one are actually asking for the same underlying architecture from two different angles.
- \[[arena-memory-model]\]: an arena-per-compilation memory strategy is great for a batch compiler (allocate, compile, throw the whole arena away) and actively hostile to an LSP (files stay open and get incrementally re-typechecked for the lifetime of an editor session, so nothing gets to be "done" and freed the way a batch run is done and exits). If arenas are chosen for the compiler's core memory strategy, the language server likely needs a different allocation discipline internally, which is friction against "same library, two entry points."
- \[[codegen-friendly-metaprogramming]\]: compile-time execution (comptime-style) means the type checker can run arbitrary user code during typechecking. An LSP needs to typecheck *incomplete, currently-being-edited* code constantly. Running arbitrary compile-time code against a half-written, syntactically-broken buffer on every keystroke is a real problem (rust-analyzer's answer for proc macros is essentially "sandbox it, cache aggressively, degrade gracefully when it can't run"). Worth flagging as a specific hard case rather than assuming metaprogramming and "LSP alongside the compiler" compose for free.

## Dead ends / cautions, recorded so they don't get rediscovered the hard way

- **"Just shell out to the compiler binary and parse its diagnostic output for the LSP."** This is bolted-on-afterward wearing a costume. It gets syntax errors and maybe type errors as text, but not hover types, not go-to-def, not completions, because none of those are things a batch compiler's stdout was ever designed to carry. Doesn't satisfy the wishlist entry even though it superficially "uses the compiler."
- **"Build the LSP once the language is stable, since building it early means rewriting it every time syntax changes."** This is the argument for deferring, and it's exactly backwards per the wishlist framing: rewriting the LSP every time syntax changes is only painful if the LSP is a second implementation (Reading B) or built against an unstable ad-hoc API. If it's a thin layer over a stable library boundary (Reading A), syntax changes are absorbed by the one shared AST, same as the batch compiler absorbs them.
- **Assuming Reading A and Reading B are the same commitment.** They aren't, and picking neither by default (building a batch compiler with no query API, then reaching for "we'll just add LSP support" without deciding which reading that means) is how "bolted on" happens even when everyone agreed in principle to avoid it.

## Threads worth pulling later

- Whether "compiler as a library with a query API" should itself be a decision record, since it's foundational enough to gate a lot of downstream tooling work (formatter, linter, and LSP in \[[dev-tooling-philosophy]\] all want the same boundary).
- What incremental typechecking actually costs to build for *this* language's type system, once that system is chosen; Hindley-Milner-style inference is notoriously non-local (one edit can change inferred types far away), which is a harder incremental story than a language with mostly-local type annotations.
- Whether the language server ships as part of the same toolchain binary (`lang server` as a subcommand, per the \[[dev-tooling-philosophy]\] "one verb space" doodle) or as a separate artifact with its own release cadence, since an LSP arguably needs to ship *faster* than the compiler (editor users want fixes now) which is in tension with one shared version number.
- Whether error-tolerant parsing (needed for the LSP) should just be *the* parser, used by the batch compiler too, rather than maintaining a tolerant-parser/strict-parser split. Rust and C# both eventually concluded one tolerant parser used everywhere beats two parsers.
