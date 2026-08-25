# Pit of success

Riffing on a driving philosophy: it should always be hard to do the wrong thing.
Not "the right thing is documented" or "the right thing is the example in the README," but structurally steered toward, the way a pit has sides.
You fall into the correct usage by doing the laziest, most obvious thing, and climbing out to do it wrong takes deliberate effort.

## Two flavors, not one philosophy

Soft pit of success: the right way is the path of least resistance, but the wrong way still compiles.
A nudge.
Example shape: the ergonomic API is three characters shorter than the dangerous one, so tired/rushed/copy-pasting code ends up correct by default, but an expert can still reach past it.

Hard pit of success: the wrong way is not expressible at all.
A wall, not a nudge.
Example shape: there is no null to dereference because there is no null.

These are not the same commitment and a single feature can't be graded "pit of success: yes/no" without picking one.
Soft degrades gracefully under deadline pressure but is escapable by accident, not just on purpose, which may defeat the point.
Hard is bulletproof but forecloses legitimate escape hatches: FFI, a perf hotpath, a deliberate prototype.
Working guess: the interesting design question per-feature is not "pit of success or not" but "which flavor, and is the door out of the hard version visible when it exists" (see the `unsafe` doodle below).

## The test to apply per feature

For any proposed feature: what's the wrong way to use this, and what happens when someone does it by accident, not maliciously, just tired or rushed or mid-copy-paste?
If the answer is "compiles fine, silently wrong at runtime," that's the red flag this whole note is about.
"Silently wrong" is worse than "loudly wrong," and "loudly wrong at compile time" is worse than "loudly wrong at runtime," and the philosophy is a bet that paying the cost as early and as loud as possible is always worth it, even when it costs expressiveness.

## Doodle: techniques that actually produce the pit

Six mechanisms, in decreasing order of how hard they are to weasel out of:

1. **Make the wrong thing unrepresentable.**
   No null type means no null-deref class of bug, full stop, not "usually caught," gone.
   This is the only technique that's bulletproof, and the only one that can't be added after the fact without a breaking rewrite of everything downstream, which is the actual argument for deciding this early rather than "cleaning it up in v2."

2. **Make the wrong thing require an incantation.**
   `unsafe { ... }` as a visible, lexical, grep-able marker, not a runtime flag or a config setting.
   The pit is still hard-walled, but there's a marked ladder out, and code review can `grep unsafe` the whole repo and count the ladders.
   Contrast with a "strict mode" opt-in (TypeScript's `strict: true`, JS's `"use strict"`): that inverts the pit, since the *default* stays loose and the thing you must opt into is the safe one.
   Everyone who doesn't know to opt in, which is every beginner pasting from Stack Overflow, lands in the unsafe default forever.
   An escape hatch has to escape *out of* safety, never *into* it.

3. **Make the correct API shorter to type than the incorrect one.**
   Asymmetric ergonomics as a design constraint on the standard library itself: if the safe function name is longer or the safe call site needs more boilerplate than the dangerous one, most people will reach for whichever their fingers find first, and that's the dangerous one.

4. **Collapse choice.**
   Already argued in [[dev-tooling-philosophy]] for formatting (`gofmt`, zero configuration): the fewer ways there are to do a thing, the fewer of those ways can be the wrong one.
   This generalizes past formatting to API surface: one canonical HTTP client, one canonical way to spawn concurrent work, one canonical serialization format for the stdlib's own types.
   Every second sanctioned way to do something is a second thing that can be the wrong one to reach for under time pressure.

5. **Make ignoring the result ugly, not just possible.**
   Rust's `#[must_use]` on `Result` is the compromise version of this: ignoring an error is still legal, but it produces a visible warning, so silence at the call site becomes a code-review flag instead of nothing.
   Go's `_ = err` shows the failure mode: the discard is exactly as easy to type as the check, so under pressure people type the shorter one, and the language has no opinion either way.
   Whatever this language does with errors, the discard path should never be the terse one.

