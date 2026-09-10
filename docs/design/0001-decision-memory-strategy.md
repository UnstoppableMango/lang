---
title: 'Decision: Memory management strategy'
stage: 1
status: active
created: 2026-08-26
updated: 2026-08-26
---

# Decision: Memory management strategy

## Motivation

Every value a program creates occupies memory that must eventually be reclaimed.
The language has to answer three questions about that reclamation, and it has to answer them the same way everywhere:

1. Who decides when a value's memory is released: the programmer, the compiler, or a runtime?
1. What prevents a program from reading memory that has already been released?
1. What does that prevention cost, and is the cost paid at compile time, at run time, or in what the programmer is allowed to write?

This is a decision record rather than a feature because the answer is load bearing for most of the wishlist.
Deterministic scope-tied cleanup, data-race freedom, structured concurrency, a C-compatible FFI boundary, and "errors as values" all read differently depending on the answer, and none of them can be designed while it is open.

Deferring the decision is not neutral.
A language accretes assumptions about memory in its syntax and its standard library long before it admits to having made a choice, and unwinding those assumptions later is the expensive kind of mistake.
Making the choice explicit, in writing, with the rejected options recorded, is the point of this document.

### Scenario: phase-shaped allocation

A compiler pass parses a source file into an AST, walks it, produces an IR, and drops the AST.
Every allocation made during the parse has exactly one lifetime, and that lifetime is the pass.
The language should make this case cheap and obvious.
Individually tracking the lifetime of each AST node is wasted work when the entire graph dies at one known point.

### Scenario: a long-lived graph with no single owner

A running program holds a registry: a cache of interned strings, a symbol table, a set of open connections, each entry referenced from several places, with no one place that obviously owns it.
Entries are added and removed while the program runs.
The lifetime of any single entry is genuinely dynamic and not derivable from lexical structure.
The language should make this expressible without the programmer contorting the data structure to satisfy the memory model.

### Scenario: releasing something that is not memory

A file handle, a socket, or a lock has to be released at a point the programmer can predict, not "eventually".
Whatever answers the memory question also constrains the answer here, because a strategy with non-deterministic reclamation cannot by itself provide deterministic resource cleanup.

These three scenarios pull in different directions on purpose.
A strategy that only serves the first is a domain-specific tool, not a language.

## Rough shape

Six candidate directions, sketched.
They are not mutually exclusive in the limit, and the question of whether the language picks one or layers two is itself open.
Syntax below is illustrative only.
No syntactic decision is implied, and none of these sketches should be read as a preferred surface.

### A. Tracing garbage collection

A runtime periodically determines which values are reachable and reclaims the rest.
The programmer writes no annotations and thinks about lifetimes only for performance.

```
fn build() -> Tree {
    let t = Tree::new()   // no annotation, no free
    t                     // returning it is fine, the collector figures it out
}
```

Buys: the simplest programmer model on offer, and the long-lived-graph scenario becomes free.
Costs: a runtime, non-deterministic reclamation (so scope-tied cleanup of non-memory resources needs a separate mechanism), pause behavior to manage, and a heavier FFI boundary because foreign code holding a pointer must be visible to the collector.

### B. Ownership and borrowing

Each value has one owner.
References are checked at compile time to not outlive the value they point to.

```
fn build(input: &str) -> Tree { ... }   // borrows, does not own
fn consume(t: Tree) { ... }             // takes ownership, t dies here
```

Buys: no runtime, deterministic reclamation, and the checker doubles as the data-race story.
Costs: the largest learning burden of any option, real programs that are correct but rejected, and the graph scenario becomes the case where the model fights hardest.

### C. Compiler-inserted reference counting

The compiler emits the increments and decrements, rather than a runtime tracing the heap.
Perceus shows this can be made precise enough to be garbage free for cycle-free programs, and to enable in-place reuse.

```
fn map(xs: List<A>, f: A -> B) -> List<B>   // compiler may reuse the input cells in place
```

Buys: deterministic reclamation without a tracing runtime, a light FFI boundary, and no annotation burden.
Costs: cycles leak unless something else handles them, count traffic costs throughput, and the design leans on the language being mostly immutable, which presumes decisions this project has not made.

### D. Regions or arenas as the primary model (rejected 2026-08-26)

Allocations go into a region, and the region is released as a unit.

```
region r {
    let ast = parse(r, src)
    lower(r, ast)
}   // every allocation in r dies here
```

