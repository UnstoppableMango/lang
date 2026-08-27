# Paradigm: answers to the unblocked open questions

Answers given by the author on 2026-08-27 to open questions 1-5 in `docs/design/0003-decision-paradigm.md`.
Questions 6-8 are blocked on the memory-strategy, compilation-model, and concurrency decisions, and were not asked.

This is a note, so it carries zero standing.
The answers are not part of the design record until they are written into the design doc itself, which is stage 2 work.
This file exists so that a real wishlist entry has something to draw on.

## The answers

1. **Is the AST the paradigm's central case:** yes.
   Judge every candidate by how well it models the exact recursive-sum-shaped value `hack/hello.lang` already compiles through today.
1. **Mutation by default:** immutable by default.
   An explicit marker is required at a binding to mutate it.
1. **One idiom or two for attaching behavior to data:** one mechanism, two spellings.
   A method call is sugar over a free function that takes the value as its first argument, Rust's `x.foo()` desugars to `Type::foo(x)` shape.
1. **Retroactive interface conformance:** pure structural satisfaction, Go-shaped, no declared association at all.
   A free function counts as one of a type's methods by having that type as its first parameter, the same substrate answer 3's method-call sugar already presupposes, and a type satisfies an interface whenever its method set matches, with no `impl`-shaped block naming the relationship.
   Revised from an initial explicit-`impl` answer after the interview surfaced that it contradicted the wishlist's existing structural-interfaces commitment; see the first resolved tension below.
1. **Is paradigm scoped to organization and dispatch only:** yes.
   The type system's nominal-versus-structural axis and the memory record's value-semantics candidate remain separate decisions this record does not presume.

## What this rules out

Everything below is derivation, not authored answer, and it is the part most likely to be wrong.

- **Candidate A, object-oriented with class inheritance: dead.**
  Answers 3 and 4 together choose free-function dispatch with method-call sugar plus explicit trait-shaped `impl` blocks, not a class hierarchy with `extends` and virtual dispatch.

- **Candidate C, procedural Go-shaped: dead as a whole, but its interface mechanism survives.**
  Answer 2 rejects unremarkable mutation outright, which was C's other defining trait.
  Answer 4, after revision, is exactly C's structural, no-declaration interface satisfaction, just running over a free-function substrate instead of Go's receiver methods.
  C does not come back as a candidate; this one piece of it is folded into B.

- **Candidate E, data-oriented as the primary paradigm: dead as a primary model.**
  Answer 1 makes a recursive sum type the central case, which is exactly the shape a batch/columnar layout fights hardest, the same conclusion the memory-strategy record reached about arenas as the primary model.
  Nothing here rules out a data-oriented layout as a compiler-internal optimization for bulk passes; it is just not the language's organizing idea.

- **Candidate D, fusion, Scala-shaped: absorbed rather than chosen as sketched.**
  D's central insight, don't force a single idiom for attaching behavior to data, survives, but in a narrower form than D originally sketched: one mechanism (free functions) with method syntax as sugar over it, not two independently first-class idioms living side by side.
  D's "immutability encouraged but not enforced" is also superseded by answer 2's firmer default.

- **Candidate B, functional, algebraic data types and pattern matching over free functions: strongest survivor, effectively chosen,** with one amendment neither original sketch of B anticipated: method-call sugar over free functions (answer 3, absorbed from D).
  Answer 4's structural conformance is not really an amendment to B; it is C's interface story, already on the wishlist, applied to B's free-function substrate.

## The shape this suggests

Algebraic data types and pattern matching are the primary way to model closed, recursive data, the AST scenario being the concrete case.
Values are immutable by default; an explicit marker is required to mutate a binding.
Behavior lives in free functions; a method call is sugar over a free function taking the value first, one mechanism, not two.
A function counts as one of a type's methods by taking that type as its first parameter, and an interface is satisfied purely by a matching method set, structurally, with no declared relationship anywhere, the wishlist's existing Go-shaped commitment carried over onto a free-function substrate.

That is closer to a real decision than a rough shape.
It is not yet a design, and the tensions below have to be resolved before it becomes one.

## New tensions the interview created

Each one is a candidate open question for the stage 2 pass.

1. **This converges on most of Rust's paradigm, minus Rust's own interface mechanism, and nobody said so out loud during the interview.**
   Immutable-by-default, ADTs plus pattern matching, and free-function dispatch with method-call sugar reads like Rust, right up until interface conformance, which is deliberately Go-shaped instead of trait-shaped.
   Worth being honest about at the stage 2 pass: name Rust explicitly for the parts that converge, and be explicit that structural interfaces are the one place this record diverges from it on purpose.

