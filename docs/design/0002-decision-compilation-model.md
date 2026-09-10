---
title: 'Decision: Compilation model'
stage: 1
status: active
created: 2026-08-27
updated: 2026-08-27
---

# Decision: Compilation model

## Motivation

A language's compilation model, whether source becomes native machine code ahead of time, is interpreted, or takes some mixed path, shapes almost every other tool built around it: how fast edit-compile-run feels, what a language server can answer and how quickly, how metaprogramming that runs at compile time is itself executed, and what a deployed program depends on to start.
The wishlist already commits to outcomes that presume an answer here without naming one: a no-ceremony way to run a single file with no separate build step, fast incremental compilation where editing one function doesn't force rechecking a whole module, and compile-time code execution as the metaprogramming model.
None of those three can be designed honestly until this record picks a direction, because each reads differently depending on whether "compile" means "emit a binary" or "start interpreting."

The current compiler already has a starting point: it targets LLVM IR through `inkwell`, and `hack/hello.lang` compiles and runs end to end through that path.
That is existing groundwork, not a foreclosed decision; `AGENTS.md` states the Rust/inkwell/nom choice is a starting point that may change if a switch proves worthwhile, and this record should not presume the native-codegen answer just because it is what exists today.

### Scenario: running one file with no ceremony

A user writes a ten-line file and wants to run it the way they would run a Python script: one command, no project scaffolding, no visible build artifact to clean up afterward.
`go run` satisfies this by compiling to a temp binary and executing it transparently; a pure interpreter satisfies it by never compiling at all.
Whichever this language does, the latency of that first run has to feel like an interpreter's, whether or not machine code is involved under the hood.

### Scenario: a language server answering "what is the type of this expression" while the user is still typing

The wishlist commits to fast, incremental compilation and to a language server exposing semantic queries as a library API.
Editing one function in a large file must not force re-elaborating the whole module, and the answer has to come back fast enough to feel synchronous with typing.
This scenario cares about incremental type-checking latency specifically, which is a different axis than how fast the eventually-produced program runs.

### Scenario: compile-time code execution producing a value baked into the binary

The wishlist lists compile-time code execution as the metaprogramming model, with Zig `comptime` and Rust `const eval` as prior art.
That requires *some* execution engine to run arbitrary language code during compilation itself, before any notion of "the program's runtime" exists.
Whatever answers the main compilation-model question has to also answer what executes comptime code, and whether that is the same engine that will later run the compiled program or a separate one.

## Rough shape

Five candidate directions, sketched.
As with the memory strategy record, these are not necessarily exclusive, and syntax below is illustrative only.

### A. Pure ahead-of-time to native code

Every build lowers all the way to machine code before anything runs; `lang run` is sugar for "compile to a temp binary, then execute it," matching `go run`.

```
$ lang run hello.lang
# compiles to a temp binary, execs it, cleans up
```

Buys: one execution semantics for the whole language, predictable startup and steady-state performance, a light story for the C FFI and embeddability wishlist items, and it is a direct continuation of the inkwell/LLVM path already in the repo.
Costs: comptime execution needs its own answer since there is no runtime yet to run comptime code *in*, LLVM codegen latency sits on the critical path of the no-ceremony single-file scenario unless compilation itself is fast, and the language-server incremental-query scenario gets no help from this choice alone since codegen speed is not the same as type-checking speed.

### B. Bytecode interpreter, no native codegen

Source compiles to a compact bytecode; a VM interprets it.
No machine code is ever emitted.

```
$ lang run hello.lang
# compiles to bytecode in memory, interprets it directly
```

