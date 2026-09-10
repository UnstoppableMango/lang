# A single canonical formatter, Go-style

\[[dev-tooling-philosophy]\] already named "Format: one canonical style, zero configuration" as one bullet in a checklist of eleven.
This note zooms into just that bullet.
Not "should there be a formatter," that question already got answered there (dead end: "ship the formatter as a separate optional install" and "configurable formatter with sane defaults" both die).
The question here is narrower and more interesting: what does "canonical" actually mean mechanically, and which of gofmt's specific design choices are the ones worth stealing versus the ones that are just what Go happened to land on.

## gofmt is less opinionated than its reputation

The popular take on gofmt is "it removes all formatting choices."
That's not quite true, and the gap is the interesting part.

gofmt does NOT:

- enforce a line length limit (no 80-col, no 100-col wrapping, ever)
- reflow a multi-line struct literal into one line, or vice versa
- decide how you group your `if` conditions across lines
- touch blank line placement beyond collapsing 2+ consecutive blanks into 1

gofmt DOES:

- normalize indentation to tabs, consistently
- normalize spacing around operators, commas, braces
- align consecutive struct tags / short var decls / map entries into columns (the tabwriter behavior)
- sort and group imports (goimports, technically a superset, but shipped as if it were gofmt to most users)
- guarantee `fmt(fmt(x)) == fmt(x)`, idempotence, non-negotiable

So the actual shape of gofmt's authority is: it owns whitespace *within* a line and *around* tokens, but it explicitly defers to the author on *where the line breaks go*.
If you write a struct literal on one line, it stays one line.
If you break it across five lines, it stays five lines (just re-indented and re-aligned).

That's a specific and narrow claim of territory.
Compare to prettier or rustfmt, which claim much more: both will reflow your code to fit a width, deciding FOR you whether a call expression fits on one line or needs to break, deciding where each argument goes.
Two genuinely different philosophies both wearing the "canonical formatter, no config" costume:

```
# gofmt-style: formatter owns tokens, author owns line breaks
formatter authority:  low reflow, high normalization
what survives review: "should this break across lines" (a real question, still asked)
what's eliminated:    tabs-vs-spaces, brace style, alignment-by-hand

# prettier-style: formatter owns almost everything
formatter authority:  high reflow, high normalization
what survives review: almost nothing formatting-related
what's eliminated:    line-break placement too, "does this call fit" is not a human decision anymore
```

Neither is obviously correct.
gofmt's approach leaves a residual bikeshed (people do still argue about when to break a struct literal across lines, gofmt just won't resolve it for you).
prettier's approach eliminates that residual but at the cost of sometimes producing reflows a human wouldn't have chosen, and diffs that change more than the semantic edit warranted (add one field to a call, prettier might reflow four other lines around it).

Worth sitting with: this is a real design fork for mangolang's formatter, not a solved question just because "have a canonical formatter" is settled.
"Zero config" doesn't specify which of these two territories the formatter claims.

## Doodle: what would each look like for a hypothetical mangolang literal

```
# gofmt-style (author's line breaks preserved)
point = Point{
    x: 1,
    y: 2,
}

# same author, one-lined, formatter leaves it alone
point = Point{ x: 1, y: 2 }
```

```
# prettier-style (formatter decides based on width, author's choice erased)
point = Point{ x: 1, y: 2 }
# ... six months later, someone adds a field ...
point = Point{
    x: 1,
    y: 2,
    z: 3,
}
# the formatter made this call, not the author, because it no longer fits some width threshold
```

The prettier-style version has a nice property: the two examples above are *the same code semantically* and a formatter-driven reflow means you never see a struct literal that "should" be multi-line but isn't, because someone wrote it before it grew a third field and nobody revisited it.
The gofmt-style version has a nice property too: the diff for adding `z: 3` is exactly one line, nothing else moves.
Diff minimality vs. consistency-under-growth. Real tradeoff.

## Idempotence and the parser/printer identity

The non-negotiable one, mentioned above without justification: `fmt(fmt(x)) == fmt(x)`.
Worth naming why this specific property matters more than it sounds like it should.

