# No null type, no null representation

Starting premise: not just "no null keyword," but no bit pattern anywhere in the runtime representation of a value that means "nothing," "unset," or "invalid," unless a type explicitly says so.
Tony Hoare's "billion dollar mistake" line is the easy pitch, but the interesting design work is in what replaces the pointer-shaped hole null used to fill.

## Two different claims, worth separating

(a) "The language has no `null` literal or `nil` keyword."
That's cheap.
Plenty of languages claim this and still have null in disguise: an uninitialized `int` defaulting to 0, an uninitialized pointer defaulting to all-zero-bits, a `String?` that's really a tagged union with a null-shaped tag.

(b) "There is no representable value that means absence unless the type says `Option<T>` or equivalent, and even then absence is a tagged variant, not a sentinel bit pattern."
This is the real claim.
It means no default-zero-value semantics for any type that doesn't itself have a sane zero (numbers do, most things don't), no uninitialized variable ever observable before assignment, and every "this field might not be there" case forced into the type.

This note is about (b).
(a) falls out of it for free.

## Candidate mechanisms

### 1. Option/Maybe as an ordinary sum type

```
type Option<T> = Some(T) | None

fn findUser(id: UserId) -> Option<User> {
    ...
}

match findUser(id) {
    Some(u) => greet(u),
    None => print("no such user"),
}
```

Pro: uniform, composes with generics, `Option<Option<T>>` is a real distinct type from `Option<T>` (unlike null, which flattens).
Con: allocation/representation question.
Is `Option<&T>` the same size as `&T`, using the fact that a valid reference is never the null bit pattern to encode `None` for free (Rust's niche optimization), or is it always a tag byte plus payload, uniform but bigger?
If the language wants "no null representation" as a hard rule, niche-packing `None` into the all-zero-bits pattern of a reference is philosophically fine (nothing is claiming that bit pattern means "a valid reference to address 0"), but it's worth being honest that it's the same trick, just contained inside the type instead of leaking out of it.

### 2. Non-nullable by default, `T?` sugar over Option

```
fn findUser(id: UserId) -> User? { ... }

let u = findUser(id)
if u != null {   // or `if let user = u`
    greet(u)
}
```

This is Kotlin/Swift's move: keep the word "null" and the `?`/`!` ergonomics, but make the compiler enforce that every `null`-shaped access is checked.
Tension with the premise: this note's whole point was banning the word and the representation, and `T?` reintroduces both at the surface syntax level, just compiler-checked.
Maybe that's fine, syntax sugar over `Option<T>` where `null` desugars to `None` and the checked-access desugars to `match`.
Or maybe keeping the word around is a trap: it invites people to think of it as "a null that's been made safe" rather than "a different type entirely," and that framing leaks into how people design APIs (reaching for `?` reflexively instead of asking whether absence is even a valid state).

### 3. No absence at all, callers must prove presence

Wilder idea: don't have `Option<T>` be a thing you pattern-match at the call site casually.
Make "this collection might be empty" or "this lookup might fail" something the type system forces you to handle at the boundary where the possibility was introduced, and forbid `Option<T>` from floating around deep in business logic.

```
fn findUser(id: UserId) -> User? { ... }   // boundary: DB lookup, might fail

fn greetUser(id: UserId) {
    let user = findUser(id) orElse { return }   // resolved immediately, User from here on
    doStuffWith(user)   // this function's body never sees Option<User>
}
```

This isn't really a different representation mechanism, more a style/lint layer on top of (1): "resolve your Options at the edges, don't let them infect deep call chains."
Interesting because it's the same shape of complaint as the `?`-propagation "eighteen frames of busywork" problem in \[[no-exceptions-explicit-errors]\] but inverted: there the complaint is too much ceremony propagating, here the risk is too little, `Option<T>` seeping everywhere and every function having to handle a None it has no opinion about.
Same underlying tension (explicit-everywhere vs explicit-at-the-boundary) showing up in both the error note and this one.
Suggests these two notes are actually one open question about the language: "how does explicit-by-type interact with propagation," not two.

### 4. Required fields, no partial construction

If there's no null, there's no such thing as a struct with an "unset" field mid-construction.
That means either:

- constructors must supply every field in one shot (no incremental builder pattern that leaves fields null in between), or
- a builder pattern exists but is typed so the builder's type changes as fields get filled in (typestate-ish), so "call `.build()` before all fields are set" is a compile error, not a null-field-at-runtime bug.

```
struct Config {
    host: String,
    port: Int,
}

// no `Config{}` with zero values, no `Config{host: null, port: null}` half-state
let cfg = Config { host: "localhost", port: 8080 }
```

Builder version gets gnarly fast if attempted with typestate.
Might be a note of its own someday (`docs/notes/typestate-builders.md`?), parking it here as a dead end for now: full typestate builders for every struct is a lot of machinery for a problem that "just require all fields at construction, and if you need a builder, write `Option<T>` fields on the builder struct explicitly" solves more simply.
The honest tension: sometimes you DO want to build a thing incrementally with real absence between steps.
That's not the same as null-the-bug, it's Option-the-feature, being used correctly.
So "no null" doesn't mean "no absence ever," it means "absence is always spelled out."

## The hard part: where does the language's own machinery need a "no value yet"

Every language implementation, even ones with zero null in the surface language, tends to need SOME internal notion of "not yet initialized" during compilation or at the memory layer.

- Arrays: `Array<T>(size: 10)` with no initializer, what's in the slots?
  Options seem to be: (i) force an initial value always (`Array<T>(size: 10, fill: defaultValue)`), (ii) arrays start empty and only grow via push (no fixed-size uninitialized allocation at all), (iii) allow uninitialized memory but make the type system prevent reading a slot before it's written (this is close to what Rust does with `MaybeUninit`, and it's real complexity).
  (ii) is the cleanest fit for "no null representation, full stop" but forecloses on a class of performance-sensitive code (pre-sized buffers you fill in a loop) unless there's an escape hatch.
