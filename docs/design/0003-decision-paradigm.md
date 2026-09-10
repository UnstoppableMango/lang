---
title: 'Decision: Paradigm'
stage: 1
status: active
created: 2026-08-27
updated: 2026-08-27
---

# Decision: Paradigm

## Motivation

A language's paradigm is the answer to two linked questions: how is data organized relative to the behavior that acts on it, and how is that behavior dispatched.
Object-oriented languages bundle the two behind a class and dispatch through inheritance or an interface.
Functional languages keep them apart: data is a shape, behavior is a free function, and dispatch is pattern matching over the shape.
Procedural languages keep them apart too, but without functional programming's emphasis on immutability or exhaustiveness.
Data-oriented design goes further and asks the question in terms of a whole batch of values rather than one value's shape at all.
This record is scoped to that organizing question only.
It is not the type system's nominal-versus-structural axis, explored in `docs/notes/duck-typing.md`, and it is not the memory-strategy record's question of who owns a value, though the two interact and neither may be presumed here.

Several wishlist entries and notes already lean somewhere without anyone having said so out loud.
Expression-oriented control flow, exhaustive pattern matching, "no null" as a real `Option` type, and structural interfaces all read most naturally in a language whose data is algebraic and whose control flow composes as expressions, which is the functional candidate below.
The design philosophy section's "few ways to do one thing" reads most naturally as a caution against the fusion candidate below, which deliberately keeps two idioms alive.
Nobody has said which of these readings is the decision, and until this record exists, every doc that touches data modeling is guessing.

### Scenario: the compiler's own AST

`hack/hello.lang` compiles today through a parser that builds an AST, a lowering pass, and LLVM IR generation.
An AST node is a small, closed, recursively nested set of variants: a literal, a binary operation, a call, and so on.
Every pass over it needs to inspect which variant a node is and act differently per variant, then recurse into its children.
This is the paradigm's central data shape in the one program that exists in this repository today, so any candidate has to say plainly how it would represent and traverse this specific value.

### Scenario: satisfying an interface you did not write against

`docs/notes/duck-typing.md` poses a `Logger` interface and a third-party `Widget` type that already has a `Log(string)` method, with the wishlist already committed to structural interface satisfaction checked statically.
Whatever the paradigm, something has to explain how a value written with no knowledge of `Logger` nonetheless satisfies it, and whether the mechanism differs when the value's behavior lives in a method versus a free function.

### Scenario: transforming many nodes at once

A type-checking or lowering pass over a large file touches every node once, in bulk, not as an isolated value handled one at a time.
Whether the language's paradigm treats "many values of the same shape processed together" as an ordinary case of the same abstraction used for one value, or as a distinct concern with its own layout rules, changes how the compiler itself would eventually be written in its own language, should it ever self-host.

## Rough shape

Five candidate directions, sketched.
As with the other two decision records, these are not necessarily exclusive in the limit, and the syntax below is illustrative only.

### A. Object-oriented

Data and behavior are bundled behind a class or struct with methods; dispatch for shared behavior across variants goes through inheritance or an interface implemented at the type's definition site.

```
class Shape {
    fn area() -> Float
}
class Circle extends Shape {
    radius: Float
    fn area() -> Float = 3.14159 * radius * radius
}
```

Buys: the most familiar model to most working programmers, and encapsulation gives a natural unit for hiding a representation behind its methods.
Costs: inheritance hierarchies are the textbook case that fights retroactive conformance, the exact complaint the duck-typing note raises against nominal designs; adding a new variant of a closed hierarchy means editing every subclass rather than being told by the compiler which `match` arms are missing; and it sits awkwardly next to wishlist items that already lean away from nominal inheritance toward structural interfaces.

### B. Functional

Values are immutable by default.
Data is modeled as algebraic data types, sums of products, and behavior lives in free functions that pattern match on shape rather than methods attached to it.

```
type Expr = Lit(Int) | Add(Expr, Expr) | Mul(Expr, Expr)

fn eval(e: Expr) -> Int = match e {
    Lit(n) => n,
    Add(a, b) => eval(a) + eval(b),
    Mul(a, b) => eval(a) * eval(b),
}
```

Buys: this is exactly the AST scenario's shape, and it is the candidate that requires fighting none of the wishlist entries already leaning this way (expression-oriented control flow, no-null as a real `Option`, exhaustive pattern matching); adding a new operation over `Expr` costs nothing, and adding a new variant is a compile error at every non-exhaustive `match`, the flip side of candidate A's failure mode.
Costs: adding a new variant later means touching every existing pattern match over that type, the flip side of the same expression-problem tradeoff; encapsulation is weaker without a language-level notion of hiding a representation behind an interface; and "mostly immutable" is a precondition the memory-strategy record's candidate C (compiler-inserted reference counting) already names as something it leans on, so choosing this candidate here narrows that record's remaining live candidates rather than leaving them open.