1. **Resolved: answer 4 initially chose Rust-shaped explicit `impl` blocks, which conflicted with the wishlist's existing "structural interface satisfaction, checked statically" entry.**
   `docs/notes/duck-typing.md` places Go's implicit, no-declaration interfaces as structural and explicit `impl Trait for Type` blocks as "nominal + static, mostly" in its own four-quadrant framing, so the initial answer presumed an answer to the nominal-versus-structural axis that answer 5 said this record must not presume.
   Revised on the author's direction back to pure structural satisfaction: a function is a type's method by taking that type first, and conformance is inferred from the method set alone, no declaration anywhere.
   This keeps the record consistent with the existing wishlist entry and with answer 5's scoping.

1. **The revised answer 4 trades Rust's coherence guarantee for Go's collision risk.**
   With no declared association, two unrelated functions with the same name and a matching first-parameter type can collide, exactly the ambiguity Rust's explicit `impl` (and its orphan rules) exist to prevent.
   Go's answer to this in practice is convention and a single package owning a type's methods; whether this language gets the same protection for free, or needs its own answer, is unexamined.

1. **Answer 2's "explicit marker to mutate" and the memory record's "regime is a property of the type's definition" have an unexamined interaction.**
   Does mutability become a per-binding annotation independent of which memory regime a value is in, or does the second, reference-counted regime (for shared or graph-shaped values) imply always-mutable-in-place semantics, since shared graph structures often need interior mutability regardless of the outer binding's marker?

1. **Answer 1 is a design aspiration proven only by an implementation in a different, already-decided paradigm.**
   The reference compiler is Rust, not self-hosted, and has no plan recorded anywhere to become self-hosted.
   Whether "the AST is the central case" needs the compiler to eventually demonstrate that in the target language itself, or stays a design principle judged only on paper, is unaddressed.

1. **Candidate E's fate here parallels arenas' fate in the memory-strategy record exactly.**
   Both were rejected as the *primary* organizing idea and both remain available only as a compiler-internal optimization for bulk, phase-shaped work.
   Whether that parallel is meaningful, the two rejected candidates describing the same underlying idea from two different records, or coincidental, is worth asking directly rather than tracking as two independent footnotes.

## Answers to the tensions

Answers given by the author on 2026-08-27 to four of the tensions raised above.

1. **Collision protection:** restrict a type's method set to functions defined in the same module as the type, Go's actual rule.
1. **Mutation and memory regime:** independent.
   The mutability marker governs mutability regardless of which memory regime a value is in; a reference-counted value can still be declared immutable.
1. **Self-hosting:** yes, eventually.
   The paradigm should be judged partly on what it would take to rewrite the compiler in its own target language someday.
1. **Naming Rust explicitly in the stage 2 doc:** no.
   Keep the candidate analysis framed on its own terms rather than naming Rust as the target; let the convergence stay implicit in the candidates and their costs.

Tension 6, whether candidate E's fate paralleling arenas' fate in the memory-strategy record is meaningful or coincidental, was not asked; it is a reflective question for whoever does the stage 2 pass, not one with a clean yes/no shape.

### What the tension answers narrow further

The collision-protection answer turns out not to cost anything against the duck-typing note's own motivating scenario.
That scenario is retroactive conformance to a *new interface* (`Logger`, defined after the fact, in your own module) over a type's *existing* methods (`Widget.Log`, already defined in `Widget`'s own module), never retroactively adding a *new method* to someone else's type from outside it.
Restricting a type's method set to its own defining module costs nothing there; it only forecloses a scenario the language was never trying to support in the first place, adding a new method to a type you don't own.
This is worth stating plainly at the stage 2 pass rather than leaving it as an apparent tension between "structural, retroactive conformance" and "methods are module-scoped."

## Not asked

Questions 6, 7, and 8 remain blocked: the memory-strategy record (`0001`), the compilation-model record's dispatch-mechanics question (`0002`), and the concurrency model, which has no decision record.
Answer 2 (immutable by default) directly matches the shape the memory-strategy interview already leaned toward without either interview referencing the other, and answers 3 and 4's free-function dispatch with structural conformance is exactly what the compilation-model record's blocked question 8 was waiting on: Julia-style multiple dispatch is not in this picture at all, so whatever candidate C or a comptime monomorphization scheme does in that record, it does it over single-dispatch, structurally-typed free functions.
