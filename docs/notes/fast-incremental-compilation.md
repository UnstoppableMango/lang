fast incremental compilation

the pitch: editing one function and re-running `mango build` should feel instant, not "wait for the whole crate to re-typecheck."
worth separating two things that get conflated under "incremental compilation":

1. incremental *rebuilds* across process invocations (edit a file, run the CLI again, only the changed part re-executes).
1. incremental *analysis* inside one long-lived process (a language server watching keystrokes, needs sub-100ms turnaround on every reparse).

these have wildly different constraints.
(1) can afford disk-backed caches, hashing, serialization.
(2) needs in-memory structures that survive edits without a serialize/deserialize round trip.
if the language ships both a batch compiler and a language server (it will, eventually, everyone does), do they share one incremental engine or two?

## approach A: file-level caching (the "make" model)

the simplest incremental story: hash each source file, if the hash is unchanged and its dependencies' hashes are unchanged, skip recompiling it, reuse the cached object/IR.

```
mango build
  for each file in dependency order:
    if hash(file) == cached_hash(file) and all deps unchanged:
      skip, reuse cached output
    else:
      recompile file, update cache
```

this is what most build systems do (make, and by extension most C/C++ toolchains) and it's dead simple to implement and reason about.
the failure mode is well known though: change one function in a file with 2000 lines and 40 unrelated functions, the whole file re-typechecks, because the cache granularity is "file," not "function."
if \[[directory-scoped-modules]\] lands, this gets worse in one specific way: a directory *is* the module, so does the cache granularity become "directory" instead of "file"?
that would mean editing `token.mn` invalidates the cached analysis of every sibling file in `parser/`, even ones that don't reference `Token` at all, purely because they're nominally "the same module" and module-level facts (name resolution, visibility) got computed as one unit.
that's a real tension between directory-scoping's zero-ceremony grouping and fine-grained incrementality, worth flagging if both ideas survive to sketch stage.

## approach B: query-based / demand-driven (the salsa model)

rustc's older architecture was file-level-ish; rust-analyzer and modern rustc internals both moved to a *query system* (salsa is the actual crate, but the pattern predates it, e.g. Adapton, self-adjusting computation more generally).
the idea: instead of a fixed pipeline (parse -> resolve -> typecheck -> codegen), define named queries with memoized results, and let each query pull its inputs on demand.

```
query type_of(item: ItemId) -> Type
query resolve(name: Name, scope: ScopeId) -> ItemId
query hir_of(file: FileId) -> Hir
```

when a file changes, you don't invalidate "the file's compiled output," you invalidate exactly the queries whose *inputs* changed, and every downstream query that transitively depended on them gets recomputed lazily, on next demand, not eagerly.
this gets you function-level (really, query-level) granularity for free, at the cost of a much more complex compiler architecture: everything has to be expressed as a pure function of its inputs, memoized, with an invalidation graph tracked automatically.
this is a much bigger up-front architectural commitment than approach A.
it basically decides "the compiler is a database" as a foundational shape, which interacts with the *implementation language* decision (a decision record still open per AGENTS.md): salsa-style incrementality is comfortable in Rust, has been ported to other languages, but writing one from scratch is a project unto itself, not a feature you bolt on later.

worth naming directly: this is arguably too big a decision to make as a language design note.
"is the compiler internally a query database" is closer to an implementation strategy than a language feature, so maybe this note's job is just to flag the *implication for language design*, not to pick the compiler architecture.
what does the language have to guarantee for either approach to work well?

## what the language can do to make either approach easier

this is probably the actually load-bearing question for a design note, versus "which compiler architecture," which is an implementation detail.

