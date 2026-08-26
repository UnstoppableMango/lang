# a parser combinator library is an initial target

The idea: pick "someone can write a parser combinator library in this language, and it feels good" as the first real target the language aims at.
Not a spec, not a milestone, a tire-swing.
The point of a target like this is that it forces a pile of open questions to become concrete at once, because a combinator library is a small program that leans on almost every feature a language has.

Why this one and not, say, "write a web server" or "write a CLI"?
A web server mostly tests the runtime and the standard library.
A combinator library mostly tests the *language*: functions as values, generics, composition syntax, error types, recursion, and how much ceremony sits between an idea and its expression.
And it is self-serving: the compiler currently parses with `nom`, so the day the language can host its own parser is the day self-hosting stops being hypothetical.

## what a combinator library actually demands

Working backwards from the smallest useful library, roughly in order of how load-bearing each is.

1. Functions as first-class values, including returning them.
   A combinator is a function that takes parsers and returns a parser.
   If that is awkward, everything downstream is awkward.
1. Closures that capture.
   `many(p)` has to hold onto `p`.
   How it holds it (copy, reference, ownership, refcount) is a question the memory decision answers, and this target makes it answer loudly.
1. Generics over the output type.
   `Parser<T>` where T varies per combinator, and `map(p, f)` turns a `Parser<A>` into a `Parser<B>`.
   Without this, the library is either stringly typed or unusable.
1. A failure-shaped return type and a way to chain it without pyramids.
   Every parser returns success-with-remainder or failure-with-reason.
   See \[[no-exceptions-explicit-errors]\] and \[[friendly-monads]\]: this is the exact shape those notes are circling.
1. Recursive definitions, including mutual recursion between grammar rules.
   `expr` calls `term` calls `factor` calls `expr`.
   Eager evaluation plus a value-level definition makes this an infinite loop at construction time in some designs, which is why so many libraries have a `lazy` or `defer` or `rec` escape hatch.

Beyond the top five, still relevant but less structural: operator syntax for sequencing and alternation, slices/views over the input that do not copy, and pattern matching on results (\[[robust-pattern-matching]\]).

## syntax doodles

Side by side, three flavors of the same tiny thing: parse one or more digits into an int.

Flavor A, method chaining, no special syntax:

```
let number = digit.many1().map(fn (ds) = int_of_chars(ds))
```

Flavor B, pipeline operator:

```
let number = digit |> many1 |> map int_of_chars
```

Flavor C, operator soup, in the tradition of parsec and fastparse:

```
let number = int_of_chars <$> many1 digit
```

Flavor A reads best cold and needs zero language support beyond methods on a type.
Flavor B needs a pipeline operator, which is a language-wide decision that should not be made because of one library.
Flavor C is dense and lovely once fluent, and is exactly the thing \[[pit-of-success]\] would push against: a newcomer cannot guess what `<$>` does.

Tentative lean: if the library can only be written comfortably in flavor C, the language failed the test.
Flavor C should be *possible* for people who want it, but the library that ships should read like A or B.

Now sequencing, which is where it gets interesting.

Flavor D, tuple-returning sequence combinators:

```
let assignment = seq3(ident, symbol("="), expr)
  |> map(fn ((name, _, value)) = Assign(name, value))
```

Flavor E, a do-notation / bind block:

```
let assignment = parse {
  let name = ident
  symbol("=")
  let value = expr
  Assign(name, value)
}
```

Flavor F, short-circuit operator inside a plain function:

```
fn assignment(input) = {
  let name = ident(input)?
  symbol("=")(name.rest)?
  let value = expr(...)?
  Assign(name, value)
}
```

D is honest and gets ugly fast at arity 5 or 6, and every library that starts there ends up with `seq7`.
E is the nicest to read and costs a real language feature, one that generalizes past parsers if it is built on an interface rather than special-cased (\[[friendly-monads]\] calls this exact fork).
F is tempting because it reuses the `?` operator the wishlist already wants, but notice the threading of `rest` through by hand.
That threading is the State monad wearing a disguise, and doing it manually is precisely the boilerplate combinators exist to remove.
So F is a dead end for the *library*, though it might be fine for a hand-written recursive descent parser, which is a different program.

