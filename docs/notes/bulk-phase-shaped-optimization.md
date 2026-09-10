# Bulk, phase-shaped optimization: arenas and data-oriented layout are one deferred idea

Two decision records independently rejected a candidate for the same reason, phrased in each record's own vocabulary, without either referencing the other at the time.

- `docs/design/0001-decision-memory-strategy.md`, candidate D, regions or arenas as the primary memory model: rejected 2026-08-26.
  The rejection reasoning names phase-shaped allocation, a batch of values that all die together at one known point, as the scenario the language is allowed to leave awkward, not the one it optimizes for.
- `docs/design/0003-decision-paradigm.md`, candidate E, data-oriented layout as the primary paradigm: rejected 2026-08-27.
  The rejection reasoning names the same shape from the opposite direction: a batch of same-typed values processed together in one pass, contrasted with the AST's recursive-sum shape as the paradigm's actual central case.

The author confirmed on 2026-08-27 that this parallel is meaningful, not coincidental: both candidates describe the same underlying idea, a batch of same-shaped values laid out together and reclaimed or processed together, arriving at it from two different records (memory reclamation in one, data organization in the other).
Tracking them as two independent rejected footnotes would lose that.

## The idea, stated once

A phase-shaped batch: values of one shape, allocated together, transformed together in bulk, and released together at one known point.
`docs/notes/arena-memory-model.md` explored the memory side of this at length before the decision record existed.
Neither record's rejection means the idea is dead.
Both mean it is not the language's *primary* organizing idea, for data or for memory, and both leave it exactly the same door: available as something the compiler does internally, or exposes as an explicit, narrow, opt-in construct, without either becoming the spine of the paradigm or the memory strategy.

## Why this stayed deferred rather than resolved now

Both records rejected this as a language-level default, not as a compiler engineering technique.
Whether it is worth building as an actual compiler-internal optimization is a question about a specific pass in a specific compiler, not about the language's design, and this project's compiler does not yet have a pass that would need it: `src/main.rs` is a small stub with no AST, no passes, and no incrementality either way, per `docs/notes/compilation-model-interview-answers.md`'s own finding when it hit the same question.

## Where to pick this back up

Revisit when the compiler has a real bulk-shaped pass, most plausibly type-checking or lowering many AST nodes at once, per `docs/design/0003-decision-paradigm.md`'s own bulk-transformation scenario, and test directly whether laying those nodes out as a batch (data-oriented) and reclaiming them as a batch (an arena) turn out to be the same mechanism in practice or two mechanisms that only look alike on paper.
Until then, this note is the single place both records' rejected candidates point back to, rather than two separate dead ends.
