# Dependency cultures

Riffing on the premise of "Dependency Cultures" (Richard Feldman, SSW 2026).
Only the title was available when this note was written, so nothing here is a summary of the talk and nothing should be attributed to it.
Zero standing, presumes outcomes of open decision records freely.

## The premise I want to steal

Ecosystems have wildly different norms about taking a dependency.
A JS developer adds forty packages to render a date and feels nothing.
A Go developer thinks about it.
A C developer vendors the file and renames the symbols.
A Nix user pins the world and never thinks about drift again.

The interesting claim is that this is not a difference in the *people*.
It is a difference in what each language made cheap.
Culture is the downstream sediment of a handful of very early technical decisions, and the sediment is basically permanent by the time anyone notices it forming.

If that is true, then "what dependency culture do I want" is a language design question that gets answered whether or not it gets asked.
Answering it late means inheriting someone else's answer.

## The levers that actually set the culture

Candidate list, roughly in order of how much I think each one matters:

1. **Size of the standard library.** Small stdlib means the ecosystem fills the gap, and it fills it with thousands of tiny packages. Large stdlib means fewer deps and a museum of 2009 API design that can never be removed.
1. **Friction of publishing.** `npm publish` in four seconds produces left-pad. A curated registry with review produces Debian, and produces packages that are three years stale.
1. **Granularity of the unit.** If the unit of dependency is a package, people ship packages. If it is a single function or a single module, the graph looks completely different, and so does the diff.
1. **Version resolution model.** SemVer plus a solver invites version ranges, which invites "works on my machine". MVS (Go) invites pinning. Content addressing (Unison) deletes the question.
1. **Whether adding a dependency shows up in a diff and costs something visible.** A one-line manifest edit that pulls 300 transitive packages is the single most asymmetric action in modern software.

Note that only (1) and (3) are language design in the narrow sense.
The rest is tooling, and the tooling is where the culture actually gets set, which is uncomfortable if you think of the language as the interesting part.

## What is culture optimizing for, really

Every dependency decision is a question of who pays and when.

- **Author time** goes down with more dependencies. This is the only cost anyone feels while typing.
- **Reader time** goes up. Understanding a program means understanding its deps, and nobody reads them.
- **Build time** goes up, invisibly, in increments too small to ever blame on one package.
- **Audit and supply-chain risk** goes up superlinearly with the transitive count.
- **Maintenance time** goes up but is deferred past the point where the decision can be revisited.

The npm culture is not irrational.
It is a perfectly rational response to a cost structure where the only bill that arrives immediately is author time.
So the design question is not "how do I make people take fewer dependencies", it is "how do I move some other bill forward in time so it arrives while the decision is still being made".

That reframing feels like the useful part of this note.

## Doodle: make the bill arrive at the call site

What if the compiler charged you, out loud, on every build?

```
$ lang build
  building app
    core          0.4s     12 KB
    http          2.1s    380 KB   (+7 transitive)
    json          0.3s     41 KB
  total          2.8s    433 KB

  http is 76% of your build time and 88% of your binary.
```

Nothing enforced, nothing forbidden, just the number on the screen every single time.
Cheap to build, and it is the only mechanism here that does not require winning an argument with the user.

Prior art in spirit: the wishlist already wants compile times treated as a budget that is measured continuously.
This is that, attributed per dependency instead of in aggregate.

## Doodle: dependencies declare their capabilities, and the compiler checks

```
package http {
  requires: net, alloc
}

package left-pad {
  requires: ()
}
```

A package that declares no capabilities cannot open a file, cannot spawn a process, cannot phone home, and the compiler proves it rather than the registry promising it.
Then the transitive count stops being the scary number.
Four hundred pure packages is fine.
One package that quietly requires `net` and `exec` is the thing worth a code review.

This is the same shape as the capability idea in \[[file-primitive]\], generalized from files to the whole dependency graph, which is a hint that "capability" might want to be one concept in this language rather than two.

The hard part is not the checking, it is that `alloc` and `time` and `rand` are all capabilities too, and if the annotation burden is high enough nobody writes it, and if there is a default it will be "everything".

## Doodle: no versions, only content

Steal Unison.
A dependency is a hash.
Names are a local alias for hashes, versions are a social fiction layered on top, and "upgrading" is explicitly rewriting your aliases rather than a solver silently picking for you.

Consequences I like:

- Diamond dependencies stop existing as a category. Two hashes are two different things and both can be present.
- Reproducible builds are not a feature you add, they are the only thing that can happen. The wishlist wants this from day one.
- `Blob`-style content addressing shows up in two unrelated notes now, so it may be a real primitive rather than a coincidence.

Consequences I do not like:

- Security patches. If a name is just an alias for a hash, "everyone who depends on this must move" has no mechanism. The whole industry runs on that mechanism working badly, but it does run on it.
- Humans cannot read hashes, so the alias layer comes back, and the alias layer is where all the version problems lived in the first place.

Not dead, but it needs an answer to patching before it is more than a doodle.

## Dead ends, recorded so I stop rediscovering them

- **"Just ship a huge stdlib."** Dies on the fact that a stdlib is where APIs go to become permanent. Every large stdlib has a `datetime` module its own maintainers tell you not to use. Deprecation with no removal is the tax, and it is paid forever.
- **"Forbid transitive dependencies."** You may only depend on packages that have zero dependencies of their own. Deliciously extreme, forces the graph flat, and dies instantly because everyone just vendors and the graph becomes invisible instead of shallow. Worse than what it replaced.
- **"Registry review board."** Solves quality and supply chain, dies on throughput, and the language has one maintainer.

## Where this leaves me, unresolved

Three tensions I cannot collapse yet:

The **stdlib size** question and the **culture** question are the same question wearing different clothes, and I want a third option that is neither "tiny stdlib, huge ecosystem" nor "huge stdlib, dead APIs".
Maybe a stdlib that is versioned and replaceable like any other package, so that removal is possible.
Maybe that just relocates the museum.

**Capabilities** are the most interesting lever here and also the one most likely to be annotated into meaninglessness.
The design lives or dies on defaults, and I have no candidate default I believe in.

And the honest one: all of this presumes there is an ecosystem, which presumes there are users, which is a long way from a repository with no compiler in it.
The reason to think about it now anyway is exactly the premise at the top.
By the time the culture is visible, the decisions that caused it are ten years old.
