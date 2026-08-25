# Haskell-strength type system

Riffing on: what would it mean to have a type system "as strong as Haskell's"?
Zero standing, nothing here is decided, and everything presumes outcomes of open decision records freely.

## First problem: "Haskell-strength" bundles at least five unrelated things

People say this and mean some subset of:

1. Hindley-Milner inference (write no type annotations, ever, and it still checks).
2. Typeclasses (ad-hoc polymorphism resolved by the compiler, not the caller).
3. Higher-kinded types (`Functor f`, `Monad m`, things generic over type constructors, not just types).
4. Purity tracked in the type (`IO a` vs `a`, "a function of type `Int -> Int` cannot launch missiles").
5. Laziness by default.

None of these imply any other.
OCaml has 1 and a weaker 2 without 4 or 5.
Rust has a form of 2 and 3 (traits, GATs) without 1 in the "no annotations anywhere" sense, without 4, without 5.
Idris has 1-4 plus dependent types and is call-by-value by default.
Haskell is really "5 was the accident that made 4 necessary," which is a story worth its own section below.

So "Haskell-strength" as a wishlist entry is five wishlist entries wearing a trenchcoat.
Worth splitting before it ever reaches a sketch.

## Doodle: what does inference actually buy, and what does it cost

The pitch: never write a type signature, the compiler figures it out.

```
let compose f g x = f (g x)
```

No annotation, and yet this has a fully general type.
That is genuinely magic the first time you see it, and it comes from a specific algorithmic property: HM types are principal, meaning there's one most-general type and inference finds it, no search, no backtracking, no heuristics.

The cost shows up the moment the type system grows past HM's sweet spot.
Add typeclasses with functional dependencies, or higher-rank polymorphism (`forall` inside an argument position), or GADTs, and principal types stop existing.
Haskell's actual answer: infer inside a function body, but require an annotation on anything exported or anything using the fancy features.
So "no annotations" was already a half-truth by 2005; `-XScopedTypeVariables` and friends exist because inference alone got overwhelmed.

Alternative framing: inference is a *local* tool (fill in the obvious stuff inside a function) not a *global* one (derive every public API from usage).
Go, Rust, Swift, Kotlin all landed here independently: infer locals, require signatures on function boundaries.
That might just be the correct answer and the "look ma no signatures" Haskell demo is a party trick that the language itself grew out of.

Tension: if the language requires signatures at boundaries anyway, is full HM even worth the implementation cost over "infer local let-bindings, require types on function params and returns"?
The latter is a much smaller algorithm (no unification-with-generalization, no let-polymorphism corner cases) and probably gets 90% of the felt benefit.

## Doodle: typeclasses, three different shapes

```
-- Haskell: classes are open, instances are global, coherence enforced by
-- "one instance per type" so the compiler never has to ask which one you meant
class Show a where
  show :: a -> String

instance Show Bool where
  show True  = "true"
  show False = "false"
```

```
-- Rust: same idea, different name, orphan rules instead of full coherence,
-- and no higher-kinded abstraction over the trait itself without GATs
trait Show {
    fn show(&self) -> String;
}
```

```
-- ML-family: modules and functors instead of typeclasses.
-- Explicit at the call site: you pick the implementation by naming a module,
-- there is no "the compiler searches for it" step at all
module ShowBool : Show with type t = bool = struct
  let show b = if b then "true" else "false"
end
```

The real fork: implicit resolution (typeclass, trait) vs explicit selection (module/functor).
Implicit resolution is why `show` "just works" everywhere and also why Haskell has a whole subfield of papers on instance resolution ambiguity, overlapping instances, and incoherence when two libraries both define an instance for the same type.
Explicit selection never has that problem because you always say which one, and pays for it in ceremony at every call site that wants polymorphism.

