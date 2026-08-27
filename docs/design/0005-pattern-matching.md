---
title: Pattern matching
stage: 1
status: active
created: 2026-08-27
updated: 2026-08-27
---

# Pattern matching

## Motivation

The wishlist already commits to "pattern matching with compiler-checked exhaustiveness, nested destructuring, guards, and or-patterns," and the paradigm record's own interview leans on pattern matching as the primary way to inspect algebraic data.
`docs/notes/robust-pattern-matching.md` found that exhaustiveness checking alone answers one narrow question, did you handle every shape, and says nothing about destructuring depth, guards, or-patterns, matching on non-sum-type values, or what happens when a library adds a variant to a type every downstream `match` already covered exhaustively.
This record exists to decide how much of that "robust" feature set the language actually commits to, not just whether exhaustiveness checking exists at all.

### Scenario: matching the compiler's own tokens

A tokenizer produces a stream of tokens, and a parser needs to react differently to several related token shapes at once, `Plus`, `Minus`, `Star`, `Slash` should all route to the same operator-parsing function, while `Number(n)` routes elsewhere, and anything else is an error.
Writing this without or-patterns means repeating the same arm body four times or falling back to an `if`-chain, exactly the ceremony pattern matching exists to remove.

### Scenario: a library adds a variant to a public enum

A library exposes a `Token` enum with five variants.
Every downstream consumer writes an exhaustive `match` over it, relying on the compiler to flag any missing case.
The library author adds a sixth variant.
Every one of those downstream matches is now a compile error, everywhere, instantly, which is exhaustiveness checking working exactly as designed and also a breaking change the library author may not have intended to make.
The language has to pick a side: exhaustiveness's strongest guarantee, or library evolution without breaking every consumer, and say which one is the default.

### Scenario: destructuring the compiler's own AST

A lowering pass matches on an `Expr` node several levels deep, wanting to bind a nested field, guard on a condition over that field, and fall through to a more general arm if the guard fails.
This is the paradigm record's own central scenario, the compiler's AST, read through pattern matching specifically: `Add(Lit(0), rhs) => rhs` folding an identity addition is exactly the shape a lowering pass needs to write directly, not simulate through nested `if`/`else`.

## Rough shape

Five candidate directions, sketched from narrowest to widest feature scope.
These are not mutually exclusive tiers so much as a menu; a real design likely takes pieces from several rather than adopting exactly one wholesale.

### A. Minimal: exhaustiveness and bare destructuring only

`match` supports one level of constructor destructuring and compiler-checked exhaustiveness over closed sum types, full stop.
No guards, no or-patterns, no `@`-bindings, no matching on tuples, ranges, or literals.

```
match opt {
    Some(n) => n,
    None => 0,
}
```

Buys: the smallest possible surface to specify and implement, and it is literally the wishlist entry's own bare example.
Costs: `docs/notes/robust-pattern-matching.md`'s own finding, a language can ship this and still feel thin the moment a real program wants to test a condition alongside a shape, or handle four variants the same way, forcing nested `match` or a fallback to `if`-chains for exactly the cases pattern matching is supposed to remove.

### B. Rust-shaped: full destructuring, guards, or-patterns, consumer-side `#[non_exhaustive]`

Nested destructuring to arbitrary depth, guards, or-patterns with a binding-consistency requirement across every branch of the or, `@`-bindings, and matching on tuples, ranges, and literals in addition to sum types.
A type author opts an enum into "consumers must always write a wildcard arm" to allow adding variants later without breaking downstream matches.

```
match shape {
    Circle { radius: r } if r > 10.0 => "big circle",
    Circle { .. } => "circle",
    Rect { width: w, height: h } if w == h => "square",
    Rect { .. } => "rect",
}
```

Buys: the most complete, most battle-tested feature set of any candidate, directly matching everything the wishlist entry names.
Costs: `#[non_exhaustive]` trades exhaustiveness's main benefit, the compiler tells you every call site needing an update, for library stability, and the guard-exhaustiveness gap, a guard can fail so a guarded arm never fully covers its pattern, means arm order becomes implicit fallback logic, which some code leans on as a decision table and other code trips over as a footgun.

### C. Closed by default, explicit opt-in to open

The inverse of B's extensibility answer: every sum type is exhaustiveness-checked as strictly as possible by default, and only a type explicitly marked extensible at its own definition requires consumers to write a wildcard arm.

```
open enum HttpMethod { Get, Post, Put, ... }   // consumers must handle a wildcard
enum Ordering { Less, Equal, Greater }          // closed; adding a variant is a breaking change, by design
```

Buys: reads more consistently with "compiler-checked exhaustiveness" taken literally as the default, strong guarantee, and puts the extensibility decision at the definition site where the type's author actually knows whether the set is meant to grow, rather than leaving every consumer to guess.
Costs: forces every type author to decide up front, honestly, whether they will ever add a variant, a question that is easy for an `Ordering`-shaped closed set and genuinely hard for a `Result<T, E>`-shaped one that might grow error variants over time; this is really a module-versioning question wearing a pattern-matching costume.

### D. Structural matching alongside nominal exhaustiveness

Destructuring works uniformly over nominal sum types, tuples, structs, and open, structurally-typed values, but exhaustiveness checking only ever applies to closed nominal sums.
Matching on a structural or open type always requires a trailing wildcard.

```
match value {
    { name, age } if age >= 18 => s"${name} is an adult",
    { name, .. } => s"${name} is a minor",
}
```

Buys: lets structural destructuring, pulling named fields out of anything with that shape, coexist with the strong exhaustiveness guarantee where it can actually apply, rather than forcing a single all-or-nothing answer across two structurally different kinds of openness.
Costs: two different reasons a wildcard might be required, an explicitly `open`-marked nominal type from candidate C and any structurally-typed value here, is two concepts for the author to track instead of one, and the two need to read as one coherent story rather than two unrelated escape hatches at stage 2.

