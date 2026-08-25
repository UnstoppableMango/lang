compiled language first, with a scripting-like alternative

the pitch: design and ship the compiled language as the "real" language.
get a scripting-shaped mode later, maybe much later, as a second front end or a second mode of the same front end.
not a separate language, a separate posture toward the same core.

## why compiled first

forces the type system and semantics to be decided honestly instead of deferred to "the interpreter will figure it out at runtime."
a compiled language that can't explain its own errors at compile time is just a slow scripting language with extra steps.
scripting languages get to punt on a lot (arity checking, type coercion, name resolution) until the line actually executes.
if compiled is the foundation, punting isn't available, which is a forcing function for rigor.

counterpoint: this also means the fun, malleable, "just try it" experience is not available on day one.
the wishlist workflow already prizes going slow and not presuming outcomes, so maybe that's fine, maybe even the point.
but there's a real cost: nobody plays with a compiler.
people play with a REPL.
if the only way to touch the language for the first year is `mango build && ./a.out`, the feedback loop that makes language design fun (and honest) is gone.

## what "scripting-like alternative" could even mean

candidate shapes, from least to most divergent:

1. **same language, interpreted front end.**
   identical syntax and semantics, but `mango run file.mn` skips the object file and tree-walks or bytecode-interprets.
   scripting-like only in invocation, not in language design.
   this is the "cargo script" / "go run" move: same type checker, same everything, just skip the ceremony of a build step.

2. **same language, relaxed mode.**
   `mango run --loose` (or a `#!/usr/bin/env mango --script` shebang) turns off some compiled-language rigor: top-level statements allowed outside a function, implicit `main`, maybe looser numeric coercion.
   danger: now there are two dialects of the same language and every design decision has to be made twice ("does this feature exist in loose mode?").
   this is how you get Python's `from __future__ import` accretion, or TypeScript's `any`.

3. **genuinely different surface syntax, shared IR/semantics.**
   two front ends compiling to the same core language, the way ReScript and OCaml share semantics but not surface syntax, or the way ClojureScript and Clojure share semantics but not host.
   the scripting front end could look nothing like the compiled one: dynamic-looking, no explicit types, semicolon-optional, whatever "feels like a script" means.
   underneath, it desugars to the same typed core, maybe with holes filled by inference or by a lightweight runtime type.
   this preserves "one language, one semantics" while giving very different day-to-day ergonomics.

4. **actually a different language that happens to interop.**
   at this point it's not "the language, but scripty" — it's two languages, like Java and Groovy, or Kotlin and Kotlin Script.
   probably out of scope for a single-language wishlist entry, but worth naming as the far end of the spectrum so we know when we've crossed it.

the wishlist framing ("with a scripting-like alternative") sounds closest to (1) or (3): same language, different entry point or different surface, not a whole second language.

## what actually makes something feel "scripting-like"

worth separating the bundle of things people mean by "scripting language," because they're not one thing:

- **no build step.** `mango run` instead of `mango build && ./out`. this one is nearly free if the compiler can also interpret its own IR, or if it just compiles to a temp file and execs it transparently. arguably compiled languages should just do this by default (see: `go run`).
- **no explicit types.** inference does most of this for free even in a compiled language. the "scripting feel" people actually want here is often just "I didn't have to write the type," which type inference already buys you without touching semantics.
- **top-level imperative code, no `main` ceremony.** this is a syntax/whitespace decision, not a compilation-model decision. could allow bare top-level statements always, compiled or not, and just desugar to an implicit `main`. doesn't need runtime interpretation to get.
- **REPL / interactive eval.** this one genuinely wants an interpreter, or at least incremental compilation with a persistent environment. can't fake this with "compile to a temp binary each time" past a certain latency.
- **loose error handling / exceptions instead of explicit results.** this is a language-semantics decision (see other notes on error handling, not written yet) and probably shouldn't be tied to compiled-vs-scripting at all. tying "scripting mode" to "different error handling" is how you end up with mode (2) above and its two-dialects problem.
- **dynamic typing / duck typing.** the most divergent one. if scripting mode means actual dynamic typing, that's not a mode, that's a different type system, which pushes hard toward option (3) or (4).

so: a lot of "scripting feel" is achievable via tooling (`mango run`) and syntax sugar (optional `main`) without touching the type system or semantics at all.
the genuinely hard part is only the dynamic-typing slice, if that's even wanted.
maybe the honest wishlist entry is narrower than "scripting alternative": it's "no-build-step invocation" plus "optional top-level code," and the rest is scope creep from the word "scripting."

## an example, doodled

hypothetical compiled surface:

```
fn main() -> Unit {
    let name = read_line()
    print("hello, {name}")
}
```

hypothetical "scripting-like" invocation of the *same* program, no `fn main` required, top-level allowed:

```
let name = read_line()
print("hello, {name}")
```

same types, same stdlib, same errors, just skips the wrapper.
this is basically option (1) + top-level sugar from option (2)'s good half, without touching type rigor.
feels achievable without forking semantics.

now the more divergent doodle, something that reads more like a shell script or Python, per option (3):

```
name = input()
echo "hello, $name"
```

string interpolation via `$name` instead of `{name}`, `echo` instead of `print`, no `let`.
this is a genuinely different surface grammar.
underneath it'd desugar to the exact same typed core as the first example, with `name` inferred as `Str`.
this is more fun to imagine but a lot more work to build and maintain in parallel with the "real" syntax, and now every syntax decision in the main language needs a "how does this look in scripty-mode" answer too.

## dead end considered

thought about "scripting-like" meaning untyped/dynamic by default, with types as an opt-in annotation layer (gradual typing, like TypeScript's relationship to JS but inverted: compiled+typed is the base, dynamic is the escape hatch).
died fast: gradual typing systems are a huge design surface on their own (see TypeScript's `any`, soundness holes, `strict` mode fragmentation) and dragging that into "just want a scripting feel" conflates two big open decisions (type system shape, and compilation model) into one wishlist bullet.
if gradual typing is wanted, it deserves its own wishlist entry and its own note, not a rider on this one.

## open tension, not resolved here

is "scripting-like alternative" actually about a *second front end/dialect*, or is it actually just "the compiled language needs to not suck to iterate on" (fast `run` command, decent REPL, no `main` ceremony for scratch files)?
those are very different amounts of work and very different commitments.
the second one might fully satisfy the itch without ever forking the language.
worth sitting with before this graduates past a note.
