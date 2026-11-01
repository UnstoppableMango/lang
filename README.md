# lang

A general-purpose programming language, designed from scratch and deliberately slowly.

There is no spec yet.
The compiler, `unmangc`, is written in Rust and compiles `hack/hello.lang` end to end; everything past that is still being designed.

## Principles

Every design choice is held to these.
My [influences](docs/influences.md) inform these choices too, but only the principles bind.

### Few ways to do one thing

Favor convention over configurable primitives, so the obvious way to write something is the only way.
A feature stays out until it proves necessary: "no" is the default, not a promise, and the bar for overturning it is that the cost of not having the feature clearly outweighs the complexity it adds.
Prior art: Go, and `gofmt`'s refusal to be configured.

### Hard to do the wrong thing by accident

The laziest, most obvious way to use a feature should be the correct one, and doing it wrong should take deliberate effort.
Where possible, make the wrong thing unrepresentable rather than merely discouraged.
Where an escape hatch must exist, mark it visibly and greppably (e.g. `unsafe`), so it escapes out of safety and never makes the unsafe path the quiet default.
A real bug is a compile error, not a warning, and documenting a footgun does not count as removing it.
Prior art: Rust's `unsafe` blocks and `#[must_use]`.

## Milestones

The language grows toward small, concrete milestones, each a program that must compile and print an expected output in CI.
They are grouped into tiers, lowest first:

- **A. Hello world**: 1 Hello, 2 Arithmetic.
- **B. Turing complete**: 3 Truth, 4 Choice, 5 Functions, 6 Recursion, 7 FizzBuzz.
- **C. Data**: 8 Records, 9 Sums, 10 Evaluator, 11 Strings, 12 Lists, 13 Closures.
- **D. Real programs**: 14 Generics, 15 Errors as values, 16 Two files, 17 Reading stdin, 18 Tokenizer, 19 Parser combinator library.
- **E. Bootstrapping**: the compiler, written in the language, compiles itself.

Each tier is a [GitHub milestone](https://github.com/UnstoppableMango/lang/milestones), and an open decision or idea is attached to the tier that first needs it.
A milestone's acceptance program lands in `milestones/` with the pull request that makes it pass.

## Open decisions

The foundational decisions are open, and nothing may presume one until its issue is closed by an ADR in `docs/adr/`:

- [Memory strategy](https://github.com/UnstoppableMango/lang/issues/97)
- [Compilation model](https://github.com/UnstoppableMango/lang/issues/98)
- [Paradigm](https://github.com/UnstoppableMango/lang/issues/99)
- [Concurrency model](https://github.com/UnstoppableMango/lang/issues/100)
- [Pattern matching semantics](https://github.com/UnstoppableMango/lang/issues/101)

Feature ideas are [`needs-triage` issues](https://github.com/UnstoppableMango/lang/issues?q=is%3Aissue+is%3Aopen+label%3Aneeds-triage), and `CONTEXT.md` is the glossary.

## Building

Nix drives everything:

- `make build` builds the compiler.
- `make test` runs the golden and unit tests.
- `make run` inside `hack/` compiles and runs `hack/hello.lang`.

See `AGENTS.md` for the full set of commands.
