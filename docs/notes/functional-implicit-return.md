functional implicit return

the pitch: in a world where the language leans functional, every function body is an expression, and the value of that expression is the return value.
no `return` keyword needed for the common case.

```
fn add(a, b) = a + b

fn classify(n) =
  if n < 0 then "negative"
  else if n == 0 then "zero"
  else "positive"
```

this is the ML/Rust/Ruby lineage.
Rust keeps `return` around for early exit but the last expression in a block is implicit.
OCaml and Haskell don't have `return` as a keyword at all (Haskell's `return` is a monadic function, not control flow).
Ruby has implicit return of the last evaluated statement, but it's not really "functional", it just falls out of everything being an expression.

## what "implicit" buys you

- less ceremony for the common case (single-expression functions, especially in something math-y or DSL-shaped)
- pushes toward expression-oriented control flow (if/match as expressions, not statements) which composes better and avoids the "forgot to reassign in every branch" class of bug
- visually, function bodies read like definitions/equations rather than procedures: `fn square(x) = x * x` looks like the math it is

## what it costs

- early return.
  if the body is one big expression, how do you bail out of the middle of a loop or a long computation?
  options:
  - no early return at all, structure everything as expression composition (may need `Result`/`Option`-style short-circuiting operators instead, `?` in Rust, monadic bind elsewhere)
  - keep `return` as an escape hatch alongside implicit tail return (Rust's approach), but then you have two ways to return and have to teach both
  - loops become expressions too and you `break value` out of them (Rust does this too)
- the "did I mean to return this?" trap.
  add a debug `print(x)` as the last line of a function, and now that's the return value instead of unit/void.
  Rust dodges this because statements (things ending in `;`) evaluate to `()`, so a trailing `print(x);` returns unit, but drop the semicolon and you've silently changed the function's return value.
  that semicolon-as-semantic-marker is subtle and is a genuinely common beginner Rust mistake.
- readability for long functions.
  implicit return is great for one-liners, gets murkier the longer and more branch-y the function body gets.
  a 40-line function where the return value is "whatever the last expression several nested blocks deep happens to evaluate to" is harder to scan than one with an explicit `return foo` you can grep for.
- interacts with whether blocks are expressions at all.
  if `{ }` is a block *statement* (procedural, C-like) then implicit return doesn't even make sense, you'd need `{ }` to be a block *expression* whose value is its last sub-expression.
  that's a bigger structural commitment than it sounds like, it ripples into what `if`, `match`, loops, and even `;`-vs-newline mean.

## tangent: implicit return doesn't require "no return"

these are two separate axes, worth not conflating:

1. is a function body's trailing expression its return value by default? (implicit tail return)
2. is there a `return` keyword for early exit? (explicit early return)

matrix:
- yes / no: pure expression-oriented, no early exit except via expression composition (some point-free/combinator-heavy styles, or a hypothetical minimal ML dialect)
- yes / yes: Rust, Ruby (kind of, `return` exists but is rarely needed)
- no / yes: C, Go, most procedural languages, everything is a statement, `return` is mandatory
- no / no: doesn't really make sense, a function with no implicit tail value and no `return` can never produce a value except via non-local exit (exceptions) or side effect

so "functional implicit return" as a wishlist bullet is really quadrant 2 (yes/yes) or quadrant 1 (yes/no), and those are pretty different languages to write in.
quadrant 1 is more aesthetically "pure" but probably too austere for anything with loops and mutation-adjacent state, quadrant 2 is the pragmatic default basically every modern-ish language picks.

## open tension

does implicit return want expression-oriented `if`/`match`/blocks as a prerequisite, or could you bolt implicit-last-expression-return onto an otherwise statement-oriented language?
I don't think you can cleanly, the moment the return value is "whatever the block evaluates to," the block has to be able to evaluate to something, which means the block is an expression, which means if/match/loops are expressions too, which is a much bigger design commitment than "no `return` keyword".

so this note started as "should return be implicit" and it's actually a proxy for "is this language expression-oriented all the way down."
that's the real question hiding under the wishlist bullet.
