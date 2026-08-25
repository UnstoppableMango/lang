# Fork and trim dependencies

Starting premise: minimizing dependencies is not a value to nudge people toward, it is a load-bearing design philosophy.
The mechanism under test: when you take a dependency, you fork its source, trim it down to what you actually use, and only the trimmed fork ever reaches the compiler.
Zero standing, presumes outcomes of open decision records freely.

\[[dependency-cultures]\] already recorded a dead end that looks adjacent to this: "forbid transitive dependencies" dies because everyone vendors anyway and the graph becomes invisible instead of shallow.
This note exists to find out whether mandatory, tooled vendoring is the same dead end wearing a disguise, or a genuinely different animal because the vendoring is forced, automatic, and visible in-repo rather than an escape hatch someone reaches for quietly.

## The basic shape

Not "add a dependency", but "import and trim" as a single verb.

```
$ lang vendor add http@a1b2c3
  fetched http (380 KB, 41 files)
  computing reachability from your call sites... 6 symbols used
  trimmed to 6 symbols + transitive closure: 4.1 KB, 3 files
  wrote vendor/http/ (source, not a binary blob)
  diff is 3 files, 190 lines. review before committing.
```

The compiler never talks to a registry, never resolves a version range, never fetches anything.
By the time `lang build` runs, `vendor/http/` is just... your source.
No special-cased "this is a dependency" concept survives past the vendor step.
That is the part I like best here: an entire subsystem (package manager, resolver, lockfile format) stops being a language design question and becomes a one-time source-transformation tool that could in principle live outside the language entirely.

## Why "trim" and not just "vendor"

Plain vendoring (copy the file in, C-style, per \[[dependency-cultures]\]) freezes a copy but keeps the whole thing, including the 95% of the package you never call.
Trimming is vendoring plus reachability analysis: dead code elimination performed at fork time, on source, before a single line reaches the compiler, rather than as a linker pass afterward.

Consequence I like: the diff you actually review is small.
"Add a dependency" today means trusting 380 KB you will never read.
Here it means reviewing 190 lines you can plausibly read in five minutes, because the tool already did the discarding for you.

Consequence I am less sure about: reviewing a trim is not the same as reviewing the package.
You are reviewing "the parts reachable from my current call sites", which is a moving target every time you add a new call.
Add one more function call next month and the trim silently grows, possibly pulling in code nobody has looked at.
So the safety property is weaker than it first sounds: it caps initial exposure, it does not cap it forever.

## Doodle: does this actually dodge the "vendoring makes the graph invisible" dead end

The old dead end: forbid transitive deps, everyone vendors by hand, and now the dependency graph is smeared across the codebase instead of declared anywhere, so nobody can answer "what do we depend on" without grepping.

Fork-and-trim might dodge this if the tool is the one doing the vendoring, and it leaves a receipt:

```
vendor/http/.origin
  source: http@a1b2c3
  trimmed: 2026-08-25
  reachable-from: src/api/client.gos:12, src/api/client.gos:45
```

That receipt is the dependency graph, it just lives next to the code instead of in a manifest.
"What do we depend on" becomes `grep -r origin: vendor/` instead of reading package.json, which is a wash on convenience but a win on honesty, since the receipt can't drift from the code the way a manifest can.

Or it might not dodge the dead end at all, if in practice `vendor/` becomes a pile of 200 divergent forked snippets that nobody re-syncs, which is just left-pad's problem with extra steps and a bigger diff.
Genuinely unsure which failure mode dominates.
The old dead end assumed *manual* vendoring; I don't have a strong intuition for whether *tooled, mandatory* vendoring inherits the same failure or not.

## Doodle: what "update" even means now

Once you've trimmed, you no longer have "the package at a newer version", you have your derivative of an old version.
Updating is not `bump 1.2.3 -> 1.2.4`, it is:

```
$ lang vendor update http
  new upstream: http@d4e5f6
  re-running reachability from your current call sites...
  new trim: 4.4 KB, 3 files
  diff against your vendored copy: 2 files changed, 12 lines
  (your vendored copy had 1 local edit that upstream does not have: vendor/http/retry.gos:8)
  merge? [y/N]
```

This reframes "update" as "regenerate and diff", closer to a lockfile regen than a git merge, which is nice because it is mechanical and re-runnable.
But it exposes the sharp edge immediately: if you ever hand-edited the trimmed fork (fixed a bug, patched a security hole yourself because upstream was slow), updating means merging your patch against a re-trim, and that is a real merge conflict, not a version bump.
\[[dependency-cultures]\] flagged "security patches" as the unresolved cost of content-addressed dependencies, where a name is just an alias for a hash and there's no mechanism for "everyone who depends on this must move".
Fork-and-trim has the opposite problem: there is a mechanism (re-vendor), but it competes with your own local edits every single time, which is a tax the content-addressing approach doesn't pay, because it never lets you edit the dependency in the first place.

## Doodle: trimming pushes the "granularity of the unit" lever without changing what a package is

\[[dependency-cultures]\] lists "granularity of the unit" as a culture lever: package-shaped ecosystems produce package-sized dependencies, function-shaped ecosystems don't.
Trimming doesn't require redesigning what a package is.
Upstream can still publish coarse, kitchen-sink packages.
But what actually lands in your repo, post-trim, is function-shaped, because reachability analysis doesn't care about the package's intended unit, only about what you call.
So you might get function-granularity dependencies for free, as a side effect of the trim step, without ever having to convince an ecosystem to publish that way.
That feels like it might be the actual point of this idea, more than the "fewer dependencies" framing on the tin.

## Where this fights the other doodles in \[[dependency-cultures]\]

Three different answers to "how do you take a dependency" now exist across these two notes: capability-checked packages, content-addressed hashes, and fork-and-trim.
They are not obviously compatible.

Capability checking wants the dependency to declare `requires: net, alloc` and have the compiler verify it against the *original* package.
Once you've trimmed, what are you checking capabilities against, the original untrimmed source (which you no longer ship) or the trimmed derivative (which might have accidentally dropped the code path that used `net`, silently changing the declared capability)?

Content-addressing wants a dependency to be a hash, so two projects that depend on "the same" thing dedupe.
Trimming produces a different artifact per project, keyed off each project's own call sites, so the same upstream package trimmed by two different projects is not the same hash, and dedup across projects is gone by construction.
That might be fine (dedup was mostly a build-time/disk-space concern) or might be the whole point being thrown away (content-addressing's other promise, "reproducible builds are the only thing that can happen," survives fine, it's just per-project reproducibility, not cross-project sharing).

Not resolving this now.
Flagging it because a language that wants all three properties (capabilities, content addressing, fork-and-trim) at once is presuming they compose, and I have no evidence they do.

## Dead ends, recorded so I stop rediscovering them

- **"Trim by hand every time."** Forces real engagement with what you're importing, in the spirit of the C vendoring culture already in \[[dependency-cultures]\]. Dies on scale: the moment a dependency has more than a handful of call sites, hand-trimming is just re-implementing reachability analysis badly, and worse, inconsistently between people.
- **"Trim once, never re-trim."** Simplest possible version, and it just becomes permanent forking with a nicer first diff. All the update problems above still exist, they just have no tooling answer at all. Worse than the tooled version, not obviously worse than the ecosystem status quo.

## Where this leaves me, unresolved

The "add is cheap, everything else is expensive" cost structure that \[[dependency-cultures]\] identifies as the root cause of npm-style culture is directly attacked here: `lang vendor add` is deliberately not cheap, it produces a diff you have to look at before the dependency exists at all.
That is the strongest part of the idea.

The weakest part is "update", which trades a solved problem (semver ranges, however badly they work) for an unsolved one (mechanical re-trim plus manual merge of local patches), and I don't have a story for why that trade is obviously worth it at scale, only that it might be worth it for a language whose whole premise is having very few dependencies in the first place, where "at scale" may never actually arrive.