### E. Active or view patterns as an escape hatch: rejected for this pass

A pattern calls a function first, then matches on the result, decoupling "what can appear in a match arm" from "what the type's real constructors are."

```
active pattern Even(n) => n % 2 == 0
active pattern Odd(n)  => n % 2 != 0

match n {
    Even => "even",
    Odd => "odd",
}
```

Buys: real expressive power, F#'s actual production feature, letting a match arm test an arbitrary computed condition as if it were a constructor.
**Rejected for this record.**
A "total" active pattern's claim to be exhaustive is an assertion the compiler cannot check without evaluating arbitrary code, F# resolves this by taking the author's word for it, which is exactly the kind of unchecked claim this language's other wishlist commitments, compiler-checked exhaustiveness chief among them, exist to avoid.
This does not rule active patterns out forever; it says the core matching feature this record specifies should not depend on an unchecked exhaustiveness assertion to be coherent.
A future note or feature can revisit active patterns once the core semantics here are settled.

## Prior art

- **Luc Maranget, "Warnings for pattern matching"** (Journal of Functional Programming, 2007; [PDF](http://moscova.inria.fr/~maranget/papers/warn/warn.pdf)).
  The algorithm actually used to detect non-exhaustive matches and useless clauses in OCaml, formulated over matrices of patterns, with explicit limitations noted for guards and GADTs.
  Directly relevant to every candidate above: exhaustiveness checking is not a vague aspiration, it is a specific, decades-old, well-understood algorithm with known blind spots, particularly around guards, that this record inherits rather than invents.

- **Rust RFC 2008, `#[non_exhaustive]`** ([RFC text](https://github.com/rust-lang/rfcs/blob/master/text/2008-non-exhaustive.md); [RFC book](https://rust-lang.github.io/rfcs/2008-non-exhaustive.html)).
  The consumer-side opt-out this record's candidate B adopts: an enum or struct author marks it as liable to grow, and every downstream match is required to carry a wildcard arm from that point on.
  Direct precedent for the exhaustiveness-vs-library-evolution tension, and for candidate C's inverse framing, closed by default with an explicit `open` marker, being read against the same problem from the other direction.

- **Rust RFC 2175, or-patterns** ([RFC book](https://rust-lang.github.io/rfcs/2175-if-while-or-patterns.html)).
  Established that or-patterns need to work consistently across `match`, `if let`, and `while let`, and that every branch of an or-pattern must bind the same variables with the same types, the binding-consistency requirement candidate B names explicitly.

- **F# active patterns** ([Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/active-patterns)).
  Distinguishes total (complete) active patterns, which the compiler trusts to be exhaustive, from partial active patterns, which return an option type and force a wildcard arm.
  The clearest shipped precedent for candidate E, and for exactly the "on trust" cost that got it rejected here: F#'s own design has lived with this tradeoff for years rather than resolving it away.

Where the search was thin: nothing found treats "should exhaustiveness checking apply to structurally-typed, open values at all" (candidate D) as a solved problem; most pattern-matching literature, including Maranget's own formulation, assumes a closed, nominal set of constructors as the starting premise.
That gap is recorded as an open question rather than folded into a candidate, since it is closer to an open research question than a settled design choice any prior art directly answers.

## Open questions

1. Does the exhaustiveness-vs-library-evolution tension get Rust's consumer-side `#[non_exhaustive]` answer (candidate B), or the closed-by-default, explicit-`open`-at-definition inverse (candidate C)?
1. Does an arm covered only by a guard still require a trailing wildcard or fallback arm, given a guard can fail, and is relying on arm order as the implicit fallback, Rust's actual behavior, accepted as intentional or rejected as a footgun?
1. What is a `match` expression's type when its arms produce different types: a hard type error, or unification to some common type?
   The wishlist's own "expressions over statements wherever possible" entry presumes `match` is an expression, which makes this question unavoidable rather than optional.
1. Does pattern matching apply to open, structurally-typed values at all (candidate D), and if so, is exhaustiveness simply inapplicable there, always requiring a wildcard, or does the language attempt some other model the prior-art search did not find?
1. What happens when a range pattern's scrutinee is `NaN`: does it fall through to a wildcard arm, is it a compile error when no wildcard is present, or are range patterns over floating-point values disallowed outright?
1. Are active or view patterns in scope for a later feature, now that they are rejected for this record's core semantics, or should this record say plainly that the language never adopts an unchecked exhaustiveness assertion in any form?
1. **Blocked on an undecided foundational decision:** the memory-strategy record (`0001`).
   Whether destructuring a value in a match arm moves, borrows, or copies its fields depends on that record's regime mechanics in more detail than its interviewed-but-not-yet-stage-2 answers currently specify.
1. **Blocked on an undecided foundational decision:** the paradigm record (`0003`).
   Whether a match arm can destructure through a value typed only by the structural interface it satisfies, or only through its own concrete underlying type, depends on that record's still-open interface-dispatch mechanics.

## Advancement record

- 2026-08-27, gate 0 → 1: sketched from the wishlist entry "Pattern matching with compiler-checked exhaustiveness, nested destructuring, guards, and or-patterns" and `docs/notes/robust-pattern-matching.md`; five candidate directions surveyed narrowest to widest, with active/view patterns rejected for this pass on an unchecked-exhaustiveness-assertion basis; prior art cited across Maranget's exhaustiveness algorithm, Rust's `#[non_exhaustive]` and or-pattern RFCs, and F#'s active patterns; open questions recorded including two blocks on undecided foundational decisions.

## Changelog
