# Nix-first build system

This repo already builds itself with Nix (`nix flake check`, `nix build .#`, `nix fmt`).
That's process tooling for the design-docs repo, not a claim about the language.
The question worth playing with: what if that's not an accident of how this project happens to be developed, but the actual answer to "how does *this language* build things"?
Not "the compiler is written in a build tool" (Nix wrapping an arbitrary Makefile-shaped compiler, which is just infrastructure), but "the language's own build model *is* the Nix model," content-addressed derivations all the way down.
Zero standing, presumes outcomes of open decision records freely.

## The provocation

[[dev-tooling-philosophy]] argues for a Go-shaped `lang build`/`lang test`/`lang fmt` monolith, one binary, one noun, first-party everything.
That whole note assumes the toolchain is bespoke: this language invents its own build graph, its own caching, its own reproducibility story.

Nix already solved reproducible, content-addressed, cached builds, twenty years of prior art, a real substrate people already have opinions about.
What if `lang build` is not a bespoke build graph at all, it's a thin frontend that emits a Nix derivation, and the actual work (caching, parallelism, remote build, reproducibility) is just Nix doing what Nix does?

```
$ lang build
  (generates a derivation, hands it to nix-build, gets a store path back)
  /nix/store/9fjp2q...-myprogram
```

The appeal: skip reinventing content-addressed build caching, which is genuinely hard to get right (ask Bazel, ask Buck, ask every language that's tried to bolt on incremental builds after the fact, see [[fast-incremental-compilation]]).
Nix's store is already a working answer to "never rebuild the same thing twice, ever, across projects, across machines."

## Doodle: what "package" even means if it's a derivation

If a build artifact is a Nix derivation, a dependency is not "a name and a version," it's a store path, content-addressed by its inputs.
This is [[fork-and-trim-dependencies]]'s content-addressing tension made literal instead of hypothetical, Nix already has the exact mechanism that note was gesturing at.

```
lang.toml (or whatever the manifest is)
  http = "nixpkgs#some-http-lib"     # or a flake ref, or a store path directly
```

Two projects depending on "the same" library dedupe automatically, for free, because they resolve to the same store path if the inputs match.
That's the cross-project dedup that [[fork-and-trim-dependencies]] noted trimming throws away.
Interesting: Nix-first and fork-and-trim might be in real tension, not just adjacent doodles.
Trimming produces a bespoke per-project artifact by construction; Nix wants to hash the *whole* input and share it whenever two consumers agree.
You can't have "the artifact is exactly what I use, nothing more" and "the artifact is exactly what upstream published, shared across everyone" as the same object.
Pick one, or accept that fork-and-trim happens *before* the Nix layer (trim first, then let the trimmed result be its own content-addressed thing, dedup only across projects that happened to trim identically, which will basically never happen and so dedup silently stops mattering).

## Doodle: the whole toolchain as flakes, not as a monolith binary

[[dev-tooling-philosophy]]'s central bet is "one binary, one noun, so tooling doesn't fragment the way C's make/cmake/meson/bazel did."
Nix-first cuts against that bet directly: instead of `lang fmt` being a subcommand of one blessed binary, it's a flake output.

```
lang#fmt
lang#test
lang#build
lang#lint
```

Each one is independently versionable, independently swappable, independently overridable by a downstream flake (`lang.override { formatter = myFmt; }`), which is exactly the "plugin ecosystem" outcome that [[dev-tooling-philosophy]] flagged as the thing that kills `gofmt`'s zero-config promise ("no config, but you can just not use gofmt, which is worse than either").
Nix's whole culture is "everything is overridable," which is the opposite instinct from Go's "everything is fixed, that's the point."
Not obviously wrong, just a different bet: Nix's ecosystem has survived fine with pervasive overridability because the *reproducibility* is what's load-bearing, not the *uniformity*.
Maybe that's a real alternative to Go's answer: you don't need one true formatter if every build is reproducible regardless of which formatter produced the input, you just need the build graph to be honest about what actually ran.
Or maybe that's coping, and code review style arguments show up anyway, just relocated from "which formatter" to "whose flake override wins."

## Doodle: build cost reporting is Nix's for free

[[dev-tooling-philosophy]] doodled a chatty `lang build` that prints per-phase timing.
Nix already has this shape baked in, mostly as noise (`these 14 derivations will be built, these 40 will be fetched`), not as the polished per-phase story that doodle wanted.

```
$ lang build
these 3 derivations will be built:
  /nix/store/...-lang-parse.drv
  /nix/store/...-lang-typecheck.drv
  /nix/store/...-lang-codegen.drv
building '/nix/store/...-lang-parse.drv'...
```

That's honest in a way hand-rolled phase timing isn't: it's telling you exactly what's cached and what isn't, derivation by derivation, which answers "why is this slow" better than a timer ever could, since the timer doesn't tell you *which part missed cache*.
But it's also Nix's UX, not language UX, derivation hashes are not phase names, and translating one into the other (making `lang-typecheck.drv` read as "typecheck 0.41s" to a user who's never heard of Nix) is nontrivial glue, not a free lunch.

## Where this fights itself: does the language now require Nix

The load-bearing risk, named directly: if `lang build` is a thin wrapper around `nix-build`, does the language require a Nix install to build anything, ever?
That's an enormous scope decision hiding inside what looked like a build-tooling doodle.
Go's whole pitch (see [[dev-tooling-philosophy]] again) is *one command, no separate thing to install first*.
`go install` doesn't ask you to already have some other package manager on your machine.
A Nix-first language asks exactly that, and Nix's own onboarding cost (multi-user install, daemon, flakes still technically experimental years in) is not nothing, it is a genuinely higher first-five-minutes cost than `curl | sh`-ing a Go-shaped toolchain.

Possible escape: vendor a tiny embedded store implementation, a "Nix, but only the 5% this language needs, statically linked into the compiler, no daemon, no multi-user install."
That's plausible (there are already projects chipping at "Nix without the ceremony") but it's a research bet, not an available-today substrate, and it reintroduces exactly the "reinvent the build graph" cost this whole doodle was trying to dodge by riding on existing Nix.
So the doodle's main appeal, don't reinvent content-addressed caching, might only be free if you're willing to pay Nix's actual onboarding cost, and might cost just as much as reinventing it if you try to hide that cost from users.

## Doodle: transparent Nix, the compiler leans on it but never says its name

Reframing the "does the language now require Nix" risk above: the lean worth exploring is not "expose Nix to users" vs. "don't use Nix at all," it's "the compiler's own internals lean on Nix" while the end user never types `nix` or knows the store exists.

```
$ lang build myprogram.gos
building...
myprogram
```

No `.drv`, no store path in the happy-path output, no flake syntax in the manifest.
Underneath, `lang build` still shells out to a vendored/embedded Nix (or links against the Nix libraries directly, evaluator and store as a linked dependency rather than a subprocess) and the actual caching, content-addressing, and dependency resolution the earlier doodles wanted is still real Nix doing real work.
The user-facing error path is the hard part, not the happy path: a broken build has to surface as a language-shaped error ("type mismatch in `parse.gos:12`"), never as a raw Nix eval trace, or the abstraction leaks exactly where it matters most (see the "friendly wrapper over a scary tool" dead end below, this doodle is a bet that the leak is containable if the wrapper owns error translation as a first-class concern from day one, not bolted on after the fact).

This changes which of the "Dead ends" below is actually live.
"Language-level package manager on top, Nix underneath, user never sees Nix" was recorded as a dead end because the wrapper either leaks the underlying tool's concepts on first failure, or reimplements enough of Nix to hide them.
Leaning toward "compiler leans on Nix, kept transparent" means picking that dead end back up and betting the leak is survivable if the compiler, not a bolted-on wrapper, owns the Nix boundary from the start, same "shared library, not shelled-out tool" instinct as the AST-sharing argument in [[dev-tooling-philosophy]].
Not proven, just the direction currently being explored, and it reopens rather than closes the "does this need a Nix-without-the-ceremony embedded store" thread above, that's now load-bearing instead of an escape hatch.

## Dead ends, recorded so I stop rediscovering them

- **"Just shell out to `nix build` and call it done."** Works today, for this repo, as infrastructure. Does not answer the actual language-design question, which is whether *programs written in this language* declare their dependencies and build graph in Nix's terms, not whether the design-docs repo happens to use Nix to lint itself.
- **"Language-level package manager on top, Nix underneath, user never sees Nix."** Tempting, but it's the same shape as every "friendly wrapper over a scary tool" project (see: every Docker Compose-for-X, every Terraform-wrapper-du-jour). The wrapper either leaks the underlying tool's concepts the first time something goes wrong (a Nix eval error is not a friendly error), or it has to reimplement enough of Nix's semantics to hide them, at which point why depend on Nix at all instead of just building the bespoke thing [[dev-tooling-philosophy]] already sketched.

## Threads worth pulling later

- Whether "content-addressed dependency" as a language-level concept (not just a build-system-level one) is actually the same idea [[dependency-cultures]] and [[fork-and-trim-dependencies]] were circling, just arrived at from the infrastructure side instead of the language-design side. If so, three notes have been converging on one idea independently, worth a unifying pass once this stops being play.
- Cross-compilation (on [[dev-tooling-philosophy]]'s checklist) is close to free if the build graph is Nix's, since Nix's whole cross-compilation story is mature and unrelated to this language's own effort. Would need to check whether that's actually true or just sounds true.
- Whether "the compiler is a Nix derivation" and "the compiler is a library other tools import" (the shared-AST doodle in [[dev-tooling-philosophy]]) are compatible at all, or whether Nix's process-boundary-per-derivation model fights the "one shared in-memory AST across fmt/vet/lsp" goal by construction, since derivations don't share memory, they share store paths.
