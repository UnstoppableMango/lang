# Arena memory model

Poking at whether "arena allocation" could be THE memory strategy, not just an optimization bolted onto GC or ownership later.

## The pitch

Bump-allocate into a region.
Free the whole region at once instead of tracking individual object lifetimes.
No refcounting, no mark-and-sweep, no borrow checker fighting the programmer.
Cheap allocation (pointer bump), cheap deallocation (drop the whole block), and the cost is paid up front in how you structure the program around regions instead of paid continuously at runtime.

This is attractive for a language "designed from scratch, deliberately slowly" because it sidesteps two of the biggest unresolved fights (GC vs manual vs ownership) by picking a third thing that is arguably simpler to reason about than any of them, if the scoping story is right.

## What "the scoping story" even means

Where do arena boundaries come from?

Option A: lexical blocks are arenas.

```
arena {
    let buf = alloc_thing()
    let list = build_list(buf)
    use(list)
} // everything allocated in this block dies here
```

Nice because it mirrors control flow the reader already sees.
Ugly because "everything allocated in this block" is a strong claim, if arenas nest you need to know which one an allocation lands in, and if a function called inside the block secretly allocates into the caller's arena that's spooky action at a distance.

Option B: arenas are explicit values, threaded like a capability/effect.

```
fn build_list(a: Arena, buf) -> List {
    a.alloc(Node { ... })
}
```

This is the Zig/Odin move (allocator-passing).
Every allocating function takes an arena parameter (or reads one from context).
Verbose, but the allocation site is never a mystery, and you can pass a different arena (or a longer-lived one) explicitly when you need to.
Feels almost like Reader-monad-for-memory.

Option C: arena-per-call-frame, invisible, freed on return, like Jai's "temporary storage" default.
Fast to write in, but silently wrong the moment you return a pointer into that frame's arena and the caller uses it after the frame is gone.
That's just use-after-free with a friendlier name unless something (the type system? a lint? a runtime tag?) catches it.

## The real problem: escaping references

Every version of this note runs into the same wall.
Arenas make freeing cheap but they don't make "does this reference outlive its arena" go away.
That question is exactly what a borrow checker exists to answer.
So either:

1. The language accepts arenas as a convention/discipline, not a checked guarantee (Zig's stance: allocator passed explicitly, use-after-free is a "you" problem, tooling like sanitizers/valgrind catches it later).
1. The language adds enough type-level tracking (lifetimes, regions-as-types a la Cyclone/Rust) to reject dangling-out-of-arena at compile time, which drags in most of the complexity ownership types were supposed to spare us.
1. Arenas are only ever used for structures provably local (never returned, never stored past the block), enforced by a much narrower and simpler check than full borrow checking: something like "no reference derived from this arena crosses this boundary," which is close to what escape analysis already does in GC'd languages, just inverted (escape analysis normally decides stack vs heap, here it'd decide arena-local vs must-be-longer-lived).

Option 3 feels like the sweet spot to chase, if it's tractable.
Escape analysis is well trodden.
The novelty would be surfacing its result as an error instead of silently promoting to a longer-lived allocation.

## What if arenas aren't the whole story

Nothing says memory strategy has to be one mechanism everywhere.
Sketch: arenas as the *default* and fast path for "this bag of stuff has one clear lifetime" (parsing a file, handling one request, one frame of a game loop), with a fallback mechanism (refcounting? a real GC? explicit heap alloc/free?) for the minority of values that are genuinely graph-shaped and don't fit a region.

That's basically what a lot of real systems do already (per-request arena in a web server, GC for everything else) but usually as a library pattern, not a language primitive.
Could the language make "this value's lifetime is one arena" the common/cheap case and "this value needs the general-purpose allocator" the opt-in/slower case, syntactically distinguished so the reader always knows which regime they're in?

## Arenas and \[[green-threads-threading-model]\]

If green threads/goroutines are in the language, arena-per-task is an obvious pairing: spin up a task, give it an arena, tear down both together.
Structured concurrency (task can't outlive its scope) and arena lifetime (allocations can't outlive their scope) are the same shape of constraint.
Wonder if one mechanism could enforce both: "this scope owns a task and an arena, neither escapes."

## Dead end: arenas as GC replacement, full stop

Tried assuming arenas totally replace GC with no fallback (pure option 1 or 3, no escape hatch).
Died fast: any language that wants persistent data structures, caches, long-lived registries, or anything shaped like a graph with unclear single-owner lifetime is going to fight this model constantly.
Arenas are great for "phase-based" memory (a compiler pass, a request, a frame) and bad for "who knows how long this lives" memory.
A real language probably needs both, which reopens the question this note was trying to dodge: still need *some* answer for the long-lived case, arenas only ever cover part of the map.

## Open tension

Left unresolved, on purpose: is the arena the primary allocation model with GC/refcounting as the rare fallback, or is GC the default with arenas as an opt-in fast path for the phase-shaped cases?
Those are very different languages to write code in, and this note doesn't know which one this language wants to be.
