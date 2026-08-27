---
title: 'Decision: Concurrency'
stage: 1
status: active
created: 2026-08-27
updated: 2026-08-27
---

# Decision: Concurrency

## Motivation

The wishlist already commits to three properties of the concurrency story without naming a mechanism for any of them: a story where data races are hard or impossible to express, structured concurrency where a child task cannot outlive its parent scope, and concurrency without function coloring, any function can block, no async-marked type infects its callers.
`docs/notes/green-threads-threading-model.md` found that these three pull against each other more than casual "goroutine-like" language admits: the no-coloring property and Go's actual unstructured lifetime model are separable, and "goroutines" as commonly understood buys no-coloring while doing nothing for race-freedom or structured lifetimes.
This record exists to pick an actual mechanism and be honest about which of the three properties it wins outright, which it wins partially, and which it has to fight for.

This record inherits real constraints from two records already interviewed, even though neither has reached stage 3.
The memory-strategy record's interview settled on two regimes, a default ownership-and-borrowing regime with no runtime cost and a second, explicitly marked, compiler-inserted-reference-counting regime for shared or graph-shaped values, and both interviews (memory, compilation model) independently answered "no runtime at all."
Any concurrency mechanism that needs a scheduler, an actor mailbox system, or similar standing machinery has to reconcile with that answer or explain why concurrency gets to be the exception.

### Scenario: a function that blocks

A function reads from a socket.
Calling it from ordinary code should not require the caller, or the caller's caller, to be marked in any special way.
This is the no-coloring wishlist commitment made concrete: the difference between `resp := http.Get(url)` and `let resp = http_get(url).await`, where the second spelling infects every function on the call stack above it.

### Scenario: two tasks racing on the same variable

Two concurrently running tasks close over the same mutable variable, one reads it, one writes it, with no synchronization between them.
The language should make this either impossible to express, or require an explicit, visible admission, not something that compiles silently and fails nondeterministically at run time under load.

### Scenario: a function that spawns two children and returns

A function spawns two tasks to fetch two URLs concurrently, combines their results, and returns.
Neither task should still be running, or leaked, once the function has returned or panicked.
This is the opposite of Go's actual answer, `go f()` is fire-and-forget with no built-in join or cancellation propagation, `context.Context` and `sync.WaitGroup` are library-level patches for a gap the language left open, and "goroutine leak" is common enough vocabulary in Go to be its own bug category.

### Scenario: tens of thousands of concurrent in-flight tasks

A server holds one long-lived task per client connection, with tens or hundreds of thousands of clients connected at once.
OS-thread stack footprint and context-switch cost genuinely dominate at this scale, which is the specific, checkable reason green threads or an async multiplexer exist at all rather than "just use OS threads."
Whether this language's target workloads include this shape is itself unresolved, and the answer changes which candidates below are even in contention.

## Rough shape

Five candidate directions, sketched.
As with the other records, these are not necessarily exclusive in the limit, and the syntax below is illustrative only.

### A. OS threads plus an async I/O multiplexer, no green threads

Every task is a real kernel thread; blocking I/O is handled by a runtime-level multiplexer (epoll, io_uring) underneath ordinary-looking blocking calls, with no custom scheduler and no stack management.

```
spawn fetch(url)   // a real OS thread, or a thread-pool task, under the hood
```

Buys: no scheduler to build, no growable-stack machinery, no GC-stack-map coupling, and thread creation on modern Linux is far cheaper than the 2009-era priors that motivated goroutines in the first place.
Costs: loses outright at the fan-out scenario above; megabyte-class OS thread stacks and context-switch cost dominate at tens or hundreds of thousands of concurrent tasks, which is the one scenario this candidate cannot answer by construction rather than by tuning.

### B. Stackful M:N green threads, Go-shaped

A scheduler multiplexes many lightweight tasks, each a real growable stack, onto a small pool of OS threads.
Memory is shared by default; channels are provided as one synchronization primitive among others.

```
spawn fetch(url1)
spawn fetch(url2)
// unstructured: nothing joins these, nothing cancels them
```

