# Development tooling philosophy

Riffing on "comprehensive tooling, Go-style" as a design stance rather than a feature list.
\[[go-simplicity]\] already touched this in one section ("tooling as part of the simplicity story").
This note is the expansion: what would it actually mean for this language to ship a `go`-shaped toolchain, and is that even the right shape to copy.

## The core claim to test

Go's actual innovation was not any single tool.
It was collapsing "install a compiler" and "install a development environment" into the same command.
`go install` gets you a compiler, a formatter, a test runner, a vet pass, a doc server, a module system, and a profiler, all speaking the same flags, all versioned together, all maintained by the same team that ships the language.

The alternative, which is most languages, is: the language ships a compiler, and the community assembles the rest.
C has a compiler and then a fifty-year archaeology dig of make/cmake/autotools/meson/bazel, each solving the same problem slightly differently.
JS has a compiler-ish thing (several) and then webpack/rollup/esbuild/vite/turbopack, a graveyard of build tools that each promised to be the last one.
Python has pip/poetry/pipenv/conda/uv, still not settled, decades in.

The claim worth testing: fragmentation is not a maturity problem that these ecosystems will eventually grow out of.
It is what happens by default when the language does not claim the space first.
Once a community has three build tools, a fourth official one does not consolidate anything, it just adds a fourth.
The window to prevent this is early, maybe only at 1.0, maybe only before 1.0.

## What "comprehensive" would actually cover

Listing the surface area Go claims, as a checklist to react to rather than a spec:

- **Build**: one command, no makefile, no build DSL to learn separately from the language.
- **Format**: one canonical style, zero configuration, so format arguments cannot happen in code review.
- **Test**: a test runner built into the toolchain, not a third-party framework with its own assertion DSL.
- **Vet / lint**: a static analysis pass that ships with the compiler and knows the language's specific footguns, not a generic linter bolted on later.
- **Docs**: documentation generated from source and comments, viewable without a separate static site generator.
- **Dependency management**: a module system with one resolution algorithm, not a choice of five.
- **Package registry story**: where code comes from when you depend on someone else's, and who runs it.
- **Cross-compilation**: build for a different OS/arch without installing a different toolchain.
- **Profiling / tracing**: performance investigation tools that speak the runtime's own format, not an external agent bolted on.
- **Editor / LSP**: a language server planned alongside the compiler (already on the wishlist), so tooling and language evolve together instead of the LSP forever chasing the compiler's internals through a second implementation.
- **Release / versioning**: how the toolchain itself is versioned and how that interacts with language version and module version (Go's own answer here, `go.mod`'s `go 1.21` directive, took over a decade to arrive and is still awkward).

That is eleven surface areas.
Go's actual team is not eleven times the size of a normal compiler team, which suggests either heavy internal reuse (a lot of these share the same parser and AST) or that some of these are much cheaper than they look once the first two or three exist.

## Doodle: the toolchain as one verb space

```
lang build
lang test
lang fmt
lang vet
lang doc
lang run
lang mod tidy
lang mod why net/http
```

The interesting property is not the individual commands, it's that they are all subcommands of one noun.
Compare to the C world's `gcc`, `make`, `ctest`, `clang-format`, `clang-tidy`, `doxygen`, `pkg-config`, each a separate binary, separate version, separate config file format, separate flag conventions (`-o` means different things in different tools).
One binary means one help system, one version number that means something, one thing to install.

Cost: the binary becomes a monolith that has to ship on every language release, and every subcommand's roadmap is now coupled to the compiler's release cadence.
Go mod took years to land partly because it had to happen inside the main toolchain's release train rather than as an independent package someone could iterate on in isolation and abandon without consequence.

## Doodle: tools that share one data structure

The reason Go's tools compose well is not the CLI, it's that `go vet`, `gofmt`, `gopls`, and the compiler all parse to the same AST shape and mostly share `go/packages` for loading a build graph.
A vet check and a compiler error are the same kind of finding running through the same pipeline, just at different confidence thresholds.

