# Compilation model: answers to the unblocked open questions

Answers given by the author on 2026-08-27 to open questions 1-6 in `docs/design/0002-decision-compilation-model.md`.
Questions 7 and 8 are blocked on the memory-strategy and paradigm decisions, and were not asked.

This is a note, so it carries zero standing.
The answers are not part of the design record until they are written into the design doc itself.
This file exists so that a real wishlist entry and a stage 2 pass have something to draw on.

## The answers

1. **One execution semantics or two:** one.
   The compiled program has exactly one way to run: ahead-of-time to native code.
   No interpreted fallback mode for running a finished program.
1. **What engine runs comptime code:** a separate interpreter, embedded in the compiler, used only for comptime.
   This is Zig's actual approach despite its "one language" framing: an internal interpreter for comptime, the real backend for everything else.
1. **No-ceremony single-file run:** compile-then-exec is fine, matching `go run`.
   The user never sees a build step or artifact; paying AOT codegen latency under the hood is acceptable.
1. **Language-server incremental-query latency:** served by the compiler's own incrementality, not a separate tooling-only layer.
   Whatever makes the compiler fast incrementally is also what serves the language server.
1. **Is a runtime acceptable:** no runtime at all.
1. **Embeddability:** not a priority right now.
   Left genuinely open, not defaulted either way.

## What this rules out

Everything below is derivation, not authored answer, and it is the part most likely to be wrong.

- **Candidate C, bytecode plus a tiered JIT: dead.**
  Answer 5 removes it outright; a JIT is a runtime component by definition.

- **Candidate B, bytecode interpreter with no native codegen: dead.**
  Answer 1 says the compiled program has exactly one way to run, ahead-of-time to native code, which is the opposite of B's premise.

- **Candidate D, interpreter for dev/comptime plus AOT for release, as originally sketched: dead.**
  Answer 3 explicitly rejects a separate interpreted dev-mode as the answer to the no-ceremony scenario, and answer 1 rejects a user-facing interpreted execution mode outright.
  The comptime interpreter from answer 2 does not resurrect D: it is compiler-internal and never runs user program logic after compilation, so it is not a second *user-facing* execution semantics.

- **Candidate E, native codegen plus a separate rust-analyzer-style incremental layer: dead as sketched.**
  Answer 4 chooses the opposite: incrementality lives in the compiler's own pipeline, not in a second tool.
  This is a real architectural commitment, not a detail: the compiler's core has to be structured as an incremental/query-based system from the start (salsa-shaped or similar), because both normal builds and the language server draw on the same incrementality.

- **Candidate A, pure ahead-of-time to native code via LLVM: strongest survivor, effectively chosen,** with two amendments neither original sketch of A anticipated: a compiler-internal comptime interpreter (answer 2), and a mandatory incremental/query architecture for the compiler itself (answer 4).

## The shape this suggests

Ahead-of-time compilation to native code through the existing LLVM/inkwell path is the only way a compiled program ever runs.
`lang run` is `go run`-shaped sugar over that same pipeline, no separate fast path.
The compiler is architected as an incremental query system from the start, since that incrementality has to serve both ordinary builds and the language server.
Comptime code executes on a separate interpreter embedded in the compiler, never on the AOT backend, and never exposed as a user-facing execution mode.

That is closer to a real decision than a rough shape.
It is not yet a design, and the tensions below have to be resolved before it becomes one.

## New tensions the interview created

Each one is a candidate open question for the stage 2 pass.

1. **Answers 1 and 2 need "execution semantics" defined precisely before they read as consistent.**
   "One execution semantics" (answer 1) and "a separate interpreter for comptime" (answer 2) are only compatible under a specific reading: the rule in answer 1 is scoped to *how a compiled program runs*, not to every execution engine that exists inside the compiler.
   The design doc needs to state that scoping explicitly, or the two answers read as contradictory.

1. **Comptime/runtime behavioral identity is now a hard requirement, not a nice-to-have.**
   Zig's comptime design treats "comptime code behaves identically to runtime code" as load-bearing precisely because it runs on a different engine than the final binary.
   The same two-engine split exists here (answer 2), so the same burden applies: either a conformance suite that runs representative programs through both the comptime interpreter and the AOT backend and checks agreement, or a comptime language subset narrow enough that divergence is structurally impossible.
   Neither is decided.