Buys: the simplest possible answer to the no-ceremony scenario, a natural home for comptime execution (comptime code just runs on the same VM before "real" execution begins, as Crafting Interpreters' `clox` demonstrates for a much simpler language), and the cheapest embeddability story since a bytecode VM is a small, portable artifact.
Costs: steady-state performance is bounded well below native code, systems-programming use cases and the C FFI boundary get more awkward when the caller lives inside a VM rather than as a native symbol, and this abandons the inkwell/LLVM work already done.

### C. Bytecode plus a tiered JIT

Bytecode runs interpreted at first; hot paths get compiled to native code at run time, the JVM/CLR/Julia shape.

```
$ lang run hello.lang
# interprets immediately, promotes hot functions to native code as they run
```

Buys: startup latency close to an interpreter's with steady-state performance approaching native code for hot paths, and Julia's "just barely ahead of time" model shows the same engine can serve comptime-style execution too.
Costs: the heaviest engineering lift of any option by a wide margin, a runtime component that most directly collides with a "no runtime" answer to the memory-strategy record, latency spikes at first invocation of any given function (documented at length as Julia's core latency problem), and W^X/codesigning restrictions on some deployment targets complicate emitting and executing code at run time.

### D. Interpreter for development and comptime, native codegen for release builds

One frontend, two backends sharing the same semantics: an interpreter (or bytecode VM) used for `lang run`, the REPL, and comptime evaluation; the LLVM path used only when the user asks for a release build.

```
$ lang run hello.lang      # interpreted, starts immediately
$ lang build --release     # AOT via LLVM, produces a native binary
```

Buys: the no-ceremony scenario gets interpreter-speed startup, comptime execution has an obvious home in the interpreter, and the release path keeps the native-codegen investment already made.
Costs: two execution semantics that must agree on every observable behavior, which is a standing conformance burden, not a one-time cost, and Zig's own comptime design explicitly requires compile-time execution to be hermetic and to behave identically to runtime execution, which is exactly the property two independently maintained backends are most likely to violate under change.

### E. Native codegen for execution, a separate incremental type-checking layer for tooling

Programs always run via AOT native codegen; the language server's incremental-query needs are served by a second, lighter-weight analysis pass (type inference and name resolution only, no codegen) that never produces an executable artifact.
This is the shape rust-analyzer uses relative to `rustc`: an independent, salsa-based incremental query engine, not a faster invocation of the real compiler's backend.

```
$ lang run hello.lang        # full AOT pipeline
# meanwhile, the editor's language server re-elaborates only the edited
# function's types, without invoking codegen at all
```

Buys: decouples "is codegen fast" from "does editing feel responsive," which the wishlist's incremental-compilation goal is really asking for, and keeps a single execution semantics, avoiding option D's dual-backend agreement burden.
Costs: does not by itself improve the no-ceremony single-file scenario, since running the program still means paying for full codegen; would likely need to be paired with A or D's answer to that scenario rather than standing alone as a complete answer to this record.

## Prior art

- **AOT vs. JIT, general tradeoffs.**
  Empirical comparisons consistently find AOT wins on startup latency and predictable memory footprint, while JIT wins on peak steady-state throughput once compilation cost is amortized over a long-running process ([Ahead-of-Time vs. Just-in-Time Compilation Trade-offs: Empirical performance studies](https://www.researchgate.net/publication/396770709_Ahead-of-Time_vs_Just-in-Time_Compilation_Trade-offs_Empirical_performance_studies); [Hybrid Execution: Combining Ahead-of-Time and Just-in-Time Compilation](https://dl.acm.org/doi/10.1145/3623507.3623554), VMIL 2023).
  The relevant lesson: the "hybrid" framing in the second paper is itself evidence that treating this as a binary choice is a simplification most mature systems abandon.

- **Crafting Interpreters (Bob Nystrom): tree-walk vs. bytecode VM.**
  [Chunks of Bytecode](https://craftinginterpreters.com/chunks-of-bytecode.html) documents the concrete engineering tradeoff behind candidate B: a tree-walker is simple and slow because it chases heap-allocated AST nodes with poor cache locality, while a bytecode VM sacrifices some of that simplicity for a dense, linear instruction stream that is far more cache-friendly, at a cost well short of native code.
  Directly relevant to candidates B and D, and the book is explicit that reaching for even more speed than bytecode buys is what a JIT is for, at the cost of candidate C's added complexity.

- **Zig: `comptime` as compile-time execution of the same language.**
  [What is Zig's Comptime?](https://kristoff.it/blog/what-is-zig-comptime/) and [Things Zig comptime Won't Do](https://matklad.github.io/2025/04/19/things-zig-comptime-wont-do.html) describe comptime execution as hermetic, reproducible, and required to behave identically to runtime execution, with no separate macro syntax.
  This is the sharpest existing example of a language answering this record's third scenario, and it is direct evidence for candidate D's central risk: Zig gets away with "the same language runs at compile time and run time" in part because it has one execution semantics, the exact property two-backend designs like D put under strain.

- **Julia: a "just barely ahead of time" JIT with precompilation.**
  Julia's JIT compiles each method specialization the first time it is called, which produces well-documented first-call latency; precompilation caches the type-inference stage to disk to blunt it ([Tutorial on precompilation](https://julialang.org/blog/2021/01/precompile_tutorial/); [Julia's latency: Past, present and future](https://viralinstruction.com/posts/latency/)).
  Relevant to candidate C as the clearest real-world account of what "startup latency close to an interpreter's" actually costs to achieve, and as a caution that the cost does not go away, it moves to a caching and invalidation problem instead.

- **rust-analyzer: an incremental query engine independent of the compiler's codegen backend.**
  rust-analyzer answers type-of and go-to-definition queries through its own salsa-based incremental analysis, entirely separate from invoking `rustc`'s actual code generation.
  Cited already in the wishlist's own prior art for fast incremental compilation, and it is the direct precedent for candidate E: the tool that makes editing feel responsive does not have to be the same pipeline that produces the executable.

Where the search was thin: nothing found treats "which engine runs comptime code" as a first-class axis of the compilation-model decision in its own right, most language design writing treats metaprogramming execution as a detail of whichever compilation model was already chosen rather than a constraint on choosing it.
That gap is one reason this record lists it as an open question rather than folding it into the general AOT/JIT tradeoff discussion above.

## Open questions

1. Does the language commit to one execution semantics or two?
   If two (candidate D's shape), what mechanism keeps an interpreter and a native-codegen backend observably identical as the language evolves, given that Zig's comptime design treats exactly this identity as load-bearing?
1. What engine executes comptime code: the same engine that runs the compiled program, a separate interpreter used only at compile time, or a restricted subset that itself compiles to native code before running?
1. Is the no-ceremony single-file scenario satisfied by "compile to a temp artifact and exec it" (candidate A/E), or does it require true interpretation with no codegen step at all (candidate B/D)?
   These have different latency floors, and the wishlist entry does not say which is acceptable.
1. Is the language server's incremental-query latency goal served by making the execution backend itself incrementally fast, or by a separate type-checking-only layer that never touches codegen (candidate E)?
   rust-analyzer's precedent suggests these can be decoupled; this record has not decided whether they should be.
1. Is a runtime component (a scheduler, a GC, a JIT's own machinery) acceptable at all?
   This bounds the candidate set before any other question here is answered, and it is the same open question the memory-strategy record asks; the two records currently ask it independently rather than sharing an answer.
1. Does embeddability (hosting the language inside another program, per the wishlist) require the execution engine itself to be a small, linkable artifact, which pulls toward B/D, or is embeddability satisfiable by exposing only the AOT-compiled output as a linkable library?
1. **Blocked on an undecided foundational decision:** the memory-management strategy record (`0001`) is itself blocked partly on this record, since its candidates A and C presume a runtime component whose cost depends on the answer here.
   The two records currently reference each other's open status without either committing first; resolving that ordering is not this record's decision alone to make.
1. **Blocked on an undecided foundational decision:** the paradigm has no decision record.
   Julia's JIT design leans heavily on multiple dispatch to decide what to specialize; a paradigm with different dispatch or genericity mechanics changes what a JIT (candidate C) or a comptime-driven monomorphization scheme (candidates A/D) actually has to do.

## Advancement record

- 2026-08-27, gate 0 → 1: sketched from a new wishlist entry ("A compilation model chosen as an explicit decision") added under Foundational decisions; five candidate directions surveyed, prior art cited across AOT/JIT tradeoffs, tree-walk vs. bytecode VMs, Zig comptime, Julia's JIT, and rust-analyzer's incremental architecture; open questions recorded including two blocks on undecided foundational decisions and a two-way block with the memory-strategy record (`0001`).

## Changelog
