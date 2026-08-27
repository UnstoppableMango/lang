# Concurrency: answers to the unblocked open questions

Answers given by the author on 2026-08-27 to a preliminary clarification and open questions 1-6 in `docs/design/0004-decision-concurrency.md`.
Questions 7-9 are blocked on the memory-strategy, paradigm, and compilation-model decisions, and were not asked.

This is a note, so it carries zero standing.
The answers are not part of the design record until they are written into the design doc itself, which is stage 2 work.
This file exists so that a real wishlist entry has something to draw on.

## Preliminary clarification: what "no runtime at all" meant

Both the memory-strategy and compilation-model interviews answered "no runtime at all," and this record's open question 1 asked whether concurrency is bound by that same answer.
The author clarified that "no runtime" in both prior interviews meant no separate binary to install and manage on a host, not zero scheduling or GC-shaped machinery inside the compiled artifact.
A runtime component is acceptable here as long as it ships as part of the single compiled binary.
This directly answers open question 1: yes, a runtime is acceptable under that definition.

This clarification is a live loose end for `docs/notes/compilation-model-interview-answers.md`, not resolved in this note.
That interview's candidate C, a bytecode VM plus a tiered JIT, was killed specifically for being "a runtime component by definition," under what may have been the stricter reading.
Whether a JIT that ships fully inside the binary survives under this looser definition is an open question for that record, flagged here, not answered here.

## The answers

1. **Is a runtime acceptable:** yes, if it ships as part of the single compiled binary. See the clarification above.
1. **Which of the three wishlist commitments is allowed to be awkward:** not asked directly; see "what this rules out" below, none of the three ended up forced.
1. **Same mechanism as the memory record, or a third discipline:** the same mechanism.
   Concurrency-safety extends the memory record's existing ownership-and-borrowing regime across a thread boundary with a trait-bound-shaped check at the crossing point, Rust's `Send`/`Sync` shape, rather than introducing a second, independent capability system.
1. **Structured concurrency, keyword or library:** keyword-level from day one.
   `scope`/`spawn`-shaped constructs are baked into the language rather than left to a library convention the way Go's `context.Context` was.
1. **Mutexes and atomics, sanctioned or not:** both available, alongside channels.
   Rust's actual answer: a lock-shaped primitive exists for the many-reader, shared-state case (a connection pool, a cache), channels or capability-checked sharing cover the common case, and the language accepts deadlock as a separate, known cost the ownership regime does not close off.
1. **High-fan-out workload target:** yes.
   Tens or hundreds of thousands of concurrent in-flight tasks is a target workload, which rules out OS-threads-only as a sufficient answer on its own.

## What this rules out

Everything below is derivation, not authored answer, and it is the part most likely to be wrong.

- **Candidate A, OS threads plus an async multiplexer, no green threads: dead as the primary model.**
  Answer 6 names high fan-out as a target workload, exactly the scenario the design doc's own motivation names as A's specific, checkable failure point.
  Real OS threads do not disappear entirely; they remain the M underneath whatever M:N scheduler answer 6 now requires.

- **Candidate C, stackless async/await: dead,** as it already was from the initial sketch.
  Nothing in this round of answers revives it; it still violates the no-coloring wishlist commitment by construction.

- **Candidate D, actor model, process isolation: dead as the primary model.**
  Answer 3 chose extending shared-memory ownership across threads over Erlang's no-sharing-at-all answer.
  D's failure-isolation story, "let it crash" supervision, has no home in the chosen shape and was not asked about; see the new tensions below.

- **Candidate B and candidate E did not resolve into "one wins, one dies."**
  Answer 3 chose E's policy: static, compile-time aliasing control extending the memory record's regime, not a runtime-checked or unchecked shared-memory free-for-all.
  But answer 6's high-fan-out requirement, combined with the no-coloring wishlist commitment already made, requires a real stack per task, since a coloring-free lightweight task is exactly what stackless models cannot provide and OS threads cannot provide cheaply at that scale.
  B is the only candidate sketched that supplies that mechanical shape.
  The result is a hybrid neither candidate described alone: **B's stackful, M:N scheduling mechanics carrying E's static aliasing discipline**, Go's scheduler and growable stacks with Rust's `Send`/`Sync`-shaped compile-time check at the task boundary instead of Go's anything-goes shared memory.