## the recursion problem, chased a bit

If parsers are values and evaluation is eager:

```
let expr = alt(binary_op, term)
let term = alt(paren_expr, number)
let paren_expr = seq3(symbol("("), expr, symbol(")"))
```

`expr` is defined before `term` exists, and `paren_expr` refers back to `expr`.
In a world with laziness this is free.
In a world with eager evaluation, the usual answers are:

- Make parsers functions rather than data, so the reference is only resolved at call time.
  Cost: allocation or indirection per call, and the combinators get slightly noisier.
- A `rec` or `lazy` marker on the definition.
  Cost: a language feature, and users forget it and get a stack overflow at startup, which is a terrible first experience.
- Let the compiler figure out the cycle among top-level bindings, the way ML-family `let rec ... and ...` does.
  Cost: the compiler has to reason about value-level recursion, and the error when it *cannot* is subtle.

Third option is the most pleasant and the most work.
Worth noticing that the wishlist already wants dependency cycles *rejected* at the module level (\[[directory-scoped-modules]\]) while this wants them *embraced* within a module.
That is not a contradiction, but it is a place where "no cycles" needs a qualifier.

## the parts that only show up when you actually write one

Things that a toy version hides and a real library exposes.

- Error messages.
  "expected one of: `+`, `-`, `(`, digit, at line 4 column 7" requires alternation to *accumulate* expectations rather than discard them, which means the failure type is richer than a string, and combinators must merge failures.
  This is the same quality bar the wishlist sets for compiler diagnostics, applied to a user library.
  Nice property: if the language makes good errors easy for a library author, it probably made them easy for itself.
- Backtracking control.
  Parsec's `try`, and the distinction between "failed without consuming" and "failed after consuming" is the single subtlest thing in the whole design space, and every user trips on it.
  Does the language help here at all, or is this purely a library API problem?
  Suspect the latter, which is useful information: not every hard part of the target is a language problem.
- Zero-copy input.
  A parser wants to hand back a slice of the original input, not a fresh string per token.
  This is a direct question for the memory decision, and \[[arena-memory-model]\] would answer it very differently than refcounting would.
  A combinator library is a good stress test precisely because it allocates many small short-lived things and a few long-lived ones.
- Performance expectations.
  Nobody minds a slow combinator library for a config file, everybody minds one in a compiler.
  If the answer is "combinators are 10x slower than hand-written descent, use them for prototyping," that is a legitimate answer, but it should be a chosen answer.

## the self-hosting angle

If the language can host a combinator library, and a combinator library can parse the language, the parser stops being Rust-shaped.
That is a long way off and worth writing down anyway, because it changes what "good enough" means.
A library that is pleasant for parsing JSON is a much lower bar than one that can carry a real grammar with decent errors and acceptable speed.
Related: \[[ast-in-public-api]\] wants the AST to be a public library, which pairs oddly well with this.
If the AST is public and the parser is a library written in the language, the boundary between "the compiler" and "a program someone wrote" gets thin, which might be the actually interesting outcome here.

## dead ends noted

- Making parser combinators a *built-in language feature* (a grammar DSL baked into the compiler).
  Died on contact with "few ways to do one thing": a built-in grammar syntax is a second language inside the language, and it removes exactly the pressure this target is supposed to apply.
  The whole value of the target is that it is written in ordinary user-level code.
- Judging the language by whether the operator-heavy flavor C is expressible.
  Died because expressible is cheap; the target is whether the *readable* version is expressible without contortions.

## open tension

Using a library as a design target risks over-fitting the language to one shape of program.
Combinators are heavily functional, heavily generic, heavily allocation-happy.
Optimize for them too hard and the language quietly commits to a paradigm that no decision record has decided yet, which is exactly the thing the workflow says not to presume.

The counter-argument: the target is not "make this library fast and beautiful at any cost," it is "use it as a probe to find where the language hurts."
A probe is allowed to be unrepresentative as long as nobody mistakes the probe for the goal.

Second-order worry that might matter more: a target this specific is a great way to notice missing features and a terrible way to notice *excess* ones.
Nothing about writing a combinator library will tell you the language has too much in it.
Might want a second target with the opposite bias, something deliberately boring, before treating this one as the compass.
