# OCaml style applicative modules

Poking at whether a module system in the ML tradition (structures, signatures, functors) fits here, and specifically the "applicative" functor semantics OCaml uses by default.

## The pitch, quickly

OCaml modules aren't just namespaces.
A module is a value: it groups types, functions, and other modules under one name.
A signature is a type for a module: it says which of those are visible and what their types are.
A functor is a function from modules to modules: `module Make (X : SIG) : SIG2 = struct ... end`.

Applicative means: `Make(A)` and `Make(A)` are the same type, structurally, as long as `A` is the same module.
Two calls with the same argument produce compatible types, not just isomorphic ones.
This is what lets you write `Set.Make(Int)` in two different files and still pass values between them.
The alternative, generative functors, mint a fresh abstract type on every application, even with identical arguments, so two calls are never type-compatible even if the argument is literally the same module.

## Why this might matter for a language without classes or interfaces-as-values

Type-before-name syntax and structural interfaces are already on the table per prior notes.
Structural interfaces get you ad hoc polymorphism over shape.
They don't get you polymorphism over a *cluster* of types and operations that need to move together, the way `Map.Make(Ord)` bundles a comparison function with the container type it produces.
A single-function interface answers "does this thing support `Draw`."
It doesn't answer "give me a whole Set/Map/Graph implementation parameterized by how equality and ordering work for my key type," where the equality function and the resulting container type are inseparable.

Rough sketch, borrowing existing type-before-name conventions:

```
module Ord for T:
    Int compare(T a, T b)

module Make(Ord for T) for Set(T):
    Set(T) empty
    Set(T) insert(Set(T) s, T x)
    Bool member(Set(T) s, T x)
```

Instantiate:

```
module IntSet = Make(IntOrd)
```

Two instantiations of `Make(IntOrd)` in different files, same signature, same underlying type.
Applicative semantics means the compiler doesn't need whole-program linking to prove that.

## Where this rubs against other settled-ish decisions

HM inference plus parametric polymorphism already gives you `List<T>`.
Do we need functors on top of generics, or is this redundant with generics plus typeclass-like constraints?
Haskell answers this question with typeclasses instead of functors: instance resolution is implicit, there's one canonical `Ord Int`, and you don't get to have two different orderings for the same type live in the same program without newtype wrappers.
OCaml's answer is the opposite: functors are explicit, you can have as many orderings as you want, you just have to say which one at each call site.
That's arguably more consistent with "explicit spelled name at use site" (matches the "no operator overloading, still open" note) but it is a second, heavier way to do polymorphism sitting next to whatever generics syntax the language ends up with.
Two mechanisms for parametric behavior in one language is a real cost, not a neutral design.

Region/arena memory adds a wrinkle no ML-family language has to think about: does a functor instantiation carry an implicit region parameter, or is memory strategy fully orthogonal to the type-level module story?
If a `Set.Make` produces nodes that need to live in a particular arena, does the arena become part of the functor argument, or part of the call site at every operation?
Genuinely unclear, feels like it wants its own note.

## What applicative buys you over generative, concretely

Generative functors are what you get "for free" if functor application is implemented like a struct constructor call: run the body, mint new abstract types, done, no memoization needed.
Applicative functors require the compiler to recognize that two applications with syntactically/semantically equal arguments should unify, which means either:

- structural comparison of the argument module (expensive, or requires argument modules to be fully known at compile time, no separate compilation of the argument), or
- purely syntactic comparison (same path, e.g. `Make(IntOrd)` textually, is applicative-equal; `Make(struct ... end)` with an anonymous struct literal is not, which is exactly OCaml's actual rule)

That syntactic-path restriction is a real wart.
It means the "same argument, same result type" property only holds when the argument is a named, already-elaborated module, not an inline one.
Worth deciding up front whether that's an acceptable rough edge or a dealbreaker before committing to applicative-by-default.

## Dead end or genuinely promising, undecided

Tried to sketch what a functor signature and a plain generic function signature would look like side by side to see if the syntax could visually distinguish them without a new keyword vocabulary.
Didn't land on anything that felt right, kept wanting to reuse `for` for both which erases exactly the distinction that matters (implicit generic instantiation vs. explicit named module application).
Parking that.

Also poked at "what if functors are just functions that take and return records of closures," i.e. desugar the whole module system into ordinary immutable-by-default values instead of a separate compile-time-only construct.
That would dodge the whole applicative-vs-generative question, since equality of two record values is just ordinary structural equality, no special functor-identity rule needed.
Cost: you lose the ability to define new abstract types inside a functor body (a `Make` that hides its internal representation type), because an abstract type isn't a runtime value, it's a compile-time fact.
So "modules as records" gets you the composition ergonomics but not the type abstraction, which might be the actual reason to want a module system rather than passing bags of functions around.
This tension (abstraction needs compile-time types, composition wants runtime values) feels like the real crux, more than applicative-vs-generative.

## Open thread

Never resolved whether this is solving a problem the language has (needing to parameterize a type by a runtime-orderable operation, e.g. containers) versus importing ML ceremony because it's elegant in OCaml specifically.
Go solved the same container-parameterization problem eventually with plain generics and no functor layer at all.
Worth asking, next time this comes up: what's a concrete container or abstraction in this language's own examples that generics alone can't express cleanly?
If no answer, functors might be solving a problem that doesn't exist here.