Buys: the phase-shaped scenario becomes both the cheapest and the most obvious thing to write, allocation is a pointer bump, and release is one operation.
Costs: this does not answer "does this reference outlive its region", which is the same question a borrow checker exists to answer, and the long-lived graph scenario has no home unless a second mechanism exists.

**Rejected as the primary model.**
The whole argument for organizing a language around regions is that phase-shaped allocation is the case worth making cheapest and most obvious.
This record's answer to open question 2 gives that case away: phase-shaped allocation is the scenario allowed to be awkward, and the long-lived ownerless graph is one that must be good.
Regions are strongest exactly where the language has decided it can afford to be weak, and weakest exactly where it has decided it cannot.
That removes the reason to build the language around them.

This rejects regions as *the spine*, not as a mechanism.
A region-shaped construct remains available as a candidate for the second regime, or as an optimization the compiler applies without the programmer naming it.
`docs/notes/arena-memory-model.md` explored this candidate at length, reached the same conclusion from the other direction (that arenas alone cover only part of the map), and is retained as the record of that exploration.

### E. Value semantics with second-class references

Values are independent, and references exist only implicitly at call boundaries, never stored in a variable or a field.
Because nothing shares mutable state, there is nothing to collect.

```
fn update(inout t: Tree, x: Item) { ... }   // reference exists for the call only
```

Buys: no runtime, no counts, no lifetime annotations, and data-race freedom falls out of the same rule.
Costs: data structures that are genuinely graph shaped cannot be written directly, which hits the second scenario hardest of all six.

### F. Runtime-checked references

Each allocation carries a generation number, each reference remembers the generation it was made from, and a dereference checks that the two still match.

```
let p = &node        // remembers node's generation
use(p)               // checks generation, aborts on mismatch
```

Buys: no borrow checker, no counts, no tracing, no annotations, and use-after-free becomes a deterministic abort rather than undefined behavior.
Costs: a check on dereference, per-allocation space overhead, memory not returned to the OS as promptly, and safety enforced by aborting a running program rather than by rejecting it at compile time.

### The hybrid question

Real systems already combine these: a per-request arena over a collected heap, or refcounting with a cycle collector.
The open question is not whether a hybrid is possible but whether the language exposes one mechanism or two.
Two mechanisms means the programmer must always know which regime a value is in, and that has to be visible in the source rather than inferred from context.

## Prior art

