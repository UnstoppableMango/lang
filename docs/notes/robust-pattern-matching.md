# Robust pattern matching

Wishlist already has "pattern matching with compiler-checked exhaustiveness."
This note is about what "robust" adds on top of the bare-minimum `match x { Some(n) => n, None => 0 }` pitch: what makes matching feel complete versus merely present.

## What "just exhaustiveness" leaves out

Exhaustiveness checking answers one question: did you handle every shape.
It says nothing about:

- destructuring depth (can you match two levels into a nested structure in one arm, or do you nest matches)
- guards (can an arm say "this shape, but only if...")
- binding the whole thing and a piece of it at once
- or-patterns (multiple shapes, one arm)
- matching on things that aren't sum types at all: tuples, structs, ranges, literals, strings
- what happens when a library adds a new enum variant and every downstream `match` was exhaustive

A language could ship exhaustiveness checking and still feel thin if match arms can only be one constructor deep.

## Doodle: nested destructuring

```
match shape {
    Circle { radius: r } if r > 10.0 => "big circle",
    Circle { .. } => "circle",
    Rect { width: w, height: h } if w == h => "square",
    Rect { .. } => "rect",
}
```

Guards (`if r > 10.0`) are the obvious want once you have matching at all, but they reopen the exhaustiveness question: a guard can fail, so `Circle { radius } if r > 10.0` doesn't fully cover `Circle`, the checker has to know that and still demand a fallback `Circle { .. }` arm (or fold the guard's "else" into the next pattern, which is what happens above by accident, relying on arm order).
Is arm-order-as-fallback intentional design or a footgun?
Rust accepts it; some guard-heavy code reads like a decision table but silently depends on top-to-bottom order, which is exactly the kind of implicit control flow [[no-exceptions-explicit-errors]] worries about in a different context.

## Doodle: or-patterns

```
match token {
    Plus | Minus | Star | Slash => parseOperator(token),
    Number(n) => parseLiteral(n),
    _ => error("unexpected token"),
}
```

Cheap to want, mildly annoying to implement well: every branch of an or-pattern has to bind the same variables with the same types, or the compiler has to reject asymmetric bindings (`Ok(x) | Err(x) =>` binding `x` to different types depending on which arm matched is a trap, not a feature).

## Doodle: binding the whole and the part (`@`-patterns)

```
match msg {
    full @ Resize { width, .. } if width > 4096 => log("huge resize", full),
    Resize { width, height } => resize(width, height),
}
```

Small feature, shows up constantly once people write real matches; easy to forget when sketching the grammar and annoying to retrofit later since it changes pattern syntax everywhere at once, not just one spot.

## The exhaustiveness-vs-library-evolution tension

Exhaustiveness checking is great until a library owner adds a variant to a public enum.
Every downstream `match` that used to be exhaustive is now a compile error, everywhere, instantly.
Two answers exist in the wild:

1. **Rust's `#[non_exhaustive]`**: the enum author opts an enum into "consumers must always write a wildcard arm," trading exhaustiveness's main benefit (compiler tells you every call site that needs updating) for library stability.
2. **Closed by default, explicit opt-in to open** (the inverse): most enums are exhaustive-checked hard, and only enums explicitly marked extensible get the wildcard requirement.

(2) seems more consistent with a "compiler-checked exhaustiveness" wishlist entry read literally, since it keeps the strong guarantee as the default and makes the escape hatch visible at the definition site rather than the call site.
But it means every enum author has to decide up front whether they'll ever add a variant, which is a hard question to answer honestly for a `Result<T, E>`-shaped type versus a closed `Ordering { Less, Equal, Greater }`-shaped type.
Undecided; flagging that this is really a modules/versioning question wearing a pattern-matching costume.

## Matching against non-sum-type values

Once matching exists, people want to reach for it on everything:

```
match n {
    0 => "zero",
    1..=9 => "small",
    10..=99 => "medium",
    _ => "large",
}

match point {
    (0, 0) => "origin",
    (x, 0) => "on x-axis",
    (0, y) => "on y-axis",
    (x, y) => "elsewhere",
}

match name {
    "get" | "post" | "put" => httpVerb(name),
    other => customVerb(other),
}
```

