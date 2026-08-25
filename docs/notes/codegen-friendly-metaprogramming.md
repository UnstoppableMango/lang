# Codegen-friendly metaprogramming

"Codegen friendly metaprogramming" is ambiguous on first read, and the ambiguity itself is worth sitting in for a second.
It could mean two different things.

1. The language has metaprogramming facilities (macros, derive, comptime), and those facilities are friendly to tooling: debuggable, inspectable, not a black box.
2. The language is a friendly *target* for external codegen, the way protoc, sqlc, or buf generate source files in some language and that language either cooperates or fights back.

These pull in different directions enough that they deserve separate threads.
Both threads below, no resolution on which one the wishlist entry actually meant.

## Thread 1: metaprogramming that doesn't hide from tools

The failure mode to design away from: a debugger stepping into macro-expanded code and showing garbage, an IDE's go-to-definition landing inside a `quote!` block, an error message pointing at the macro invocation site instead of the line that actually broke.
Rust proc-macros are the canonical bad example here, `cargo expand` exists as a third-party tool precisely because the compiler doesn't make expansion visible on its own.
That's a smell: if introspecting your own macro expansion needs a community-maintained side tool, expansion should have been a first-class compiler command from day one.

Three shapes to compare:

**Text substitution (C preprocessor style).**
Dead on arrival, not worth doodling syntax for.
No hygiene, no scoping, breaks every tool that isn't literally the preprocessor.
Noted here only so it's on record as ruled out early.

**Hygienic AST macros (Lisp, Scheme `syntax-rules`, Rust `macro_rules!`/proc-macro).**
Code as data sidesteps some of the text-substitution problems, no serialization round-trip since the macro operates on the real AST.
But hygiene is exactly the mechanism that makes expansion opaque to a human or a debugger, the whole point is that expanded identifiers don't collide with or shadow anything visible, which is also what makes them hard to *look at*.
Powerful, but the tooling cost is real and historically under-paid-for (see Rust's ecosystem-bolted-on `cargo expand`).

**Comptime-as-plain-execution (Zig style).**
```
fn derive_eq(comptime T: type) type {
    return struct {
        fn eq(a: T, b: T) bool {
            // reflect over T's fields, compare each
        }
    };
}
```
No macro language, no quoting, no hygiene problem to solve because there's no text splicing at all, just a normal function that happens to run at compile time and return a type.
Debug info is comparatively sane because it's literally normal code execution, not a separate expansion pass.
This reframes "metaprogramming" as "programming, at compile time, with types as values" rather than as a distinct macro sublanguage.

Tentative read: option 3 is the most codegen-tool-friendly of the three, precisely because it refuses to be a separate thing.
Fewer moving parts for a debugger, an LSP, or a formatter to special-case.

## Thread 2: being a friendly target for external generators

Separate question: if `protoc` or a schema compiler is going to vomit source files in this language, what makes that pleasant versus painful?

- Insignificant whitespace, or at least whitespace that's easy to emit correctly without a generator tracking indent state (see [[insignificant-whitespace]]).
  Python codegen tools spend real effort just getting indentation right; that's accidental complexity a generator shouldn't have to carry.
- A canonical formatter as an escape hatch (see [[single-canonical-formatter]]).
  Generator emits ugly-but-correct output, `lang fmt` cleans it, output is diffable and reviewable in a PR without the generator itself needing a pretty-printer.
- A `// Code generated, DO NOT EDIT` convention (Go's), cheap and it works.
- Source-mapping from generated output back to whatever schema or template produced it, so a compile error in generated code points somewhere a human can act on, not at a `.lang` file nobody is supposed to hand-edit.
- A grammar simple enough that "emit valid syntax" doesn't require a generator to understand precedence or associativity, over-parenthesizing should always be a safe fallback for a naive emitter.

This thread barely touches "metaprogramming" as a language feature at all.
It's closer to: don't make external tools reverse-engineer the formatter's opinions.

## Where the two threads collide: compile-time execution and reproducibility

Comptime-as-plain-execution (thread 1's tentative favorite) means arbitrary user code runs during compilation.
That's a build input, and an unpinned or unsandboxed one is exactly the kind of thing that breaks reproducible builds (see [[nix-first-build-system]]).
Rust proc-macros have the same property and mostly get away with it because the ecosystem doesn't prioritize bit-for-bit reproducibility as hard as this project apparently wants to.
Open question, not resolved: does comptime execution need to be sandboxed, deterministic-by-construction, or otherwise fenced off to keep the reproducibility story intact, or is "compilation is a build step and build steps are already trusted" enough of an answer.

## Dead end: no macros at all, only a closed set of compiler-builtin derives

Tempting for a minute: what if the friendliest metaprogramming is none, bake `derive(Eq)`-shaped things into the compiler as builtins, and never expose a generation surface to users at all.
Zero codegen surface, zero hygiene problem, zero tooling burden, because there's nothing for a debugger or formatter to special-case.
Died on: it's not actually metaprogramming, it's just a fixed list of compiler features wearing a metaprogramming costume, and the wishlist item presumably wants users to be able to write their own derive-shaped things, not petition the language author to add one.
Worth remembering as the limit case on one end of the spectrum, though, since "how close to this can we get while still being extensible" is a reasonable design pressure.

## Open tension

Do we need a macro system at all, separate from strong generics plus comptime reflection?
A lot of what macros are reached for (derive impls, boilerplate elimination, loop unrolling) is arguably just "generic code with enough type-level reflection," which Zig mostly proves you can get without a distinct macro sublanguage.
Rust needed `macro_rules!`/proc-macros harder than Zig needs macros partly because Rust's const-generics and reflection story is weaker.
Not resolved here: whether this language's answer is "strong comptime + reflection, no separate macro layer" or "macro layer, kept honest by mandatory first-class expansion tooling."
Also unresolved: which of the two threads above the original wishlist phrase was even pointing at.
