# C-compatible function interface

Wishlist already has "a C-compatible FFI as the lowest common denominator."
This note picks at what "C-compatible" actually commits us to, since the phrase hides a lot of decisions.

## What "C-compatible" could mean, ranked by how much it costs

1. Callable from C, with a C header generated for our functions.
   This is the cheap version: our compiler emits `.h` files, our ABI matches the platform C ABI for exported symbols only.
   Internal calling convention can be anything.

1. Callable from C AND able to call C, bidirectionally, but only through explicit `extern "C"`-style boundary blocks.
   Everything inside those blocks pays the FFI tax (marshaling, no fancy types); everything outside is free to be weird.
   This is Rust's approach, essentially.

1. C ABI as the native calling convention for all functions.
   Zig leans this way for extern-facing things but not necessarily for internal calls.
   Nobody serious makes their WHOLE calling convention C's, because C's ABI is bad (no multiple return values, no fat pointers, poor struct-passing rules that vary by platform).

Option 2 is almost certainly right.
"C-compatible" should mean "has a boundary that speaks C," not "is secretly C inside."

## The parts of C ABI compatibility that are actually hard

Calling convention (which registers, stack alignment, who cleans up the stack) is the part everyone thinks of first, and it's also the easy part: it's just following the platform's psABI document, already solved by LLVM/cranelift/whatever backend we pick.
The hard parts are all about what crosses the boundary conceptually, not mechanically.

### Strings

C strings are null-terminated bytes with no length.
Our strings (if we have a real string type) are probably length-prefixed, maybe UTF-8-validated, maybe rope-like internally.
Crossing the boundary means either:

- exposing `(ptr, len)` pairs and asking C callers to use a companion length param (fine for our own headers, useless for THEIR headers)
- eating the null-terminator scan cost and validation gap when receiving a `char*`
- never passing strings directly, only opaque handles with getter functions

Feels like strings are the single biggest tell for whether an FFI story is real or aspirational.
Every language's C interop docs spend the most words here.

### Ownership and memory across the boundary

Who frees a pointer that crosses the FFI line?
If our memory strategy ends up being GC'd (still an open decision record), then any pointer we hand to C must either:

- be pinned/never-move for the lifetime C might hold it (GC languages hate this)
- be copied into a C-owned allocation on the way out, freed by an explicit `free` function we also export
- be borrowed only for the duration of the call, never stored (this is the sane default, but C code doesn't always respect "don't store this")

This is actually downstream of \[[arena-memory-model]\] and the not-yet-written memory strategy decision.
Can't fully design the FFI story until that lands, but the note is worth keeping because the FFI requirement should probably be a design CONSTRAINT feeding the memory decision, not an afterthought bolted on after.
Arenas make the "who frees this" question easier in one direction: hand C a pointer into an arena, and as long as the arena outlives the call, no ownership transfer needed at all, since nobody frees anything until the arena dies.

### Struct layout

`repr(C)` in Rust exists because Rust's default struct layout is unspecified (field reordering for packing).
We'll need the same escape hatch if our own default layout is ever allowed to differ from C's, which it probably should be (we might want niche-filling enums, tagged unions, etc. that have no honest C representation).

So: two struct-layout modes.
Default layout, whatever's most efficient/idiomatic for us.
`extern` (or whatever the keyword is) layout, C-compatible, usable across the FFI boundary, with restrictions on what field types are even legal there (no tagged unions, no fat pointers, no GC references probably).

### Enums and sum types

C enums are just ints.
Our sum types, if we have pattern-matching-with-exhaustiveness (also on the wishlist), are richer than that.
A sum type with payload data has no direct C representation; it has to get flattened to a tagged struct manually, which means the FFI boundary REJECTS certain types outright rather than silently doing something lossy.
That's probably correct: an explicit compile error at the boundary ("this type is not FFI-safe") beats a silent lossy conversion.

### Error handling across the boundary

If \[[no-exceptions-explicit-errors]\] settles on `Result<T, E>`-shaped errors, that's another non-C-representable type at the boundary.
C convention is sentinel return values plus `errno`-style out-of-band state, or a `(bool ok, T value)` out-param pair.
Feels like the FFI layer needs its own error convention, distinct from whatever the "native" error story is: something like every `extern "C"` function returns an int status code and takes an out-pointer for the real return value.
That's what gRPC/protobuf-adjacent C APIs do, what SQLite does, what basically every serious C library does.

### Function pointers / callbacks

C callbacks are just pointers, no captured environment.
If we have closures, a closure crossing the boundary needs either:

- to be restricted to non-capturing functions only when going to C (compile error otherwise)
- to be represented as a `(fn_ptr, void* userdata)` pair, which is the universal C convention for "callback with context" (see: `qsort_r`, most C callback APIs)

The second option is more useful and is the actual standard pattern, so probably that.

## A tangent: do we even want raw C ABI, or a shimmed one?

What if instead of exposing our real functions to C, the "FFI" is always through a codegen'd C shim layer, like SWIG/cxx/uniffi do for other language pairs?
Then the compiler never has to make our REAL calling convention C-compatible at all; it just has to be able to EMIT C-compatible glue functions on request.
This decouples "what's efficient for us internally" from "what C sees" completely, at the cost of an extra hop (probably inlined away, but a hop conceptually).

This is basically option 2 from the top of the note, just phrased as implementation strategy instead of language design.
Kept as a reminder that the shim can live entirely in codegen, no new syntax needed beyond marking which functions are exported.

## Open tension, deliberately unresolved

Whether the C ABI boundary is a language-level concept (an `extern` keyword, FFI-safe type checking, the works) or a tooling-level concept (a `lang-cbind` codegen tool that reads our public API and emits a shim, no compiler involvement) is a real fork.
Rust picked language-level.
Most scripting languages with C extension stories (Python's ctypes/cffi, Lua's C API) picked tooling-level, or a hybrid.
Given the "no compiler yet" state of this whole project, tooling-level is tempting: it can be prototyped without touching the type checker at all.
But it means FFI-safety checking (the struct layout / sum type rejection stuff above) either doesn't exist or has to be reimplemented outside the compiler, which is worse.

No conclusion here.
Just noting that this fork exists and that the two paths diverge early.