### C. Procedural, Go-shaped

Structs are plain data with no attached methods of their own; functions are free; interfaces are satisfied structurally, with no declared relationship and no class hierarchy.
Mutation is unremarkable and is not fenced off behind a purity or ownership system by default.

```
struct Circle { radius: Float }

fn area(c: Circle) -> Float {
    3.14159 * c.radius * c.radius
}
```

Buys: the smallest learning burden of any candidate, directly matching the design philosophy's "few ways to do one thing" and Go's cited design stance; it sidesteps the expression problem's sharpest edge by keeping data and behavior only loosely coupled through structural interfaces, matching the duck-typing note's own conclusion; and it imposes the least on the memory-strategy and concurrency records, since a mutation-by-default answer and an immutability-by-default answer are different defaults on the same mechanism rather than different mechanisms.
Costs: gives up the ergonomic wins the wishlist already lists that read best under immutability and exhaustiveness; no-null's "absence is a real type, not a sentinel" story and exhaustive pattern matching both compose more naturally in an expression-oriented, sum-type-heavy language than a statement-oriented one; and "unremarkable mutation by default" is in direct tension with the concurrency wishlist item wanting data races to be hard to express by default, not merely possible to avoid.

### D. Fusion, Scala-shaped

Functions are values.
Data can carry attached methods, but pattern matching over algebraic data is equally first-class, so a caller can reach for either idiom depending on which one fits.
Immutability is encouraged but not enforced.

```
type Shape = Circle(Float) | Rect(Float, Float)

fn Shape.area() -> Float = match self {
    Circle(r) => 3.14159 * r * r,
    Rect(w, h) => w * h,
}
```

Buys: does not force a choice between B and C; algebraic data and pattern matching stay available for compiler-internal, closed-shape data like the AST scenario, while method syntax and structural interfaces stay available for open, extensible, externally-facing polymorphism; Scala is a direct existence proof that this reconciliation can be made to work in one typed language rather than staying an unresolved tension.
Costs: two idioms for the same job, a free function taking a value first versus a method dispatched through it, is exactly the "few ways to do one thing" violation the design philosophy section warns against; and Scala's own well-documented reputation for a large, hard-to-predict feature interaction surface is a direct consequence of keeping every combination available at once rather than choosing.

### E. Data-oriented

The primary unit of design is not a single value's shape but a batch of values laid out for the transformation being performed.
Types describe layout; transformations operate over columns or arrays rather than one value at a time.

```
struct ExprKinds { kinds: Array<Tag> }
struct ExprOperands { lhs: Array<NodeId>, rhs: Array<NodeId> }

fn evalAll(kinds: ExprKinds, ops: ExprOperands) -> Array<Int> {
    // one pass over every node, not one call per node
}
```

Buys: the closest fit to the bulk-transformation scenario, and a natural pairing with an arena or region memory-strategy candidate already on the table in the memory-strategy record, since a batch of same-shaped values in one region is exactly the case regions are cheapest for.
Costs: the least explored candidate in this project's own notes, with no note dedicated to it the way arenas, duck typing, or functional implicit return each have one; it actively fights the AST-as-recursive-sum-type shape the first scenario and most of the wishlist's type-system items presume; and it forces every API to expose its bulk layout rather than a single value's shape, a much larger and stranger surface for a general-purpose language to ask of every author than for a game engine's hot loop, which is the domain data-oriented design comes from.

## Prior art