**purity / referential transparency of semantic analysis.**
if resolving a name or computing a type can have side effects (macros with arbitrary compile-time code execution, e.g.), memoization gets a lot harder to trust, because "same inputs" no longer guarantees "same outputs."
languages with hygienic-but-effectful macro systems (Rust's proc macros doing arbitrary computation) have had real bugs here.
if macros/metaprogramming stay off the table, or get sharply restricted to something declarative, that's a point in incremental compilation's favor, independent of whether that restriction gets justified any other way.

**explicit, small dependency edges between items.**
if `type_of(f)` depends on `type_of(g)` only when `f` literally calls `g`, invalidation is precise.
if there's some ambient global inference step where everything can affect everything (unrestricted Hindley-Milner unification across a whole compilation unit is somewhat notorious for this, a single change can shift inferred types arbitrarily far away), incrementality degrades toward "recompute everything," because the dependency graph is effectively complete.
this is a real tension against the wishlist's "type inference strong enough that annotations are for communication, not for the compiler": whole-program HM inference and fine-grained incremental recompilation pull in somewhat opposite directions.
local type inference (infer within a function body, require signatures at function boundaries, à la Rust, TypeScript-with-noImplicitAny-ish, most MLs' module boundaries) buys back incrementality by making the function signature the cache-invalidation boundary: change a function body, only that function's callers who depend on inference-through-the-body (none, if the signature is explicit) need to recheck.
this seems like the actual design lever: how much can be inferred *within* an item without touching how other items get analyzed.

**module boundaries as cache boundaries.**
\[[directory-scoped-modules]\]'s open tension (directory-as-module vs finer per-item boundaries) turns out to matter here too, not just for visibility.
a module system where "what's visible to importers" is a small, explicit, syntactically obvious set (option B in that note, `pub` per item) gives the compiler a natural incremental boundary: recompiling a module's *internals* never has to re-typecheck importers, only recompiling its *exported* signatures does.
if visibility is implicit or capitalization-based, the compiler has to conservatively assume more surface area changed than actually did.

## approach C: don't build a novel incremental engine, reuse LLVM/existing infra

tangent worth naming and setting aside: if \[[compiled-first-scripting-alt]\] wants a `mango run` no-ceremony mode, does that mode even want incrementality, or does it want *fast enough from-scratch compilation that caching doesn't matter*?
a sufficiently fast non-incremental compiler (think: single-pass, no separate optimization pass, tree-walk-adjacent codegen for the "run a script" case) sidesteps the whole caching problem for the common case, and only the "large project, `mango build --release`" case needs real incrementality.
this bifurcation (interpreter-fast-path for `run`, incremental-cached-path for `build`) mirrors how a lot of language ecosystems actually ended up (a bytecode VM for dev-loop speed, AOT/LTO for release builds) without anyone deciding it as one unified strategy up front.
worth remembering as a possible "you don't need one answer" escape hatch if approach A vs B stalls out as a false binary.

## content hashing vs timestamps

small implementation-flavored tangent, but it affects UX: mtimes are what make/most naive tools use, and they lie constantly (git checkout resets mtimes to checkout time regardless of content, clock skew across machines/CI runners, editors that touch-and-rewrite on save).
content hashing (hash the actual bytes, compare hashes) is strictly more correct and is what modern build systems (Bazel, Buck2, Nix itself) converged on, at the cost of hashing being nontrivially slower than a stat() call on very large trees.
given the project already leans on Nix for everything (per AGENTS.md's `nix flake check`/`nix build`), content-addressing as the default mental model is already the house style elsewhere in this repo's tooling, so it'd be a strange inconsistency for the compiler's own incrementality to fall back to mtimes.
low-stakes note, but "prefer content hashing" seems like a easy default to write down now before an mtime-based cache gets prototyped out of habit.

## dead end considered: incremental as a bolt-on optimization pass, decided later

considered treating "make the compiler incremental" as pure implementation work to defer entirely, no language-design implications, revisit once there's an actual compiler to profile.
the reasoning against that: the *cost* of incrementality gets paid by language design choices made way upstream of any compiler existing (macro purity, inference scope, module visibility, all discussed above), and those choices are exactly the kind of foundational thing this project is trying to make deliberately rather than "by accident of the first implementation" (echoing the wishlist's own language about the memory strategy decision).
if inference is unrestricted whole-program HM and macros are effectful, no amount of clever query-caching architecture bolted on afterward fixes the invalidation graph being effectively complete.
so this isn't purely deferrable, even though the actual *engine* (salsa-alike vs make-alike vs skip-it-be-fast-instead) probably is.

## open tension, not resolved here

whole-program type inference (maximally ergonomic, minimally annotated) vs fine-grained incremental recompilation (wants small, explicit, syntactic dependency edges) are in real tension, and the wishlist currently has an entry wanting the former without acknowledging the tradeoff.
doesn't mean "give up on strong inference," but any future inference design doc probably owes a paragraph on where inference is scoped (function-local? module-local? whole-program?) with this incrementality question as one of the reasons that scoping choice matters, not just an ergonomics question.