Buys: the strongest no-coloring story of any candidate, a real stack means any function can block without a marker.
Costs: this is the "just copy goroutines" baseline the note names as a dead end: shared-by-default memory does nothing for the race-freedom scenario (`go vet -race` catches races at run time, not statically), the model is unstructured by design so the third scenario needs a library patch bolted on the way Go's own `context`/`WaitGroup` are, and the scheduler, growable stacks, preemption, and blocking-syscall handling are exactly the standing runtime machinery in tension with the "no runtime at all" answer both other records have already leaned toward.

### C. Stackless async/await, Rust/JS-shaped

Tasks are compiled to state machines rather than real stacks; an `async` marker on a function and an explicit await point at each suspension replace a scheduler with a library-level executor.

```
async fn fetch(url: &str) -> String {
    let resp = http_get(url).await
    resp.body
}
```

Buys: no custom scheduler or stack-management machinery at the language-runtime level, an ordinary library can serve as the executor.
Costs: dies directly on the no-coloring wishlist commitment already made; `async`/`.await` is exactly the marking this project has already said it does not want, and the color infects every caller up the chain by construction, not by an implementation accident that could be fixed later.

### D. Actor model, process isolation, Erlang-shaped

Each unit of concurrency owns its own heap; nothing is shared, and communication is exclusively by message, which copies.

```
actor Fetcher {
    on Fetch(url) -> reply(http_get(url))
}
```

Buys: races become inexpressible by construction, there is no shared mutable memory to race on, and "let it crash" supervision gives a real, load-bearing answer to failure isolation that the other candidates do not address at all.
Costs: copying overhead with no path to zero-copy sharing even when the programmer can prove it would be safe, and an actor's mailbox and scheduler are, again, standing runtime machinery.

### E. Shared memory with static aliasing control, Pony-shaped

Memory can be shared and mutated, but a second, compile-time-checked discipline, a reference capability or similar, tracks which references may cross a task boundary and rejects the ones that would let two tasks alias the same memory mutably.

```
fn spawn_with(iso data: Buffer, fn: (Buffer) -> Unit) { ... }
// `data` must be uniquely referenced (isolated) to cross into the new task;
// the compiler rejects a caller that keeps a second live reference to it
```

Buys: keeps zero-copy sharing while making races a compile-time error rather than a runtime possibility, no coloring at all since this is a constraint on aliasing, not a marker on functions that can suspend, and no actor mailbox or copying cost.
Costs: a second static discipline riding alongside whatever the paradigm record's type system already does, and it is close enough to the memory-strategy record's own regime system, ownership and borrowing as the default, that the two may be describing overlapping mechanisms from two different records rather than genuinely separate ones.

## Prior art

- **Pony: reference capabilities and deny properties.**
  Pony's type system is founded on deny capabilities, denying certain operations on a reference rather than granting them, which lets the compiler guarantee data-race freedom with no locks and no runtime check ([Pony papers](https://www.ponylang.io/learn/papers/); [Introduction to the Pony programming language](https://opensource.com/article/18/5/pony)).
  Direct precedent for candidate E, and the clearest existing evidence that "compile-time race freedom without a tracing collector" is buildable, at the cost of a capability system the programmer has to learn.

- **Project Loom: virtual threads and structured concurrency as separable JEPs.**
  Virtual threads reached production in JDK 21 (JEP 444) as JVM-managed lightweight threads with automatic yielding; structured concurrency shipped as a preview alongside it and has iterated across JEP 428, 453, 480, and 505 without being finalized as of this record.
  Relevant because the two shipped as genuinely separate proposals years apart, direct evidence for the green-threads note's own finding that stackful/no-coloring and structured lifetime are orthogonal properties, not one feature.