- **Language design as engineering tradeoffs, functional versus object-oriented.**
  Object-oriented decomposition is typically operation-first, an interface is declared and classes implement it, while functional decomposition is data-first, a shape is declared and functions pattern match over it ([Programming language trade-offs](https://www.garfieldtech.com/blog/language-tradeoffs); general survey of multi-paradigm languages including OCaml, Swift, Rust, Scala, F#, Kotlin).
  The relevant lesson: a language committing to structural interfaces and immutability, as this project's wishlist already does in places, is implicitly leaning data-first without having said so, which is exactly the gap this record closes.

- **Odersky and Rompf, "Unifying Functional and Object-Oriented Programming with Scala"** (CACM, April 2014, [full text](https://cacm.acm.org/magazines/2014/4/173220-unifying-functional-and-object-oriented-programming-with-scala/fulltext)).
  Scala's thesis is that functional and object-oriented programming are two sides of one coin, to be identified as much as possible, with function values themselves being objects.
  Direct precedent for candidate D, and its costs are equally well documented: Scala's combinatorial feature surface is the standing counterargument.

- **The Rust Programming Language, "Object-Oriented Programming Features"** ([doc.rust-lang.org/book/ch18](https://doc.rust-lang.org/book/ch18-00-oop.html)).
  Rust deliberately has no classes or implementation inheritance; shared behavior is expressed through traits implemented on plain structs and enums, interface inheritance rather than implementation inheritance.
  Relevant to candidates A and C as a worked example of a systems language that chose structural, trait-based dispatch over class hierarchies, without going as far as candidate B's algebraic-data-and-pattern-matching emphasis.

- **Algebraic data types and pattern matching in language design** ([kindatechnical.com, Programming Language Design and Evolution](https://kindatechnical.com/programming-language-design-evolution/algebraic-data-types-and-pattern-matching.html)).
  Sum types plus product types give a concise way to model composite and recursive data, and pattern matching cooperates with exhaustiveness checking to remove whole failure classes, no downcasting, no null-shaped gaps, at the cost of the expression problem's other edge: adding a variant means revisiting every match.
  Directly describes candidate B's mechanics and its named cost.

- **Data-oriented design** ([Wikipedia](https://en.wikipedia.org/wiki/Data-oriented_design); Mike Acton, "Data-Oriented Design and C++", CppCon 2014).
  Motivated by CPU cache behavior rather than domain modeling: a parallel array or structure-of-arrays layout is contrasted with the array-of-structures shape object-oriented design defaults to, with Acton's framing that object-oriented modeling optimizes for an abstract real-world picture of the system rather than for the data actually being transformed.
  The source domain is game engines, not general-purpose languages, which is exactly why candidate E is flagged as the least explored fit here.

Where the search was thin: nothing found treats "which paradigm should a language's own self-hosting compiler be written in" as a factor in choosing that language's paradigm, most language design writing treats the reference implementation's language as an accident of the author's preference rather than a constraint the target language's paradigm should satisfy.
That gap is recorded as an open question rather than folded into the candidates above, since this project's compiler is not self-hosted and may never need to be.

## Open questions

1. Is the compiler's own AST, an unavoidable recursive sum-shaped value today, the paradigm's central case, or is the reference compiler free to model data however it likes internally regardless of what the surface language exposes to its own users?
1. Is mutation unremarkable by default, as in candidates C and D, or fenced behind an explicit marker, as candidate B's immutable-by-default stance implies, given that the concurrency wishlist item wants data races hard to express and several notes already lean on immutability for unrelated reasons?
1. Does the language commit to exactly one idiom for attaching behavior to data, a free function taking the value first or a method dispatched through it, or does it allow both, as candidate D does, at the cost the design philosophy section explicitly warns against?
1. How does retroactive interface conformance, already committed to structurally per the wishlist, work under whichever candidate is chosen: a method syntax and a free function each need their own concrete account of how a value satisfies an interface it was not written against.
1. Is "paradigm" scoped, as this record assumes, to data/behavior organization and dispatch only, leaving the type system's nominal-versus-structural axis and the memory-strategy record's value-semantics candidate as separate decisions this record must not presume?
1. **Blocked on an undecided foundational decision:** the memory-strategy record (`0001`) is narrowed by whichever candidate wins here.
   Its candidate C, compiler-inserted reference counting, already names "mostly immutable data" as a precondition, and its candidate E, value semantics with second-class references, is close enough to a paradigm choice that the two records may be deciding the same thing from different directions without either admitting it.
1. **Blocked on an undecided foundational decision:** the compilation-model record (`0002`) names this record as a block on its own open question about dispatch, since Julia's JIT design leans on multiple dispatch to decide what to specialize, and a paradigm with different dispatch or genericity mechanics changes what that record's candidate C or a comptime-driven monomorphization scheme has to do.
1. **Blocked on an undecided foundational decision:** the concurrency model has no decision record.
   Whichever candidate answers open question 2 here, unremarkable mutation or immutable-by-default, directly shapes whether the concurrency wishlist item's "data races hard to express" goal is inherited for free or has to be built separately.

## Advancement record

- 2026-08-27, gate 0 → 1: sketched from a new wishlist entry ("A programming paradigm chosen as an explicit decision") added under Foundational decisions; five candidate directions surveyed, prior art cited across functional/OOP tradeoffs, Scala's fusion, Rust's trait-based non-class design, algebraic data types and pattern matching, and data-oriented design; open questions recorded including three blocks on undecided or narrowly-overlapping foundational decisions.

## Changelog
