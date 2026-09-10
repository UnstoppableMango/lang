# Green threads / threading model

Riffing on the wishlist's concurrency section, specifically the two entries: "data races hard or impossible to express" and "structured concurrency, where a child task cannot outlive its parent scope."
Neither entry commits to a threading model.
This note is about M:N green threads specifically: what buying that model actually costs, and whether it fights the two entries above.

## What "green threads" means, precisely

Three shapes worth keeping distinct, because "lightweight concurrency" gets used for all three:

- **1:1 (OS threads).** Every task is a real kernel thread. Simple runtime, expensive tasks (megabyte stacks, syscall-cost context switches).
- **M:N (green threads, Go/Erlang/Loom).** A scheduler multiplexes many lightweight tasks onto a small pool of OS threads. Needs a scheduler, needs stack management, needs a story for blocking syscalls.
- **1:N stackless (async/await, Rust/JS/Python asyncio).** Tasks are state machines, not real stacks. No custom scheduler needed at the language-runtime level (an executor library suffices), but functions get colored: async fn can only be awaited from another async fn, and the color infects every caller up the chain.

The wishlist doesn't name a shape.
"Goroutine-like" language is common in casual conversation about this project but that's prior art leaking in, not a decision.

## The actual selling point: no coloring

The reason people reach for "goroutines" as shorthand for "the good version of this" is specifically the no-coloring property.
A goroutine is a real stack.
Any function can block on a channel, a mutex, an I/O call, without being marked specially, and its caller doesn't need to be marked either.
Compare:

```
// stackful, no coloring (Go-shaped)
func fetch(url string) string {
    resp := http.Get(url)   // blocks this goroutine, not the OS thread
    return resp.Body
}

// stackless, colored (Rust/JS-shaped)
async fn fetch(url: &str) -> String {
    let resp = http_get(url).await;   // "async" and ".await" both required,
    resp.body                          // and now every caller must be async too
}
```

Erlang processes and Java's virtual threads (Loom) both bought the same property for the same reason: blocking I/O ergonomics without a second, parallel type system for "things that can suspend."
That's a real, well-documented win.
It's also not free, which is most of this note.

## What buying stackful M:N actually costs

Things that "just spawn a goroutine" hides, that a spec would have to actually specify:

- **A scheduler.** Not a library concern, a runtime concern: who decides which green thread runs on which OS thread, and when.
- **Growable stacks.** Green threads start with a small stack (Go starts at 2KB) and must grow it as call depth increases. Go's history here is a real cautionary tale: segmented stacks (cheap to grow, but broke C interop assumptions about stack contiguity, "hot split" pathological case where a loop straddled a segment boundary) replaced by copying contiguous stacks (requires rewriting every pointer into the stack on move, which requires the language to know precisely where every stack-relative pointer lives, i.e. requires precise GC-style stack maps).
- **Preemption.** Cooperative-only scheduling (yield at function calls, channel ops, allocations) means a tight numeric loop with no calls in it can starve every other green thread pinned to that OS thread. Go shipped without a fix for this until 1.14, which added signal-based async preemption. A spec written today gets to decide this up front instead of retrofitting it years in.
- **The FFI/blocking-syscall boundary.** A green thread that calls into C, or makes a blocking syscall the runtime doesn't know how to make async, blocks the OS thread carrying it. Go's answer is `sysmon`, a background thread that detects a goroutine has been in a syscall too long and spins up a replacement OS thread so the rest of that thread's runqueue isn't stuck. That's nontrivial standing machinery, not a footnote.
- **GC interaction.** Root scanning has to walk every live green-thread stack, and stacks are moving/growing, so precise stack maps and the GC are coupled. This connects hard to the still-open memory strategy decision record: stackful green threads are a well-trodden path under a tracing GC (Go), and a much less trodden one under region/arena or ownership-based memory. Not clear this has been done cleanly anywhere.

None of this is disqualifying.
It's the price tag for "no coloring," and it should be priced in on purpose rather than discovered mid-implementation.

## Where this fights the wishlist, not just costs something

The "data races hard or impossible to express" entry is the sharper tension.
Goroutines share memory by default.
Channels are provided, but nothing stops two goroutines from closing over the same mutable variable and racing on it; `go vet -race` catches this at runtime-ish (instrumented builds), not statically.
So "goroutine-like" as a phrase smuggles in "races are easy to write, we just also give you channels," which is close to the opposite of the wishlist entry.