For this language, the doodle is: define the AST/IR once, in a form the compiler, formatter, linter, and language server all import as a library, not a form each tool re-derives from source text independently.
This is not really a CLI decision, it's an architecture decision that has to happen before any of the CLI surface can be built cheaply.
If the compiler is a library first and a CLI second, the eleven tools above become thin wrappers.
If the compiler is a monolithic binary from day one, each tool has to shell out and re-parse, and now there are eleven parsers to keep in sync (see: the years-long lag of third-party formatters and linters for languages that didn't expose their parser as a library).

This connects to the wishlist's "language server planned alongside the compiler, not bolted on afterward": that entry is really asking for exactly this, a shared library boundary, phrased as a tooling item.

## Doodle: the toolchain reports what it's doing, not just the result

Stolen from the \[[dependency-cultures]\] build-cost doodle, generalized past dependencies to the whole toolchain:

```
$ lang build
  parsing        0.02s
  typecheck      0.41s
  codegen        0.88s
  link           0.12s
  total          1.43s   (+0.3s vs last build)

$ lang test
  12 packages, 340 tests, 338 passed, 2 failed  (2.1s)
  FAIL  parser/lexer_test.gos:44  unexpected token '}'
```

Go's actual tools are fairly quiet by default, arguably too quiet, `go build` on success says nothing at all.
That's a legitimate design choice (no output is good output, silence means success), but it also means the cost information from the dependency-cultures doodle would need to live somewhere, and "nowhere, by default" is Go's actual answer.
Worth deciding on purpose rather than by silence: is this language's tooling chatty-by-default (Cargo's colored, verbose-ish output) or silent-by-default (Go's near-total quiet on success)?
Those produce very different developer-experience cultures and it's not obvious Go's is the one to copy just because the rest of this note is Go-flavored.

## Where "comprehensive" fights itself

Comprehensive and opinionated are almost the same axis, and that's the thing to be honest about.
`gofmt`'s zero-configuration stance is only tolerable because there is exactly one, official, always-installed formatter.
The moment there's a plugin ecosystem or a second formatter, "no config" becomes "no config, but you can just not use gofmt," which is a worse outcome than either configurability or true uniformity.
So a comprehensive toolchain either has to stay a closed, first-party system forever, or it has to accept that "comprehensive" decays into "comprehensive-ish, plus whatever the ecosystem bolted on," which is exactly the fragmented outcome this whole note is trying to avoid.

Also worth naming: Go's toolchain is comprehensive partly because Go the language stayed small for a decade (see \[[go-simplicity]\]).
A small language surface is cheap to build eleven tools around.
If this language ends up with a large or unsettled surface (macros, comptime, an effect system, whatever lands from the wishlist), "one canonical formatter, one canonical vet pass" gets harder in direct proportion to how much the language itself still has open decision records.
Comprehensive tooling might be a *late* project, achievable only once the language design has mostly stopped moving, rather than an early one.
That's in tension with "planned alongside the compiler, not bolted on," which wants it early.
Not resolved: maybe the resolution is that the *architecture* (shared AST/IR as a library) has to be early, while the *comprehensiveness* (eleven polished subcommands) is necessarily late, and conflating "start the plumbing early" with "ship all eleven tools early" is the actual mistake to avoid.

## Dead ends, recorded so I stop rediscovering them

- **"Ship the formatter as a separate optional install."** Kills the entire point. An optional gofmt is a gofmt nobody's codebase agrees on, which is the pre-Go status quo. Formatter has to be inseparable from `lang build` or it doesn't do the cultural work.
- **"Let the community pick the build tool, keep the compiler minimal."** This is explicitly the C/JS/Python path this note opened by criticizing. Tempting because it's less upfront work, and that's exactly the trap, the upfront work is the whole value.
- **"Configurable formatter with sane defaults."** Sounds like a reasonable compromise, dies the same way `.prettierrc` dies: the existence of the knob guarantees someone turns it, and now code review has style arguments again. Either no config or don't bother.

## Threads worth pulling later

- Whether the toolchain binary should itself be self-hosting (written in this language) before or after 1.0, since a self-hosted `lang fmt` is the strongest possible proof the language is pleasant to write real programs in.
- What "vet" actually checks for *this* language specifically, once the memory model and error story are decided, since Go's vet checks (printf format strings, struct tags, lock copying) are downstream of Go's specific footguns, not a generic template.
- How module/dependency tooling here reconciles with the capability-checking doodle in \[[dependency-cultures]\], since "lang mod tidy" and "prove this package can't open a socket" are two different tools today in every existing language and maybe shouldn't be here.
- Release cadence for the toolchain itself: tied to language version (Go's model, eventually) or decoupled (each tool versioned independently, more Rust/Cargo-shaped)?