## The shape this suggests

Tasks are stackful and scheduled M:N onto a pool of OS threads, a real stack per task, so any function can block without a marker, satisfying the no-coloring commitment the same way Go does.
Moving or sharing a value across a task boundary is checked at compile time by extending the memory record's existing ownership-and-borrowing regime, a `Send`/`Sync`-shaped bound at the crossing point rather than a new capability system, satisfying the race-freedom commitment without a runtime check and without giving up zero-copy sharing.
`scope`/`spawn` are language keywords from day one, so a child task cannot outlive its parent's scope, satisfying the structured-lifetime commitment without a library-level patch.
Lock-shaped primitives and channels both exist; the language accepts deadlock, not race, as the resulting known-but-unclosed hazard.

Read against the design doc's own open question 2, none of the three named wishlist commitments had to become the awkward one.
The actual residual cost, deadlock, was never one of the three named commitments in the first place, which is worth being explicit about at the stage 2 pass rather than letting the three-for-three result read as a free lunch.

That is closer to a real decision than a rough shape.
It is not yet a design, and the tensions below have to be resolved before it becomes one.

## New tensions the interview created

Each one is a candidate open question for the stage 2 pass.

1. **The chosen shape is a hybrid the rough-shape candidates did not name.**
   Stage 2 should either name it as its own candidate, distinct from both B and E as sketched, or restructure the candidates section so the mechanics axis (how a task is represented and scheduled) and the safety axis (how sharing across tasks is checked) are presented as two separate questions from the start, since this interview answered them independently.

1. **This pairing is not proven anywhere the prior-art search found.**
   Pony's reference capabilities were built around Pony's own actor-based scheduler, and Rust's `Send`/`Sync` pair with OS threads or stackless async tasks, not with a language-level M:N green-thread scheduler.
   Grafting a `Send`/`Sync`-shaped static check onto Go-shaped stackful green threads combines two things that, as far as this record's research found, have not been combined in a shipped language, worth flagging plainly rather than assuming the combination is as well-trodden as either half is individually.

1. **Deadlock is the accepted residual hazard, and nothing here addresses it.**
   Choosing both lock-shaped primitives and channels means a lock-ordering cycle between two tasks remains a runtime-only bug, in a design that otherwise promotes memory and concurrency safety to compile time everywhere else.
   Whether this is acceptable as stated, or whether the language wants a narrower answer here too (a lock ordering discipline, or lock-shaped primitives scoped narrowly enough that cycles are structurally rare), is unexamined.

1. **Candidate D's failure-isolation story has no successor in the chosen shape.**
   "Let it crash" supervision was Erlang's answer to "what happens when one task panics while others are running," and the chosen shape does not have an equivalent.
   The wishlist's existing "recoverable errors as values, distinct from unrecoverable panics" entry is adjacent but does not by itself say what happens to sibling tasks when one panics inside a `scope`.
   This may belong to this record, to a future error-handling record, or to the `scope`/`spawn` semantics specifically; undecided.

1. **This record now depends concretely on how the memory-strategy record's stage 2 pass specifies crossing a thread boundary.**
   Answer 3 presumes a `Send`/`Sync`-shaped bound exists to extend, but `0001` has not yet specified the mechanism for checking or transferring a value across any boundary in that level of detail, boundary-crossing there was discussed only for the C FFI case.

## Not asked

Questions 7, 8, and 9 remain blocked: the memory-strategy record (`0001`), the paradigm record (`0003`), and the compilation-model record (`0002`).
Nothing above should be read as deciding any of the three, though answer 3 leans on the memory record's regime existing in a specific, not-yet-stage-2 shape, and the no-coloring result leans on the paradigm record's free-function-plus-sugar dispatch answer without depending on its exact mechanics.