Ranges and tuples are easy structurally (a tuple pattern is just nested destructuring on an anonymous product type).
String literal matching is deceptively expensive at scale: naive codegen is a chain of string-equality checks, a good compiler wants to build a trie or jump table, and that's a code-quality bet, not a semantics one, so it's more a "does the compiler team care yet" question than a language-design one.
Range patterns on floats are the interesting trap: `0.0..=1.0` looks fine until NaN, and matching semantics have to say explicitly what happens when the scrutinee is NaN (falls through to `_`? compile error if there's no `_`? silently unreachable?).
Parking as an open wrinkle, not solved here.

## Doodle: matching through structural/duck-typed values

[[duck-typing]] leans toward structural typing checked statically.
Does pattern matching work on structural types the way it works on nominal sum types?

```
match value {
    { name, age } if age >= 18 => s"${name} is an adult",
    { name, .. } => s"${name} is a minor",
}
```

This reads fine for a single shape, but exhaustiveness on structural types is a much harder problem than on closed nominal enums: there's no fixed list of "constructors" to enumerate, just an open set of possible shapes.
Practically this probably means: structural destructuring exists and is useful (pull `name` and `age` out of anything with those fields), but *exhaustiveness checking* only applies to closed nominal sums, and matching on a structural/open type always requires a trailing `_` or equivalent.
That's consistent with the `#[non_exhaustive]` discussion above: openness and exhaustiveness are in tension no matter which axis (nominal vs structural, or closed vs extensible-nominal) produces the openness.

## Doodle: view patterns / active patterns as an escape hatch

F# active patterns and Haskell's (proposed, contentious) view patterns let a pattern call a function first, then match on the result:

```
active pattern Even(n) => n % 2 == 0
active pattern Odd(n)  => n % 2 != 0

match n {
    Even => "even",
    Odd => "odd",
}
```

Interesting because it decouples "what can appear in a match arm" from "what the type's constructors literally are," which is powerful and also exactly the kind of feature that makes exhaustiveness checking undecidable in general (the compiler can't know `Even`/`Odd` are complementary and exhaustive without evaluating arbitrary code).
F# handles this by allowing *partial* active patterns to force a wildcard requirement, and *total* active patterns (ones the author asserts are exhaustive, like a real partition) to participate in exhaustiveness checking on trust.
"On trust" is doing a lot of work there: it's an unchecked assertion, the compiler takes the author's word for it.
Feels like a wilder, later-stage feature, not core matching.
Parking here rather than promoting it.

## Match as expression vs match as statement

Given the wishlist's "expressions over statements wherever possible" entry, `match` presumably has to be an expression (`let x = match y { ... }`), which drags in the question of what type a match expression has when arms return different types, and whether that's a hard error or unifies to a common supertype/union.
That's a type-system question this note shouldn't try to resolve, just flagging the coupling: robust pattern matching and expression-oriented control flow aren't independent decisions, they're the same decision seen from two notes.

## Where this connects

- Matching on `Option<T>` and `Result<T, E>` is the load-bearing use case referenced in [[no-null-type-or-representation]] and [[no-exceptions-explicit-errors]]; whatever "robust" means here has to make those two notes' sketches read naturally, not just the toy `Some`/`None` examples.
- Structural destructuring intersects [[duck-typing]]'s structural-vs-nominal axis directly, see above.
- Exhaustiveness-vs-library-evolution is really a module/versioning question; might deserve its own note someday rather than living as a subsection here.

## Dead end: pattern matching as the *only* control flow

Briefly tempting (Erlang/Elixir go partway here: function clauses are pattern matches).
Died on: `if`/`while`/loops still want to exist as distinct, readable constructs for the common case, and forcing every branch into a `match` arm just to test one boolean is ceremony without benefit.
`if let` / `while let`-style sugar (match a single pattern, fall through otherwise) seems like the right middle ground, worth a doodle of its own later, not resolving the syntax here.