6. **Sane zero values.**
   Go's actual pit-of-success win: every type has a zero value that's safe to use unintialized, so "forgot to initialize" degrades to "acts like the empty case" instead of "reads garbage memory."
   Cheap to copy, easy to forget it's a deliberate decision rather than an accident of how zeroed memory happens to look.

## Doodle: warnings don't work, so don't design around them

"Warn instead of error" sounds like the moderate compromise, and it isn't one.
`-Wall` gets turned off at the project level the first time it's inconvenient.
ESLint rules get an inline disable comment the first time they block a merge.
A warning is a suggestion with a snooze button, and at scale, snoozed is the default state.
If a check is really a pit-of-success feature, meaning the thing it catches is a real bug every time, it has to be a hard compile error, not a linter opinion that ships disabled-by-default in someone's CI three years later.
The corollary: don't add a check as "just a warning for now, we can make it stricter later," since later never comes, the warning ships forever exactly as easy to ignore as it was on day one.

## Doodle: "just document it" is a proven non-solution

Python's mutable default argument (`def f(x=[]):`) has been documented in every Python book and every "gotchas" blog post since the language existed, and it is still the single most common footgun new Python programmers hit.
Documentation is read once, if at all, usually after the bug already happened, never before the muscle memory forms.
A pit-of-success philosophy has to treat "we'll explain it in the docs" as equivalent to "we accept this will keep happening," because that's what the Python case actually demonstrates across decades of real usage.

## Where this fights itself

Pit of success and "hard things should be possible" are in real tension, not a fake one.
Elm's guarantee of no runtime exceptions is the hard-pit-of-success philosophy taken all the way: there is no escape hatch, not even a marked unsafe one, and the tradeoff is a language that structurally cannot do certain things at all (no partial functions, no throwing, no catching, ever).
That's either the purest version of this note's philosophy or a dead end depending on whether this language ever needs an FFI boundary, a perf-critical unsafe block, or an interop story with a language that doesn't share the guarantee.
Not resolved: is "no escape hatch, ever" the goal, or is "an escape hatch that's louder and uglier than the safe path" the goal?
These read similar in a pitch but produce very different type systems.

Also worth naming as tension: pit of success can tip into paternalism, the language refusing to let an expert do a thing they actually know is fine in their specific context.
Every hard-walled pit is also a wall the 0.1% expert case slams into.
The `unsafe` block doodle above is this note's best current answer (wall stays, ladder is visible and auditable), but it's worth stress-testing against a case where even the marked escape hatch isn't enough, before deciding this is the universal answer rather than just the best one found so far.

## Dead ends, recorded so I stop rediscovering them

- **"Document the footgun clearly."** Doesn't work at any observed scale. Python's mutable default argument is the standing counterexample: documented for decades, still the most common gotcha, because documentation is read after the bug, not before the keystroke.
- **"Ship it as a warning, tighten later."** Warnings get disabled the first time they're inconvenient (`-Wall` off, inline ESLint disables), and "tighten later" never arrives because the warning already shipped as ignorable and nobody wants to break existing code that relies on ignoring it.
- **"Opt-in strict mode for people who want the safe behavior."** Inverts the pit instead of removing it. The unsafe default is still the default, and it's now the default for everyone who didn't already know enough to opt out of it, which is precisely the population the philosophy exists to protect.

## Threads worth pulling later

- Whether "hard pit of success" (unrepresentable) is even affordable before the type system and memory model decision records land, since technique 1 above is really just a slogan for whatever those decisions turn out to be.
- Whether `unsafe`-as-visible-marker (technique 2) implies unsafe code has to be lexically scoped and non-composable with safe code in the same function, the way Rust does it, or whether a looser marker (a naming convention, a required doc comment) gets most of the grep-ability at less type-system cost.
- Whether "collapse choice" (technique 4) is really the same claim as [[dev-tooling-philosophy]]'s formatter argument generalized, or a distinct argument that happens to rhyme; they're filed as the same idea here but that might be sloppy.
- A concrete stdlib API (something like "read a file" or "spawn a task") worth designing twice, once the fast/obvious way and once the pit-of-success way, to see how far apart they actually land in character count and mental model.
