# Language AST exposed in the public API

Prompted by \[[dev-tooling-philosophy]\]'s doodle: "define the AST/IR once, in a form the compiler, formatter, linter, and language server all import as a library."
That doodle assumed the shared AST question without asking it directly.
This note asks it directly: what does it actually mean for this language's AST to be a *public*, versioned, importable artifact, not just an internal implementation detail the compiler happens to have?

## The claim to poke at

Two very different things get called "the AST is public":

1. The compiler is structured as a library, and that library's AST types are part of its public API surface, the way `go/ast` is importable from any Go program.
1. The AST's shape is part of the language's own specification, the way it might show up in `docs/design/NNNN-*.md` as a stable, documented data structure a spec reader could implement against.

These are not the same claim and they don't have to be decided together.
A language could have (1) without (2): the reference compiler's AST types are importable Go/Rust/whatever structs, but nothing says another implementation has to produce the same shape.
A language could also have (2) without (1) in principle, though it's hard to imagine wanting to spec an AST shape without also wanting to use it, so this direction feels like a dead end worth naming and moving past.

## Doodle: what "public" could buy you

```
import lang/ast

tree, err := ast.Parse(src)
for node := range ast.Walk(tree) {
    switch n := node.(type) {
    case *ast.FuncDecl:
        ...
    }
}
```

Prior art doodle, `go/ast` shaped.
This is the shape that makes eleven tools (from the tooling note) cheap: a formatter, a linter, a doc generator, and a codemod tool are all just "parse, walk, maybe rewrite, maybe print."
Rust's `syn` + `proc-macro2` is the other lineage: crates outside the compiler itself parse Rust-shaped tokens, and an entire macro ecosystem exists downstream of that being possible without forking rustc.

## Doodle: the AST as the thing macros/comptime see

The wishlist already has "compile-time code execution rather than a separate textual macro language."
If that lands, the AST-as-public-API question and the metaprogramming question are the same question wearing two hats.
A `comptime` block or a derive macro needs *some* view of the syntax tree to operate on.

```
derive Eq, Show for Point { x: Int, y: Int }
```

For `derive` to generate an `Eq` impl, something needs to walk `Point`'s field list.
If the AST is public and stable, that something can be ordinary user code written against the same library the compiler uses internally.
If the AST is a private compiler implementation detail, `derive` has to be a compiler intrinsic, hardcoded per-derivable-trait, extensible only by patching the compiler itself.
Zig's comptime and Rust's proc-macros disagree on exactly this axis: Zig leans toward "the compiler already speaks your language at compile time, no separate AST-walking API needed," Rust leans toward "here is a crate (`syn`) that hands you a real tree to pattern-match on."

## Doodle: versioning the AST is versioning the language, twice

Every AST node is a promise.
Add a field to `FuncDecl`, and every external tool matching on that struct either breaks or silently ignores the new field depending on the language's exhaustiveness rules.
Go dodges this partly by not having exhaustive switch by default, so `go/ast` has grown fields for two decades without breaking every downstream `switch`.
This language's wishlist already wants "pattern matching with compiler-checked exhaustiveness," which is in direct tension with a public AST that needs to grow new variants over time without breaking every consumer.

```
match stmt {
    IfStmt(c, t, e) => ...
    ForStmt(...) => ...
    // add MatchStmt later: every exhaustive match on Stmt in the
    // entire ecosystem now fails to compile until updated
}
```

Possible outs, none obviously right:

- Non-exhaustive-by-default for AST-shaped enums specifically, carving out an exception to the general exhaustiveness story.
- A closed AST core plus an explicit "unknown/extension node" variant, so exhaustive matches stay exhaustive but new syntax rides in a bucket that has to be handled generically anyway (which defeats a lot of the benefit of exposing structured nodes to begin with).
- Version the AST crate independently of the language (`lang-ast v3` vs `lang v1`), so an old formatter can keep working against an old tree shape while the compiler internally uses a newer one and translates at the boundary. Cost: now there are two things to keep in sync forever, which is exactly the "shared AST" doodle's whole premise, undone.

This tension doesn't resolve here, it's a live thread between this note and whatever decision record eventually covers exhaustiveness.

## Tangent: does "public AST" even mean "public syntax tree," or something looser?

Tree-sitter's answer is interesting because it sidesteps the exhaustiveness problem entirely: the grammar is public, but consumers query it with an S-expression pattern language, not by matching on generated types in their own language.
A tree-sitter query never "breaks" the way an exhaustive match breaks, it just fails to match new node kinds it wasn't written for.
That buys forward-compatibility at the cost of losing static exhaustiveness checking entirely, the query language can't tell you "you forgot to handle the new node kind," it can only silently not match it.

Worth asking: is the goal "expose the AST as native language values" (the `go/ast` lineage) or "expose the AST as a queryable external structure" (the tree-sitter lineage)?
These read like the same idea from a distance and are quite different in what they cost and what they buy.

## Tangent: this is also an embedding question

The wishlist has "embeddability of the language or its runtime into a host program," Lua and Wasm components as prior art.
An embedded interpreter that only exposes "run this string, get this value back" never needs a public AST.
An embedded interpreter meant for things like "let the host program inspect or rewrite scripts before running them" (Lua's `luaL_loadstring` plus bytecode dumps sit adjacent to this) does.
This is a second, mostly independent reason a public AST might matter, separate from the tooling-ecosystem reason this note started from.
Two different constituencies want the same artifact for different reasons: tool authors (formatter, linter, LSP) and embedders (host programs).
Not obviously true that satisfying one satisfies the other, an AST shaped for `gofmt`-style rewriting and an AST shaped for a sandboxed host to safely introspect a guest script may want different guarantees, especially around whether the tree can express the whole language including any parts a host would want to forbid a guest from constructing.

## Dead ends, recorded so I stop rediscovering them

- **"AST as the spec."** Tried this thought and it collapses immediately: a spec that says "the AST looks like this Rust struct" ties the specification of the language to one implementation language's type system forever. If a second implementation exists in a different host language, "the AST" can't literally be that struct, at best it's isomorphic to it. Spec-vs-reference-compiler (\[[spec-vs-reference-compiler]\]) already covers this ground more generally, this is just a specific instance of it.
- **"Solve exhaustiveness-vs-public-AST by never adding new syntax."** Not a real answer, the language will grow, the wishlist alone has a dozen half-formed features that would each need new node kinds.

## Threads worth pulling later

- Whether "AST" and "IR" should even be the same public artifact, or whether the AST stays private/unstable (close to source, churns with syntax changes) while a separate, more stable IR is the actual public/tool-facing layer. Go doesn't really have this split (`go/ast` is close to the source tree); Rust sort of does (`syn`'s tree vs rustc's HIR vs MIR, with only the outermost being an ecosystem-facing public crate).
- Whether the derive/comptime story (once that decision record exists) settles the exhaustiveness tension above one way, since "macros can pattern-match the AST" and "the AST can grow without breaking macros" are in tension and metaprogramming design will have opinions.
- What a "public AST" costs implementation-language choice: an AST-as-library story is easiest in a language with cheap tagged unions/ADTs and reasonable serialization story, which loops back to the still-fully-open "implementation language" foundational decision.
