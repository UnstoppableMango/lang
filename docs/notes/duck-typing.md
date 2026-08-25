# Duck typing

If it walks like a duck and quacks like a duck, it's a duck.
The classic framing is runtime: you don't check the type, you just call the method and let it fail if it can't.

But "duck typing" as a *feeling* is separable from dynamic dispatch.
Structural typing gives you the same ergonomics statically: if a value has the right shape, it satisfies the interface, no declared relationship required.
Go's interfaces are the obvious existence proof: implicit satisfaction, checked at compile time, no `implements` keyword anywhere.
So the real question for this language isn't "duck typing yes or no", it's "structural or nominal", and separately "static or dynamic".
Those are two axes, not one.

## Four quadrants

- Nominal + static: Java, C#, Rust traits (mostly). You declare `impl Foo for Bar`. Safe, explicit, but you can't retroactively satisfy an interface you don't own without a wrapper.
- Structural + static: Go, TypeScript, OCaml's row-typed variants sort of. Shape-matching, checked ahead of time. Feels like duck typing but with a compiler backstop.
- Nominal + dynamic: rare combo. Python-with-ABCs-that-actually-check maybe? Ruby's `is_a?` checks if you use them. Mostly nobody wants this, it's the worst of both.
- Structural + dynamic: Python, JS, Ruby by convention (not enforced, just cultural). This is "duck typing" in the classic sense, no compiler in sight.

## What actually attracts people to duck typing

Not the lack of types.
It's two things bundled together that don't have to travel together:

1. **Retroactive conformance.** I have a `Logger` interface and a third-party `Widget` type, and `Widget` already has a `Log(string)` method, so it just works. No adapter, no `impl` block I have to write and maintain. This is the structural part.
2. **Minimal interface surface.** You only need the methods you call, not the whole type. A function that calls `.Read()` doesn't care what else the argument can do. This is more about interface segregation than typing discipline at all, arguably you get this in nominal systems too if people write small interfaces (Go proves this, ironically, despite being structural — the discipline could exist in Rust with tiny traits, but culturally doesn't).

So maybe the wishlist item isn't "duck typing", it's "structural interface satisfaction, checked statically, with narrow inferred interfaces at call sites."
That's a Go-like or OCaml-row-polymorphism-like thing, not a Python-like thing.

## Where it gets weird: what "quacks" means for a value type language

If arena/value semantics wins the memory model decision (see [[arena-memory-model]]), does structural typing even look the same?
In Go, interface satisfaction is about pointer-or-value receivers and it gets subtle (a `T` satisfies an interface if all methods have value receivers, `*T` needs `*T` methods too).
If this language leans hard into values-by-default, method dispatch through an interface might need to always copy, or the interface itself might need to specify by-ref vs by-value at the boundary.
Untouched territory, not resolving it here.

## Dead end: full Python-style duck typing (structural + dynamic)

Tempting because it's maximally flexible and "just works" for scripting-style code.
Died on: it fights every other instinct so far toward compile-time guarantees (see [[no-null-type-or-representation]], [[no-exceptions-explicit-errors]] for the general vibe of "catch it before runtime").
A dynamically-typed duck-typing system means `.Quack()` on a value that doesn't have it is a runtime crash, which is exactly the failure mode this project seems to want to design away from by default.
Not ruling it out as an *opt-in* mode (see [[compiled-first-scripting-alt]]) but it's not the default-path answer.

## Open tension

Structural static typing is easy to want and hard to keep coherent once generics show up.
TypeScript's structural typing gets weird fast with excess property checks, variance, and the fact that `{}` basically means "anything with no required shape."
Go dodges most of this by not having a real generics-shaped structural system until recently, and even now interfaces and generics interact in ways people find surprising (type sets, `~T`, comparable).
No conclusion here, just a flag: "structural typing" sounds simple in the two-line pitch and stops being simple the moment you add generic interfaces, variance, or nested structural constraints.
Worth returning to once (or if) generics gets its own note.