- **Tofte and Talpin, region inference (1997), and Cyclone.**
  Cyclone added a region type system to C, then found regions alone were not enough and layered unique pointers and reference counted objects on top, reported in "Experience with Safe Manual Memory Management in Cyclone" (ISMM 2004, [DOI 10.1145/1029873.1029883](https://dl.acm.org/doi/10.1145/1029873.1029883); papers collected at [cyclone.thelanguage.org](https://cyclone.thelanguage.org/wiki/Papers/)).
  The lesson to take is the layering itself: the project that tried hardest to make regions the whole answer concluded it needed more than one mechanism.
  Cyclone's region and lifetime work is also the most direct ancestor of Rust's borrow checker, so candidates B and D share a root.

- **Rust: ownership and borrowing.**
  The existence proof that compile-time lifetime checking can carry a systems language with no runtime, and that the same analysis can be reused for data-race freedom.
  The cost is equally well documented: [The Usability of Advanced Type Systems: Rust as a Case Study](https://arxiv.org/pdf/2301.02308) and [Garbage Collection Makes Rust Easier to Use](https://arxiv.org/pdf/2110.01098) both measure the learning burden, the latter with a randomized controlled trial in which a garbage collector measurably improved outcomes.
  Directly relevant to a wishlist that lists "hard to do the wrong thing by accident" next to "pit of success".

- **Koka and Perceus.**
  "Perceus: Garbage Free Reference Counting with Reuse" (PLDI 2021, [paper](https://xnning.github.io/papers/perceus.pdf), [extended version](https://www.microsoft.com/en-us/research/wp-content/uploads/2020/11/perceus-tr-v1.pdf)) shows compiler-inserted counting can be precise rather than conservative, and that garbage freedom unlocks in-place reuse.
  The important caveat for this project: the result leans on Koka being mostly immutable with a strong effect system, which is a paradigm decision this language has not made.

- **Project Verona.**
  Ownership over *groups* of objects rather than single objects, organizing the heap into a forest of isolated regions ([project site](https://microsoft.github.io/verona/), ["Reference Capabilities for Flexible Memory Management"](https://arxiv.org/pdf/2309.02983), OOPSLA 2024).
  This is the closest known art to the "arenas plus a checked escape rule" shape in `docs/notes/arena-memory-model.md`, and it is evidence that the shape is tractable, at the cost of a reference capability system.
  Verona also fuses memory and concurrency into one model, which is the same fusion the note guesses at.

- **Hylo (formerly Val): mutable value semantics.**
  [Introduction to Hylo](https://hylo-lang.org/introduction/) and ["Implementation Strategies for Mutable Value Semantics"](https://arxiv.org/pdf/2106.12678) argue that making references second class removes the need for a collector or reference counting outright.
  Worth taking seriously as the option that solves the problem by deleting it, and worth being honest that it is also the option that most restricts what data structures can be written.

- **Vale: generational references.**
  ["Vale's Memory Safety Strategy: Generational References and Regions"](https://verdagon.dev/blog/generational-references) documents a runtime-checked approach with published benchmarks against reference counting, and it deliberately combines the runtime check with regions for the cases where the check can be elided.
  The relevant idea is not the technique alone but the pairing: a cheap default plus a static analysis that removes the check where it can be proven unnecessary.

- **Go: a tracing collector as a design constraint.**
  Go treats collector latency as a language-level promise rather than an implementation detail, which is the counterexample to "GC means unpredictable".
  Relevant because Go's simplicity stance already appears in the wishlist, and the collector is a large part of what buys that simplicity.

Where the search was thin: nothing found makes a language-level commitment to arenas as the *primary* model with a narrow escape check as the only static analysis, which is option 3 in `docs/notes/arena-memory-model.md`.
Verona is the closest, and it differs by generalizing to a full reference capability system rather than staying with a single narrow check.
That gap is recorded as evidence of the search, and it is now moot for this record: candidate D is rejected as the primary model on its own merits, not for want of prior art.
The gap would become live again only if a later record revisits regions as a spine.

## Open questions

1. Is the language's memory strategy one mechanism or two?
   If two, what makes the regime of a given value visible in the source?
1. Which scenario is allowed to be the awkward one?
   Every candidate is bad at one of the three motivating scenarios, and no candidate is bad at none of them.
1. Is a runtime acceptable at all?
   This bounds the candidate set before any other question is answered, and it has not been decided.
1. Is memory safety enforced by rejecting programs at compile time, by checking at run time, or by neither?
   Candidate F makes this an explicit third answer rather than a spectrum between the first two.
1. Does the strategy have to provide deterministic cleanup for non-memory resources, or is scope-tied cleanup a separate mechanism layered on top?
1. What is the annotation budget?
   Concretely: what fraction of function signatures may carry memory-related annotations before the design is considered to have failed the "pit of success" goal?
1. How does the strategy behave at the C FFI boundary, given that a foreign pointer is invisible to any static analysis and to a tracing collector alike?
1. **Blocked on an undecided foundational decision:** the compilation model (ahead-of-time, interpreted, or mixed) has no decision record.
   Candidates A and C presume a runtime component; a pure ahead-of-time decision does not rule them out but changes their cost.
1. **Blocked on an undecided foundational decision:** the paradigm has no decision record.
   Candidate C's precision result depends on mostly-immutable data, and candidate E is a paradigm choice as much as a memory choice.
1. **Blocked on an undecided foundational decision:** the concurrency model has no decision record.
   Candidates B and E each fuse memory and concurrency into one mechanism, so deciding memory first may decide concurrency by accident, and the ordering of these two decision records is itself unresolved.
   Candidate D fused them too, and its rejection removes one instance of this risk without removing the risk.

## Advancement record

- 2026-08-26, gate 0 → 1: sketched from the wishlist entry "A memory strategy chosen as an explicit decision"; six candidate directions surveyed, prior art cited across regions, ownership, reference counting, value semantics, and runtime-checked references, open questions recorded including three blocks on undecided foundational decisions.

## Changelog

- 2026-08-26: candidate D, regions or arenas as the primary model, rejected.
  The author's answer to open question 2 names phase-shaped allocation as the scenario allowed to be awkward and the long-lived ownerless graph as one that must be good, which inverts the case regions are built to win.
  Regions remain available as a possible second regime or as a compiler optimization; only the spine claim is rejected.
  The doc remains at stage 1.
  Recording a rejected direction with its reason is one of the artifacts the 1 → 2 gate requires, but the gate also requires a *chosen* direction, and none has been marked.
