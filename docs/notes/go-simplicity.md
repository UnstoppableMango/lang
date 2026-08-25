# Go's simplicity philosophy

Poking at what "simple" meant to Go's designers, and whether any of it is a shape worth stealing (or explicitly rejecting) here.

## What Go actually means by simple

Not "few features" in the abstract.
Go has plenty of surface area (goroutines, channels, interfaces, reflection, cgo).
The simplicity is more like: few ways to do the same thing, and a strong bias toward one obvious idiom per problem.

Rob Pike's framing was often "less mechanism, more convention."
gofmt is the purest expression of this: remove a whole axis of bikeshedding by fiat, then never revisit it.
There is no config.
That's not a feature request that got rejected, it's a design stance: the tool decides, you don't get a vote.

Compare to Rust or C++, where "simple" would mean "powerful and orthogonal primitives that compose."
Go's bet was the opposite: primitives that are individually a little blunt, but few enough that everyone holds the whole language in their head at once.

## The "no" as a design tool

Go's simplicity is mostly visible in what shipped without:

- generics (for a decade, on purpose)
- exceptions (error values instead)
- inheritance (embedding instead)
- macros / metaprogramming
- operator overloading
- optional/named parameters, default arguments
- a expression-based ternary (if/else only)

Each omission has a critic's list of workarounds it forces (the `if err != nil` chorus, `interface{}` soup pre-generics).
Go's designers seem to have decided that a slightly worse local ergonomics story is worth it if it keeps the global "how many ways can this be written" number small.

Worth noting: they eventually reversed on generics.
So "no" wasn't sacred, it was a default that could be overturned once the cost of not having it clearly outweighed the cost of the complexity it would add.
That's a different posture than "never," it's "no until proven necessary," which is a much higher bar to clear than "no because we said so."

## Tooling as part of the simplicity story

A lot of what reads as "Go is simple" is actually "Go shipped a batteries-included toolchain that hides decisions."
gofmt, go vet, go test, go mod, a single canonical build command.
The language spec is short, but a huge amount of what would otherwise be bikeshedding (formatting, import ordering, project layout, dependency versioning) got absorbed into tooling instead of language features or community convention.

This is interesting because it's not really a *language* design lesson, it's a *project* design lesson: you can buy simplicity in the language by spending complexity in the tool.
If mangolang's tooling is thin or absent for a long time (see: no compiler yet), that tradeoff isn't available yet, so whatever "simple" looks like in this language has to hold up without a gofmt to lean on.

## Where Go's simplicity has critics

- Error handling: `if err != nil { return err }` repeated everywhere is simple to read line by line and exhausting to write. Simplicity of the individual construct, cost paid in repetition at scale.
- Nil: nil maps, nil slices, nil interfaces holding non-nil underlying values (the classic "typed nil" gotcha) — the simple mental model (nil is nil) breaks down exactly at the interface boundary, which is precisely where you need it not to.
- Zero values as the only "empty" story: elegant until a zero value is a legitimate value and you can't tell "unset" from "set to zero."
- Package-level simplicity vs. large-codebase simplicity: Go optimizes hard for "any Go programmer can read any Go file," which is a simplicity claim about the *reader*, not about the *language*. Those are different axes and can trade off against each other (see: not having generics meant `interface{}`-based container code that was individually simple but systemically repetitive).

So "simple" in Go is really "simple to read cold, for a stranger, without much context," bought partly at the expense of "simple to write" and "few footguns."

## Does that framing transfer here?

Open question worth sitting with rather than answering: is mangolang's audience the "stranger reading cold" reader, or someone who's opted in and is willing to hold more context?
Go's target was famously "the average Google engineer, three years out of a bootcamp, joining a team with a million-line codebase they didn't write."
That's a specific and somewhat unusual reader profile to optimize for.
If this language doesn't have that reader (huge company, huge codebase, high turnover, code review as the primary quality gate), the case for Go-style "convention over expression, one way to do it" is weaker, and a language that trusts the reader more might be the better bet.

Dead end (for now): tried to separate "simplicity of the spec" from "simplicity of the idiom" as two independent axes, but they're pretty entangled in Go's actual history, gofmt and vet are doing spec-adjacent work without being in the spec. Might be worth a cleaner framework later, not resolved here.

## Threads worth pulling later

- "no until proven necessary" as an explicit stance for this project's own decision records, distinct from "no forever."
- whether tooling-provided simplicity is even on the table without a compiler, and if not, what a language-only substitute would look like.
- the zero-value / nil-interface gotcha as a case study in "simple abstraction, complex boundary."
