# Memory strategy: answers to the unblocked open questions

Answers given by the author on 2026-08-26 to the seven unblocked open questions in `docs/design/0001-decision-memory-strategy.md`.
Questions 8, 9, and 10 are blocked on the compilation model, paradigm, and concurrency decisions, and were not asked.

This is a note, so it carries zero standing.
The answers are not part of the design record until they are written into the design doc itself, which is stage 2 work.
This file exists so that work has something to draw on.

## The answers

1. **One mechanism or two:** two, and which regime a value is in is visible in the source.
1. **Which scenario may be awkward:** phase-shaped allocation.
   The long-lived ownerless graph and non-memory resource release both have to be good.
1. **Is a runtime acceptable:** no runtime at all.
1. **How safety is enforced:** compile-time first, with a runtime check where the static check cannot prove safety.
1. **Deterministic cleanup of non-memory resources:** same mechanism as memory, opt-in per type.
   A type declares it has a release action and the memory model runs it at reclamation.
1. **Annotation budget:** under roughly 10% of function signatures, and only where a value crosses a module or FFI boundary.
1. **C FFI boundary:** pin or transfer.
   Values handed to C move into a regime with a stable address and a lifetime the programmer controls explicitly.

## What this rules out

Everything below is derivation, not authored answer, and it is the part most likely to be wrong.

- **Candidate A, tracing garbage collection: dead.**
  Answer 3 removes it outright, and answer 5 would have removed it anyway, since a collector cannot run a release action at a point the programmer can predict.

- **Candidate D, regions or arenas as the *primary* model: rejected on 2026-08-26,** recorded in the design doc's Changelog.
  Answer 2 is the decisive one.
  Regions are the phase-shaped answer, so naming phase-shaped allocation as the case allowed to be awkward says directly that regions are not what the language is organized around.
  Worth flagging plainly: `docs/notes/arena-memory-model.md` is the note that motivated opening this decision record, and this interview pushed against its premise.
  Regions survive as a possible second regime or as an optimization, not as the spine.

- **Candidate E, mutable value semantics: weakened.**
  It is the option that is worst at the long-lived ownerless graph, and answer 2 says that scenario must be good.

- **Candidate C, compiler-inserted reference counting: strongest survivor.**
  It needs no tracing runtime, it reclaims deterministically so answer 5 works, it is good at graphs, and the place it hurts is exactly high-churn phase-shaped allocation, which answer 2 has already accepted as the awkward case.
  Every part of the interview points here.

- **Candidate B, ownership and borrowing: alive as one of the two regimes,** but see the annotation tension below.

- **Candidate F, generational references: alive as the runtime fallback in answer 4,** but see the runtime tension below.

## The shape this suggests

A default regime that is statically checked and costs nothing at run time, plus a second, explicitly marked regime for values that are genuinely shared or graph shaped, reclaimed by compiler-inserted reference counting, with a per-type release action.
Answer 4 then reads as: the count *is* the runtime fallback, and it applies only inside the second regime.

That is a coherent language.
It is not yet a design, and four things have to be resolved before it becomes one.

## New tensions the interview created

These are the useful output.
Each one is a candidate open question for the stage 2 pass.

1. **Answer 3 and answer 4 need "runtime" defined before they agree.**
   "No runtime at all" and "a runtime check where the static check cannot prove safety" are only compatible under a specific reading: emitted inline instructions are not a runtime, a heap-scanning collector is.
   Where a reference count table, a generation table, or an allocator support library falls on that line is undecided, and answer 4 cannot be implemented until it is.

1. **Answer 1 and answer 6 pull against each other.**
   If the regime a value belongs to must be visible in the source, and the annotation budget is under 10% of signatures and boundaries only, then regime membership cannot be spelled in every signature that touches the value.
   The available resolution: make the regime a property of the *type definition* and the construction site, not of every signature that passes the value along.
   That is a real design idea and it should be tested rather than assumed.

1. **Answer 6 is a hard constraint on candidate B.**
   Whether Rust clears a 10% bar depends entirely on whether `&` and `&mut` count as memory annotations.
   If they do, Rust is far over budget and the default regime cannot be borrow-checker-shaped in the usual sense.
   The stage 2 pass should answer the counting question first, because it decides whether B is viable at all.

1. **Cycles are now the sharpest unanswered problem.**
   Reference counting leaks cycles, and answer 3 forbids the tracing collector that would normally clean them up.
   The remaining options are: leak them and say so, forbid constructing them in the type system, or provide a weak-reference-shaped tool and make cycles the programmer's problem.
   None of these is obviously right, and the long-lived registry scenario in the design doc is exactly where cycles show up.

1. **Answer 7 constrains the default regime.**
   Pin or transfer requires a stable address, so the language must be able to promise that a value does not move.
   Whether the default regime is allowed to relocate values is now a question, and it was not one before.

## Not asked

Questions 8, 9, and 10 remain blocked.
Nothing above should be read as deciding the compilation model, the paradigm, or the concurrency model, though answer 3 leans hard toward ahead-of-time compilation and the shape above leans toward mostly-immutable data.
Those leanings are pressure on the blocked decisions, not answers to them.