Doodle for a hybrid: implicit resolution but *only* for a small closed set of compiler-known classes (equality, ordering, display, maybe `Iterable`), and explicit passing (dictionaries as plain values, i.e. "just pass the function") for anything user-defined.
This is roughly Go's answer post-1.18 for the built-in comparable constraint plus normal generics for everything else, and it dodges coherence entirely by never letting instance search be open-world.
Cost: the built-in classes feel special and user classes feel second-class, literally.

## Doodle: higher-kinded types, and whether they're worth their weight

```
class Functor f where
  fmap :: (a -> b) -> f a -> f b
```

`f` here is not a type, it's a type *constructor* — `Maybe`, `[]`, `IO`, something that becomes a type once applied to an argument.
This is the thing most mainstream languages skip (Java generics, Go generics, TypeScript, Rust pre-GATs) because it's a kind system on top of a type system, and kind errors ("you used `Maybe` where `Maybe Int` was expected") are a new category of mistake to explain to users.

What you actually lose without it: `Functor`, `Applicative`, `Monad`, `Traversable` as *reusable abstractions*.
You can still have `map` on `List` and `map` on `Option` and `map` on `Result`, you just can't write one function that's generic over "any of these," because there's nothing to quantify over.
Every language without HKTs eventually reinvents this pain: C#'s LINQ is Monad-shaped but hand-duplicated per type via extension methods; Rust's `Iterator` trait is a single hardcoded monad because the trait system can't express "generic over which container."

Question worth sitting with: is that duplication actually bad?
Rust proves you can get extremely far by hand-specializing the handful of shapes people actually use (`Option`, `Result`, `Iterator`, `Future`) and never generalizing.
Users arguably *prefer* `.map()` meaning something concrete over `.map()` meaning "apply `fmap` for whichever `Functor` instance resolves here," because the error messages for the concrete version are a type mismatch and the error messages for the generic version are a wall of constraint-solver output.

Dead-end-shaped thought: HKTs might be a case where the "strength" is real but the *cost is born by error messages*, and nobody has fully solved making HKT error messages good.
That is a serious strike against importing them wholesale, independent of whether the feature is theoretically nice.

## Doodle: purity in the type, and the laziness entanglement

Haskell's `IO a` is the thing people actually mean by "Haskell-strength" half the time: a function with no `IO` in its type cannot do IO, full stop, checked by the compiler.

```
pure   :: Int -> Int          -- cannot touch the network
impure :: Int -> IO Int       -- might
```