- **Erlang: the actor model and "let it crash."**
  Each Erlang process owns its own heap and communicates only by copying messages, so two processes cannot race on shared memory because there is no shared memory to race on ([The actor model in 10 minutes](https://www.brianstorti.com/the-actor-model/)).
  Supervision trees built on "let it crash" are Erlang's answer to failure isolation, a property none of this record's other candidates address as a first-class concern.

- **Nathaniel J. Smith, "Notes on structured concurrency, or: Go statement considered harmful."**
  The essay already cited in the wishlist's own structured-concurrency entry, and Trio's existence is the proof this record leans on: structured concurrency was bolted onto a stackless, colored model (Python's `async`/`await`), which is direct evidence that the structured-lifetime property does not require any particular answer to the stackful-vs-stackless axis.

- **Go's asynchronous preemption, `golang/go` issue 10958.**
  Go shipped for five years, 1.0 through 1.13, with purely cooperative scheduling, where a tight loop with no function calls could starve every other goroutine pinned to its OS thread, before Go 1.14 added signal-based preemption to fix it.
  Cited here as a concrete cost candidate B has to price in from the start rather than retrofit: a spec written today gets to decide this up front, but only if it is named as a real design question rather than assumed away by "goroutines are simple."

Where the search was thin: nothing found treats "which of these five models best coexists with an already-chosen two-regime memory strategy" as a solved problem elsewhere, most concurrency-model writing treats memory management as a fixed backdrop (Go's GC, Erlang's per-process heap, Pony's own reference-capability-based memory model) rather than a separately-decided, still-open axis the way this project's own process requires.
That gap is recorded as an open question below rather than folded into a candidate.

## Open questions

1. Is a runtime component, a scheduler, an actor mailbox system, or similar standing machinery, acceptable at all for concurrency specifically, given that the memory-strategy and compilation-model interviews have both already leaned toward "no runtime at all"?
   This bounds the candidate set before anything else here is answered: if the answer is no, candidates B, D, and most implementations of E are dead before their other tradeoffs matter, and only A survives cleanly.
1. Which of the three wishlist commitments, no coloring, race-freedom, structured lifetime, is the one allowed to be the awkward case, if the chosen candidate cannot cleanly win all three at once?
1. If candidate E is chosen, is its capability discipline the same mechanism as the memory-strategy record's two-regime system, an extension of the default ownership-and-borrowing regime across a task boundary, or a third, separate discipline layered on top of both existing regimes?
1. Is structured concurrency, a child task cannot outlive its parent scope, a keyword-level construct built into the language from day one, or a library convention that risks becoming Go's own `context.Context`-shaped patch for a gap the language left open?
1. Are mutexes and atomics sanctioned primitives alongside channels or message passing, or does the race-freedom wishlist commitment push toward minimizing, or forbidding, shared-memory synchronization primitives at the language level entirely?
1. Does the chosen candidate need to serve the high-fan-out, tens-of-thousands-of-concurrent-tasks workload shape, or is that explicitly out of scope for this language's target use cases, which changes whether candidate A's simplicity is a real contender or a non-starter?
1. **Blocked on an undecided foundational decision:** the memory-strategy record (`0001`) has an interviewed but not-yet-stage-2 answer (two regimes, no runtime), and this record's candidate E in particular cannot be fully specified until that record says how a value moves, or is checked to move, across a thread boundary.
1. **Blocked on an undecided foundational decision:** the paradigm record (`0003`) has an interviewed but not-yet-stage-2 answer (immutable by default), and immutability-by-default is exactly the property that makes shared-memory races cheaper to rule out for free; this record should not presume that answer is final, even though it would make several candidates here easier.
1. **Blocked on an undecided foundational decision:** the compilation-model record (`0002`) has an interviewed but not-yet-stage-2 answer of "no runtime at all," which open question 1 above depends on directly; that record's own open question 5 asks the identical question independently, and the three records currently do not share one answer.

## Advancement record

- 2026-08-27, gate 0 → 1: sketched from the wishlist's existing concurrency entries (data races hard to express, structured concurrency, no function coloring) and `docs/notes/green-threads-threading-model.md`; five candidate directions surveyed, prior art cited across Pony's reference capabilities, Project Loom's separable virtual-threads and structured-concurrency JEPs, Erlang's actor model, structured concurrency's own founding essay, and Go's asynchronous-preemption history; open questions recorded including a load-bearing block shared with both other interviewed-but-not-yet-stage-2 records on whether a runtime is acceptable at all.

## Changelog