1. **Answer 4 commits the compiler's core architecture before a compiler exists yet.**
   Building the compiler as an incremental query system from day one (to serve both builds and the language server) is a bigger up-front architectural bet than "compile whole module, ship it, add incrementality later."
   Whether the current inkwell-based compiler in the repo is already shaped for this, or would need restructuring, is unexamined.

1. **The comptime interpreter and answer 4's incrementality have to interact, and how is unexamined.**
   If comptime execution is expensive and the language server re-elaborates on every keystroke, comptime results need caching keyed on something, the same problem Julia's precompilation solves for its JIT.
   Whether that cache is part of the same incremental query system from answer 4, or a separate mechanism, is open.

1. **Embeddability (answer 6) was left open, not resolved.**
   Nothing here rules an embeddable interpreter in or out; it just was not allowed to drive this decision.
   If it becomes a priority later, it re-opens candidate B or D in a narrow form: an embeddable interpreter that runs *alongside* the AOT compiler for host-embedding use cases specifically, not as the language's primary execution model.

## Answers to the tensions

Answers given by the author on 2026-08-27 to the tensions raised above.

1. **Semantics scope:** implementation detail only.
   The spec describes one semantics for the language; the comptime interpreter's existence is unconstrained by the spec as long as observable behavior matches.
1. **Agreement mechanism:** a restricted comptime subset.
   Comptime code is a deliberately narrower language than the full language, no unchecked FFI, no undefined-behavior sources, so that divergence between the interpreter and the AOT backend is structurally impossible rather than merely tested for.
1. **Comptime caching:** yes, comptime results live in the same incremental query graph as everything else.
   Editing untouched code reuses cached comptime values through the same invalidation machinery as the rest of the compiler.

Tension 3 (whether the existing compiler already fits an incremental architecture) turned out not to need an answer.
`src/main.rs` is a roughly 90-line stub that hardcodes parsing a single `print "..."` line and emits LLVM IR directly: no AST, no passes, no incrementality either way.
There is nothing built yet to reconcile the incremental-query commitment against; it is a greenfield architectural choice, not a retrofit.

### What the tension answers narrow further

Answer 2 (restricted comptime subset) is now in mild tension with answer 1 (spec treats the interpreter as an implementation detail): if comptime is a restricted subset, that restriction is user-visible (some full-language constructs are rejected inside a comptime block), so the spec cannot stay silent about the interpreter's existence entirely.
What it can stay silent about is *how* comptime is executed; it cannot stay silent about *what is allowed* inside comptime.
The spec needs a defined comptime subset (a real stage 2 artifact for the comptime feature itself), even though the execution engine behind it is unspecified.

## The runtime-definition clarification does not revive candidate C

`docs/notes/concurrency-interview-answers.md` (2026-08-27) clarified that "no runtime at all," as answered across this record and the memory-strategy record, meant no separately-installed binary, not zero inline machinery, a runtime component is fine if it ships as part of the single compiled binary.
That clarification was flagged there as a live loose end for this record, since answer 5's rejection of candidate C, a bytecode VM plus a tiered JIT, cited "a runtime component by definition" as the reason.

Checked against this record's own answers, the loosened definition does not change candidate C's fate.
Answer 1 is the stronger and independent objection: "the compiled program has exactly one way to run: ahead-of-time to native code. No interpreted fallback mode for running a finished program."
A tiered JIT does not run a program ahead-of-time to native code from the start, it starts by interpreting and lazily promotes hot paths, which contradicts answer 1 regardless of whether the JIT engine itself would be statically linked into the binary with nothing separate to install.
The same reasoning kills candidates B and D again too: both were already rejected primarily on answer 1's "one execution semantics" grounds, not on the runtime objection, so neither is revived either.

This also clarifies why the compiler-internal comptime interpreter (answer 2) was never actually in tension with "no runtime at all" in the first place, under either definition: it runs inside the *compiler*, at build time, and never ships as part of the *compiled program*, so it was never the kind of runtime component answer 5 or answer 1 were ruling out.

Candidate C stays dead.
The loose end is closed, not by the runtime definition, but by a second, independent, and unaffected objection this record already had on record.

## Not asked

Questions 7 and 8 remain blocked: the memory-strategy record (`0001`) and the paradigm record.
Nothing above should be read as deciding either, though answer 5 (no runtime at all) directly matches memory record answer 3 (also "no runtime at all"), which substantially resolves the mutual-blocking tension both records flagged, in the same direction, without either formally depending on the other.