If a formatter isn't idempotent, running it twice in a row changes the file, which means "is this file formatted" is not a stable yes/no question, which means CI's `fmt --check` step can flap depending on formatter version drift or ordering of passes.
gofmt gets this almost for free because of *how* it's built: it's not a text-to-text transform, it's parse-to-AST, then pretty-print-AST-to-text, and the pretty printer is a pure function of the AST.
Feed the output back through the parser, get the same AST (comments aside, see below), pretty-print again, get the same text.
Idempotence isn't a property gofmt tries to have, it falls out of the architecture as long as parse and print are both deterministic.

This connects hard to \[[dev-tooling-philosophy]\]'s "tools that share one data structure" doodle: the formatter isn't really a separate tool from the compiler's front end at all, it's `print(parse(x))` using the exact same parser the compiler uses.
If mangolang's parser is a library (already argued for elsewhere), the formatter might not be eleven-tools'-worth of new work, it might be a few hundred lines of pretty-printer sitting on top of a parser that already has to exist.
That reframes "build a formatter" from "build a tool" to "write a printer for an AST you already have," which is a much smaller lift and makes "ship it on day one, not as a later toolchain item" more plausible.

## The hard part nobody markets: comments

gofmt's genuinely hard problem, the one that still generates bug reports after 15 years, is comment attachment.
Comments aren't part of the AST in the semantic sense (they don't affect program meaning), but they have to survive formatting attached to *something*, and figuring out what a comment is "about" from position alone is a heuristic, not a fact.

```
foo(
    a, # comment about a
    b,
)
```

Is that comment attached to `a`, or is it a trailing comment on the line?
What happens if `a` and its comment get reordered by some future refactor?
gofmt's answer is roughly "attach comments to the nearest AST node by position, and carry them along if that node moves," which works until it doesn't (comments have been known to jump to unexpected lines across large refactors, or duplicate, in edge cases involving `//go:generate`-style directives and build tags).

If mangolang wants comments to be structurally part of the grammar rather than trivia bolted on after parsing, that's a bigger, earlier decision (does the AST have comment nodes, or does the parser emit a side-channel of trivia the printer re-attaches by position), and it's exactly the kind of thing that's cheap to decide before there's a parser and expensive to retrofit after.
Not resolved here, just flagged as the actual hard 20% behind the easy-sounding 80% of "one canonical formatter."

## Doodle: formatting enforced at the parser, not just advised by a tool

A more aggressive idea than gofmt's: what if the compiler simply does not accept non-canonical whitespace at all, the way Python enforces indentation consistency (a mild, partial version of this) or the way some esolangs make whitespace semantically load-bearing?
Not "gofmt -l fails CI," but "the parser itself rejects a two-space indent."

Attractive because it closes the gap gofmt still has (nothing stops you from committing unformatted code locally, CI just yells at you after the fact).
Ugly because it destroys the "paste sloppy code and iterate" workflow, every incomplete draft mid-edit would fail to even parse, and pasting an example from docs or Stack Overflow written in a slightly different style becomes a hard error instead of a `fmt` away from fine.
Probably a dead end for exactly that reason: it conflates "canonical for committed/reviewed code" with "canonical for code that exists at all," and those don't need to be the same bar.
Gofmt's actual answer, advisory tool plus convention plus CI check, keeps write-time and commit-time as different bars on purpose.
Worth remembering next time this idea seems clever again.

## Threads worth pulling later

- Which of the two territories (gofmt's narrow "own the tokens, defer on line breaks" vs. prettier's wide "own the reflow too") mangolang's formatter should claim, this note deliberately didn't decide.
- Whether comments belong in the AST proper or as re-attached trivia, a decision that's cheap now and expensive after the parser exists.
- Whether alignment (gofmt's tabwriter-style column alignment for struct tags and consecutive assignments) is worth keeping. It's the one place gofmt does something visually loud and opinionated, and it's also the one gofmt behavior most frequently cited as annoying in large diffs (one new longer field name reflows alignment for every sibling line).
