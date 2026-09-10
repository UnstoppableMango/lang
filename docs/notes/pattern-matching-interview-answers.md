# Pattern matching: answers to the unblocked open questions

Answers given by the author on 2026-08-27 to open questions 1-6 in `docs/design/0005-pattern-matching.md`.
Questions 7 and 8 are blocked on the memory-strategy and paradigm decisions, and were not asked.

This is a note, so it carries zero standing.
The answers are not part of the design record until they are written into the design doc itself, which is stage 2 work.
This file exists so that a real wishlist entry has something to draw on.

## The answers

1. **Extensibility mechanism:** closed by default, with an explicit `open` marker required at a type's own definition before consumers must write a wildcard arm.
   Candidate C, not B's `#[non_exhaustive]`.
1. **Guard fallback:** an explicit fallback arm is required whenever a guard is present; arm-order-as-implicit-fallback is rejected as a footgun, not adopted from Rust.
1. **Match expression typing:** tentatively, arms infer a common ad hoc union type when they do not already agree, with a hard type-mismatch error kept explicitly as the fallback design if ad hoc union inference causes friction in practice.
   This is the one answer given as a hedge rather than a settled position.
1. **Structural matching:** allowed, and always requires a trailing wildcard arm, since exhaustiveness only ever applies to closed nominal sum types.
   Candidate D's structural-matching piece, adopted alongside candidate C's extensibility answer rather than instead of it.
1. **NaN in range patterns:** a float range pattern always requires a wildcard arm, and `NaN` is defined to fall through to it.
1. **Active patterns' future:** revisit later as a separate, possibly restricted feature.
   The rejection in the design doc's candidate E stays scoped to "not core matching," not "never."

## What this rules out

Everything below is derivation, not authored answer, and it is the part most likely to be wrong.

- **Candidate A, minimal matching only: dead.**
  Nothing in these answers walks back the wishlist's own commitment to guards and or-patterns; the full feature set is still in scope.

- **Candidate B's extensibility answer, `#[non_exhaustive]`: dead,** superseded by candidate C's closed-by-default, explicit-`open` mechanism.
  B's other features, nested destructuring, or-patterns with binding consistency, `@`-bindings, and matching on tuples, ranges, and literals, are all still adopted; only its specific consumer-side extensibility answer dies.

- **B's implicit guard-fallback behavior: dead.**
  Answer 2 chose an explicit fallback requirement instead of Rust's actual arm-order behavior, a real, deliberate divergence from the language B was otherwise modeled on.

- **Candidate C: chosen,** for the extensibility question specifically.

- **Candidate D: chosen,** for the structural-matching question specifically, composed alongside C rather than as a competing alternative.
  The rough shape's own framing, "a menu, not mutually exclusive tiers," turned out to be exactly right: the final answer takes C's extensibility mechanism, most of B's feature list, and D's structural piece, together.

- **Candidate E: not chosen, not permanently rejected either.**
  Answer 6 keeps the door open for a future, narrower version once core semantics are settled, matching the design doc's own "not core matching" framing rather than a stronger "never" stance.

## The shape this suggests

Sum types are exhaustiveness-checked by default; a type author marks one `open` at its definition to require consumers to write a wildcard arm from then on, rather than the language defaulting to permissive and asking authors to opt into strictness.
Nested destructuring, or-patterns with consistent bindings across every branch, `@`-bindings, and matching on tuples, ranges, and literals are all in scope, matching the wishlist's fuller ambition.
A guarded arm always requires an explicit fallback; the compiler never lets arm order alone stand in for coverage.
Structural destructuring works on any value with a matching shape, always behind a required wildcard, since there is no fixed constructor list to check exhaustiveness against.
A float range pattern is never exhaustive on its own for the same underlying reason, an uncountable domain with `NaN` as the sharp edge, and always requires a wildcard too.
A `match` expression's type is tentatively inferred as a common ad hoc union across its arms, with a hard type-mismatch error kept in reserve as the fallback if that inference proves troublesome once real programs are written against it.

That is closer to a real decision than a rough shape.
It is not yet a design, and the tensions below have to be resolved before it becomes one.

## New tensions the interview created

Each one is a candidate open question for the stage 2 pass.

1. **The tentative ad hoc union answer reaches outside pattern matching into the broader type system.**
   Structural union type inference is new type-system machinery this record did not otherwise need, and the hedge, "unless it causes friction," has no stated criterion for what friction would look like or how it would be measured.
   Stage 2 should either commit to the union answer with a real type-system home for it, or name the specific failure mode that would trigger falling back to a hard error, rather than leaving the hedge open-ended.

1. **Three independent triggers now require a wildcard arm: an explicitly `open`-marked nominal type, any structurally-typed match, and any float range pattern.**
   Each has its own justification, but a reader has to learn three separate rules that all resolve to the same syntactic requirement.
   Worth presenting as one unified "when do you need a wildcard" rule at stage 2 rather than three special cases.

1. **"An explicit fallback whenever a guard is present" needs a precise rule.**
   Does one guarded arm anywhere in a `match` require the compiler to prove the remaining, unguarded arms already cover the rest of the space without relying on it, or is a single trailing wildcard sufficient regardless of how many guarded arms precede it?
   The interview answer settles the principle, not the precise checking rule.

1. **Candidate C's own hard question, inherited unresolved from the note, still stands.**
   A type author now has to decide, at definition time, whether their type is closed or `open`.
   That is easy for an `Ordering`-shaped closed set and genuinely hard to answer honestly for a `Result<T, E>`-shaped type that might grow error variants over time; choosing C over B did not make this question easier, it just moved where the decision has to be made.

## Not asked

Questions 7 and 8 remain blocked: the memory-strategy record (`0001`) and the paradigm record (`0003`).
Nothing above should be read as deciding either.
