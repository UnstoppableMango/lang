# No exceptions or panics, explicit error handling

Starting premise: no `throw`, no `try/catch`, no unwinding control flow that jumps across stack frames invisibly.
Every function that can fail says so in its signature, and every caller has to look the failure in the eye.

## Why ban them at all

The pitch for banning exceptions is usually "every call site is a potential exit point, and you can't see it."
A function typed `(User) -> Profile` might actually be `(User) -> Profile | throws NetworkError | throws ParseError | throws OutOfMemory`, and nothing in the signature tells you.
Explicit error handling makes the type honest.

Counter-argument: exceptions exist because most errors are not handled at the call site anyway.
They're handled three frames up, at a boundary (top of a request handler, top of a task).
Forcing every intermediate frame to thread an error value through is busywork if that frame has no opinion about the error and is just passing it along.
This is the classic Go complaint: `if err != nil { return err }` eighteen times in a row.

So "no exceptions" alone isn't a design, it's a constraint that needs an answer to "then how does the 90% case (just propagate) not become ceremony."

## Candidate mechanisms

### 1. Result type + `?` propagation operator (Rust-style)

```
fn readConfig(path: Path) -> Result<Config, IoError> {
    let text = readFile(path)?      // propagates IoError automatically
    parse(text)
}
```

Pro: the common "just bubble it up" case is one character, not eighteen lines.
Con: the `?` is itself a tiny bit of magic, a control-flow jump disguised as an operator.
Is that a contradiction of "no invisible control flow"? Arguably not, because the jump is still visible in the type (`Result` in the return type) even if the operator elides the branch at each call site.
Feels like the right kind of invisible: the reader knows failure is possible from the signature, they just don't have to see every relay station.

### 2. Multiple return values, error last (Go-style)

```
config, err := readConfig(path)
if err != nil {
    return err
}
```

Pro: dead simple, no new type, no operator.
Con: `err` is a plain nilable value, not part of a sum type, so nothing stops you from ignoring it (`config, _ := readConfig(path)`).
Explicit error handling that can be silently discarded isn't really explicit, it's just optional and visible.
If we go this route we'd want the type system to make discarding an error require an explicit gesture (a lint-level "unused fallible value" or a `must_use`-style annotation).

### 3. Effect-style, errors as an algebraic effect

Wilder idea: instead of `Result<T, E>` infecting every signature, errors are an effect that a function can perform, and the effect system tracks it the way async/await tracks async-ness.

```
fn readConfig(path: Path) -> Config raises IoError {
    ...
}

fn main() {
    handle IoError as e {
        log(e)
        exit(1)
    } in {
        let cfg = readConfig("app.toml")
        ...
    }
}
```

This gets you the "?" ergonomics without a wrapper type showing up in every intermediate signature, at the cost of needing a whole effect system, which is a much bigger foundational bet than "how do errors work."
Feels like it belongs in a "what if the language has algebraic effects" note more than this one, but worth remembering the two ideas are coupled: effects make explicit-errors-without-Result-boilerplate much more natural.

### 4. Checked exceptions, but honest about it

Java tried "exceptions, but the signature has to declare them" (`throws IOException`).
It's arguably the closest existing thing to "explicit error handling without a Result type."
It died in practice because:
- checked exceptions don't compose through higher-order functions (what's the throws clause of `map`?)
- people papered over it with `throws Exception` or wrapping in unchecked exceptions anyway

Worth noting as a cautionary tale rather than a candidate: the failure mode wasn't "explicit is bad," it was "explicit but the mechanism didn't compose."
Any design here needs to answer the higher-order-function question up front: what is the type of `map: (A -> Result<B, E>) -> List<A> -> Result<List<B>, E>` and does it fall out naturally or need special-casing.

## Panics: separate question from errors

"No panics" is a different claim than "no exceptions."
Panics are usually reserved for programmer errors (index out of bounds, assertion failure, unreachable code reached), not for expected failure modes (file not found, network timeout).
Rust has both: `Result` for expected failure, `panic!` for "this should be impossible, and if it happens, something is corrupted, stop the world."

Question worth sitting with: does "no panics" mean:

(a) no unrecoverable process-aborting failure mode exists at all, full stop, or
(b) no *implicit* panics (array index out of bounds must be a `Result` or a compile error, not a runtime abort), but an explicit `fatal("unreachable")` escape hatch still exists for the programmer to hit intentionally

(a) is a much bigger claim.
It means array indexing returns `Option<T>` or requires a proof the index is in bounds, integer division by zero is a typed error not a trap, out-of-memory is a `Result` somehow (how do you even allocate the error value if you're out of memory?).
That last one is a real wrinkle: some failures are fundamentally hard to make "just another Result" because handling the error requires resources that might not exist.
This might be a case where "no panics" has to bottom out in "no panics for anything except a short fixed list of truly unrecoverable host failures," which is (b) with a narrower list than Rust's.

## Sketch: what "no exceptions or panics" programs feel like to write

Trying a tiny hypothetical program under mechanism (1), Result + `?`:

```
fn main() -> Result<(), AppError> {
    let cfg = loadConfig("app.toml")?
    let conn = connect(cfg.dbUrl)?
    let rows = conn.query("select * from users")?
    for row in rows {
        print(row)
    }
    Ok(())
}
```

Reads fine.
Now the same thing where three different error types (`IoError`, `DbError`, `ParseError`) need to unify into `AppError` at the `?` sites.
That's the part every language in this space struggles with: automatic error-type conversion (`From`/`Into` in Rust) versus manual wrapping at every boundary.
If conversion is automatic, `?` can silently launder a `DbError` into a vague `AppError`, which reintroduces some of the "can't tell what actually failed" problem exceptions have.
If conversion is manual, you're back to boilerplate at every propagation site, undercutting the whole point of `?`.

Dead end, or at least an open wound: no answer here yet, just noting that "explicit error handling" doesn't fully specify whether error *types* are explicit at every hop or only error *presence* is.
Those are different levels of explicitness and this note conflates them until now.

## Loose thread to maybe revisit

If the language ends up lazily evaluated (open decision elsewhere), does "no exceptions" even mean the same thing?
A lazy `Result` might not surface its failure until forced, at which point the "call site" that observes the error could be far from the call site that produced it, which smells like the exact invisible-control-flow problem this whole note is trying to avoid.
Not resolving that here, just flagging the two decisions aren't independent.