- Recursive/self-referential structs: a `Tree` node with a `parent` pointer set after construction, during a tree-building pass, has to be "not yet pointing at anything" for a moment.
  `Option<&Node>` handles this, but now every read of `.parent` is a match, forever, even after the tree is fully built and the invariant "every non-root node has a parent" holds.
  Is there a way to have a type that's `Option<T>` during construction and `T` after, without two different types and a conversion step?
  This smells like the same "resolve Option at the boundary" idea from mechanism 3, applied to time (construction phase) instead of space (call boundary).
- FFI: C doesn't have this rule.
  A null pointer coming in from a C function is a fact of life at the boundary.
  Does every FFI-imported pointer type become `Option<Ptr<T>>` automatically at the binding layer, forcing a check before first use?
  That seems right and is probably non-negotiable if the language wants to claim "no null representation" as a real invariant rather than a fair-weather one that C interop quietly violates.

## Open tension: does "no null" survive contact with performance

The Rust niche-optimization trick (mechanism 1) is popular precisely because "no null representation, but still one machine word for `Option<&T>`" sounds like a free lunch.
It's not entirely free: it only works when the underlying type has an unused bit pattern to steal (a reference that's guaranteed never zero, an enum with fewer variants than its tag byte can represent).
For `Option<Int32>` there's no spare bit pattern (all 2^32 values are legal `Int32`s), so it has to be a real tag plus payload, doubling the size in the worst case (tag byte, padded to alignment, so `Option<Int32>` might be 8 bytes not 4).

If the language cares about this (does it, foundational decision not yet made), "no null representation" as stated might need a companion decision about whether the compiler is allowed/expected to niche-pack automatically, or whether that's left to a `#[repr]`-style opt-in, or whether it's a non-issue because the language isn't chasing C-level memory density in the first place.
Parking this because it's downstream of memory-strategy and compilation-model decisions that are still open.

## Dead end: "just use empty collections and sentinel values instead of Option"

Tempting shortcut: instead of `Option<User>`, return an empty `List<User>` (0 or 1 elements) for "maybe a user," or use a sentinel like `UserId(-1)` for "no id."
Both are null wearing a costume.
Empty-list-as-maybe loses the ability to distinguish "definitely zero" from "haven't computed yet" in some contexts, and sentinel values reintroduce exactly the bug class null causes (nothing stops a future caller from constructing `UserId(-1)` "for real" and having it silently mean something else).
Rejecting this outright, but worth having written down why, since it's the obvious naive alternative someone will suggest.

## Loose thread

If the language ends up going the effect-system route floated in \[[no-exceptions-explicit-errors]\], does "no null" get to fold into the same machinery as "no exceptions"?
Absence and failure aren't the same thing (a `None` isn't an error, a lookup miss can be a totally expected, non-exceptional outcome), but both are "this call didn't give you the straightforward thing you asked for," and some languages (Zig sort of, with error unions) mash them into one mechanism.
Not clear that's right here.
Keeping them conceptually separate for now: `Option<T>` for "might not exist," `Result<T, E>` for "might fail," and a function that can do both is `Result<Option<T>, E>` or `Option<Result<T, E>>` depending on which failure mode is "outer."
Which of those two nestings is more common in practice is itself worth watching for once real example programs get written.