Two ways out, neither obviously right:

1. **Erlang-style isolation.** Each process has its own heap; nothing is shared, message passing copies. Races become inexpressible by construction (there's no shared mutable memory to race on) at the cost of copying overhead and no path to shared-memory speed even when the programmer knows it'd be safe.
1. **Static aliasing control (Rust-shaped).** Keep shared memory, but the type system tracks which references can cross a thread boundary and rejects the ones that would alias mutably. Kills the coloring win's sibling problem (now you have a second static discipline, borrow-shaped, riding alongside whatever the base type system already does) but keeps zero-copy sharing.

Both are legitimate answers.
Both are also *not* "copy goroutines," which is worth naming plainly since "goroutine-like" is the phrase that keeps coming up in casual description of this corner of the design.

## Doodle: structured spawn as the actual point of divergence from Go

The second wishlist entry, "a child task cannot outlive its parent scope," is something goroutines explicitly do not give you.
`go f()` is fire-and-forget.
There's no built-in join, no built-in cancellation propagation; `context.Context` and `sync.WaitGroup` are library-level patches for a gap the language left open, and "goroutine leak" is common enough vocabulary in the Go world to be its own category of bug.

A structured version might look like:

```
scope {
    spawn fetch(url1)
    spawn fetch(url2)
    // scope block does not exit until both spawned tasks finish,
    // and a panic in either cancels the scope's siblings
}
```

versus Go's actual answer, which is closer to:

```go
go fetch(url1)  // nothing tracks this. if fetch panics, nobody's watching.
go fetch(url2)  // if the caller returns first, these are still running.
```

The stackful/no-coloring property and the structured-lifetime property are orthogonal.
Trio (Python) proved you can bolt structured concurrency onto a stackless/colored model.
Nothing stops pairing structured lifetimes with stackful green threads either, it's just that Go, the most famous stackful example, happens not to have done it, so "goroutine-like" as shorthand quietly drops this half of the wishlist unless said otherwise.

## An alternative that skips green threads entirely

Worth naming since it dies fast but for a specific, checkable reason: OS threads plus a good async I/O multiplexer (io_uring, epoll under a runtime), no custom scheduler, no stack management, no GC-stack-map coupling.
Thread creation on modern Linux is cheaper than the 2009-era priors that motivated goroutines in the first place.
This is a real position some newer systems languages take.

Where it actually loses: very high fan-out I/O, tens or hundreds of thousands of concurrent in-flight tasks (a connection-per-client server being the canonical case), where OS thread stack footprint (megabyte-class, even with `ulimit` tuning) and context-switch cost genuinely dominate.
Green threads exist because that scenario is real, not because someone wanted a scheduler for its own sake.
So this isn't a dead end so much as a reminder that "just use OS threads" is a legitimate default that only loses on a specific, nameable workload shape, and whether this language's target workloads include that shape is itself an open question.

## Dead ends, recorded so I stop rediscovering them

- **"Just copy goroutines."** Dies on the structured-concurrency wishlist entry specifically: goroutines are unstructured by design, and the fix (WaitGroup/Context) is a library patch for a language-level gap, not something worth re-importing on purpose.
- **"Goroutines but make channels the only sync primitive, ban shared mutable state at the language level."** Attractive, but "ban shared mutable state" is doing all the work that entry 2 (static aliasing control) or entry 1 (process isolation) would have to do properly anyway. Saying it as a constraint doesn't design the mechanism that enforces it.

## Threads worth pulling later

- How this interacts with the still-open memory strategy decision record. Stackful green threads under a tracing GC is well-trodden (Go); under region/arena or ownership-based memory it's much less charted, and the growable-stack/precise-stack-map machinery above may not port cleanly.
- Whether cancellation is a language-level construct baked into `scope`/`spawn` from the start, or a library convention bolted on later the way Go's `context.Context` was. The wishlist's structured-concurrency entry argues for baking it in, but that's a note observation, not a decision.
- Whether mutexes and atomics exist alongside channels, or channels are the only sanctioned primitive, given the races-hard-to-express entry pushes toward minimizing the shared-memory sync surface generally.
- Whether "goroutine-like" should stop being used as shorthand in future notes and design docs, since this note's conclusion is that the phrase actually means two separable things (no coloring, and Go's specific unstructured lifetime model) and only one of them is clearly wanted.