Here's the tangle worth chasing: Haskell *needs* this distinction because it's lazy, and laziness makes evaluation order unobservable, and unobservable evaluation order makes "when does the side effect happen" a real question with no good answer unless effects are pushed into a type that sequences them explicitly (`IO`, via the monad's bind forcing an order).
Purity-tracking wasn't a preference, it was forced by laziness having already been chosen.

So: does a language want the `IO` type without wanting laziness?
Yes, probably — this is close to what Koka, Unison, and effect-system languages do (algebraic effects instead of a monad stack), and none of them are lazy by default.
In a strict language, evaluation order is already observable and sequenced by the syntax, so you don't *need* a type to smuggle ordering back in.
You could still *want* one, purely for the "the type signature tells you if this can do IO" property, decoupled entirely from why Haskell has it.

This suggests: "purity tracked in types" and "laziness" are not just separable, one of them is a historical cause of the other, and only one of the two (purity-in-types) is likely to be independently worth wanting.
Laziness-by-default without the pull of "otherwise how do you sequence IO" is a much harder sell: space leaks that require a PhD to debug (`foldl` vs `foldl'`), and the joke that half of Haskell education is "here is how to force things when the language's core feature bites you."

## Doodle: what "strength" costs at the error-message layer

Every feature above trades compile-time guarantees for compile-time *legibility*.
A monomorphic, first-order, strict language with no typeclasses has type errors that are always "expected `Int`, got `String`, right here."
Add let-polymorphism: errors can now point at a *use site* far from the actual mistake, because inference generalized too eagerly.
Add typeclass constraints: errors become "no instance for `(Foldable t, Monoid m) => ...`" which is unreadable to a newcomer and not much better for an expert three weeks removed from writing the code.
Add HKTs: errors gain a kind dimension, so a mismatch can be "you have `f a` but needed `f`," a sentence that requires already understanding the feature to parse.

This is the actual price of "Haskell-strength," and it's not expressiveness, it's *pedagogy*.
Rust spent enormous effort on this exact problem (best of any mainstream statically typed language, arguably) and even Rust's trait-bound errors get long.
Any decision record proposing typeclasses or HKTs should probably budget "and here is the error-message strategy" as a first-class deliverable, not an afterthought, because the feature is genuinely worthless if nobody can debug why their code doesn't typecheck.

## A different cut: what if "strength" means totality, not expressiveness

Tangent, but worth writing down before it's lost: nothing above touches whether functions have to be total (defined for every input, no crashes, no infinite loops as a type-level fact).
Idris and Agda push here; Haskell famously does not (`head []` blows up at runtime, `Maybe`'s existence didn't stop anyone from calling `fromJust`).
A language could plausibly claim to be "stronger than Haskell" by making partiality a type-level marker (`partial fn head`) rather than by having typeclasses at all.
That's a completely orthogonal axis from everything doodled above, arguably a more honest reading of "strong type system" (rules out more programs at compile time) than typeclasses or HKTs, which are about *code reuse* more than about *safety*.

Filing this because "strength" as raw safety (totality, effect tracking, refinement types catching array-bounds) and "strength" as expressive power (typeclasses, HKTs, letting one function serve many types) are pulling in different directions and get conflated by the word "Haskell."

## Dead ends worth keeping

- **Import the whole GHC feature set as the bar.** Dies immediately: GHC has three decades and a dozen extension flags (`RankNTypes`, `TypeFamilies`, `DataKinds`, ...) that are each their own research contribution. "Haskell-strength" cannot mean "GHC," it has to mean some specific, nameable subset, chosen on purpose.
- **Laziness by default, imported for free alongside purity-tracking.** Dies per above: it's Haskell's specific solution to a problem (sequencing effects under unobservable evaluation order) that a strict language doesn't have in the first place. Taking laziness because it "comes with" the strong type system mistakes a historical accident for a package deal.
- **Implicit global instance resolution for user-defined typeclasses.** Not fully dead, but suspect: the coherence/orphan-instance mess is a 30-year-old unsolved-to-everyone's-satisfaction problem in the Haskell/Rust ecosystems both. A hybrid (implicit for a closed builtin set, explicit for everything else) dodges it by construction and might be worth wanting on its own, independent of "as strong as Haskell."

## Open tensions I did not resolve

- Whether full HM ("no signatures anywhere") is worth its implementation and error-message complexity over "infer locals, require signatures at boundaries," which most modern languages converged on independently.
- Whether HKTs are worth adopting given nobody has good error messages for them yet, versus hand-specializing the 3-4 shapes people actually reuse (`Option`, `Result`, `Iterator`-analog).
- Whether "purity in types" is wanted for its own sake (readable function signatures, no-IO guarantees) decoupled from monads-as-the-mechanism, i.e. could this be effects-system shaped instead of `IO`-monad shaped.
- Whether "strength" for this language should mean expressive power (typeclasses/HKT, code reuse) or safety (totality, effect tracking), because the wishlist entry that prompted this note doesn't say and the two pull toward different feature sets.

## Prior art to actually read

- Hindley-Milner and principal types: why inference works until `RankNTypes`/GADTs break principality, and what Haskell requires annotations for as a result.
- Rust traits vs Haskell typeclasses vs OCaml modules, as three answers to "how does ad-hoc polymorphism get resolved," with orphan rules / coherence as the actual crux.
- Koka and Unison for algebraic effects as a strict-language, non-monadic answer to "track what a function can do in its type."
- Idris/Agda for totality checking as an orthogonal "strength" axis from typeclasses/HKTs.
- Elm's error-message design and Rust's trait-error diagnostics, as the two best-documented attempts at making constraint-based type errors legible to non-experts.
