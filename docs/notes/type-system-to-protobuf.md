# Type system compiles to protobuf

Riffing on: the type system compiles to protobuf definitions, and `.proto` files are a toggle-able build output.
Zero standing, presumes outcomes of open decision records freely.

## First fork: is proto an output, or is proto the model?

The one-line pitch hides two completely different features.

**Reading A, proto as a build target.**
`lang build --emit proto` walks the type declarations and writes `.proto` files next to the binary.
Same shape as `--emit llvm-ir`, `--emit asm`, `--emit docs`.
Purely additive, changes nothing about what programs are legal.

**Reading B, proto as the semantic ceiling.**
The type system is *designed* so that every type is expressible in protobuf, because proto is the interchange story.
This is a foundational constraint: no sum types richer than `oneof`, no generics, presence semantics inherited from proto3, open enums.

A is a tooling feature that could ship in a weekend once the type graph is walkable.
B is a decision record that constrains the language forever.
The phrasing "type system compiles to protobuf" leans B, the phrasing "toggle-able build output" leans A, and the rest of this note is mostly about how they fight.

## The toggle principle worth extracting early

If flipping `--emit proto` can cause a *compile error*, the flag is not an output toggle, it is a dialect switch.

Consider `lang build` succeeding and `lang build --emit proto` failing with "type `Parser<T>` cannot be expressed in protobuf."
Now the set of legal programs depends on a build flag.
Everyone who might ever want proto output writes proto-expressible code always, and the "optional" output has silently become reading B by social pressure.

Two ways out:

1. **Emission is total.** Every type gets *some* proto encoding, even a lossy or opaque one (`bytes` blobs for anything unrepresentable). Never errors, never constrains. Cheap, and produces a schema nobody would want to consume from another language.
1. **Opt-in per type.** Only types marked for it must be proto-expressible, and that mark is checked on *every* build, not just proto builds. The check is unconditional, the *emission* is what toggles.

Option 2 looks right and it is the whole answer to the framing tension: the constraint is per-type and always on, the file writing is what has a flag.

```
@wire
struct User {
  name: String
  email: String
}

struct ParserState { ... }   // no @wire, never checked, never emitted
```

## The field number problem, which is the actual hard part

Protobuf's entire compatibility model rests on field numbers, which are human-assigned stable identity.
A type declaration in a normal language has no such thing.
Something has to invent them, and every option is bad in a different way.

**Declaration order.**
`name` is 1, `email` is 2.
Reordering two lines in a struct is now a wire-breaking change with no visible diff in meaning.
A formatter that sorts fields alphabetically would be a compatibility disaster.
This one dies fast, and it dies harder against the wishlist entry about a canonical formatter and reformatting never producing a semantic diff.

**Explicit annotation.**

```
@wire
struct User {
  @1 name: String
  @2 email: String
}
```

Honest, and exactly what protobuf-net does in C# (`[ProtoMember(1)]`).
Cost: protobuf's numbering vocabulary is now part of the language's surface syntax, and users have to know about the reserved 19000-19999 range and the 1-15 single-byte tag optimization to write good types.
Also every field of every wire type carries a magic number forever, including the ones nobody ever renumbers.

**Derive from the name.**
Hash `"email"` into the field number space.
Stable under reordering, which fixes the formatter problem.
Breaks on rename, which is arguably correct (rename is a breaking change) and arguably terrible (rename is not supposed to be a wire event).
Collisions are a real problem in a 2^29 space with birthday math over a large schema, and there is no good recovery when two fields in one message collide.

**A lockfile.**
The compiler maintains `wire.lock` mapping fully-qualified field names to numbers.
New fields get the next free number, deleted fields become `reserved` entries so the number is never reused.
Checked into version control, reviewed like any lockfile, and a merge conflict in it is exactly the signal you want when two branches both added a field.

The lockfile answer is clearly the best one and it is not even close.
It also makes the compiler responsible for a persistent, human-reviewable side artifact, which is a bigger commitment than it sounds and lands right next to \[[version-control-integration]\] and its stable per-definition identity idea.
If that identity mechanism exists for blame/diff purposes anyway, wire numbers might just be a projection of it rather than a second scheme.

## Where the type systems disagree

Assume the language wants the things already on the wishlist.
Proto3 does not have them.

| Language wants | Proto3 has | Result |
| --- | --- | --- |
| No null, optionality as a real type | Every field optional-ish, defaults on absence | `Option<T>` needs a representation, presence semantics leak |
| Sum types with exhaustive matching | `oneof`, which can also be *unset* | Every emitted sum grows a "none" case |
| Generics | Nothing | Monomorphize into `OptionInt32`, `OptionString`, ... name explosion |
| Closed enums | Open enums, unknown values pass through | Decoded enums are not exhaustively matchable |
| Structural typing | Nominal, identity by field number | Two identical shapes are two different messages |

The `oneof` mismatch is the sharpest one.

```
enum Shape {
  Circle(f64)
  Rect(f64, f64)
}
```

becomes a `oneof` of two generated message types, plus a third state the language does not have: nothing set.
So decoding cannot produce a `Shape`, it produces a `Result<Shape, DecodeError>`, which is fine and consistent with \[[no-exceptions-explicit-errors]\].
But then the *forward compatibility* story protobuf is famous for evaporates: an old binary receiving a new variant gets a decode error rather than gracefully ignoring it, which is the whole point of proto in a rolling deploy.
The alternative is an implicit `Unknown(bytes)` arm on every decoded sum type, which means the type you match on is not the type you declared.

No resolution here, but the shape of the problem is clear: protobuf's compatibility model assumes the receiver tolerates unknowns, and exhaustive matching assumes the receiver knows everything.
Those are not compatible goals, they are the same tradeoff appearing in two places.

## Two notions of "breaking change", now both in the compiler

Once proto is emitted, the compiler knows about API evolution along two axes that disagree constantly:

- Rename a field: source-breaking, wire-compatible (number unchanged).
- Reorder fields: source-neutral, wire-breaking under order-derived numbering, wire-neutral under a lockfile.
- Widen `i32` to `i64`: wire-compatible by protobuf's varint rules, source-breaking for anyone pattern matching on the width.
- Add an enum variant: source-breaking (exhaustiveness) and wire-... it depends on whether the receiver is old.

This is either the best argument *for* the feature or the strongest argument that it belongs outside the compiler.
For: a `lang vet` that says "this commit is a wire-breaking change to `User`" is a genuinely great tool, and only the compiler has the information to say it.
It is `buf breaking` without needing to hand-maintain the `.proto`.
Against: teaching the compiler two incompatible definitions of compatibility is a lot of conceptual weight for an optional output.

## The counter-argument that might kill the whole thing

The wishlist already wants the AST/IR as a public, versioned library (\[[ast-in-public-api]\]).
If that exists, "emit proto" is an external tool that imports the type graph and prints text.
Maybe three hundred lines.
It does not need to be a build flag, it does not need to be in the compiler, and it does not need a decision record.

And once it is external, why protobuf specifically?
The same walk emits JSON Schema, Avro, OpenAPI, Cap'n Proto, FlatBuffers, or a TypeScript `.d.ts`.
Baking protobuf into the compiler is a permanent coupling to one company's wire format at a moment when the language has not even chosen a paradigm.

The counter-counter: emission has to be *coherent with the build*.
It needs the same name resolution, the same module graph, the same monomorphization decisions, and it needs to be hermetic and cached like every other output (\[[nix-first-build-system]\]).
A third-party tool re-deriving all of that will drift.
"Toggle-able build output" might be exactly the right compromise: the compiler owns emission because it owns the type graph, but emission is one plugin-shaped consumer of a generic interface rather than a hardcoded protobuf backend.

That reframing feels like the real idea worth keeping: **serialization schema emission as a general build output, with protobuf as the first backend.**

## What is actually attractive here

Worth naming, because the mechanism above is complicated enough to lose the motivation.

1. **The language is the IDL.** Schema-first without a schema language. Define types once in real code, and every other language in the org consumes them through their existing protoc pipeline. Same pitch as tRPC for TypeScript, without needing everyone else to be TypeScript.
1. **Wire format for free**, which the wishlist already wants under derivable serialization boilerplate.
1. **Compatibility linting**, per above.
1. **A polyglot on-ramp.** A brand new language nobody uses can still be dropped into an existing gRPC mesh, because the thing it emits is the thing everyone already speaks. That is a real adoption argument for a language with zero ecosystem.

Point 4 is the strongest one and it is not a type system argument at all, it is a distribution argument.

## Doodles on the surface syntax

```
// A: attribute, lockfile-numbered
@wire
struct User { name: String, email: String }
```

```
// B: attribute, explicit numbers, protobuf-net shaped
@wire
struct User { @1 name: String, @2 email: String }
```

```
// C: a marker interface, structural satisfaction does the work
struct User : Wire { name: String, email: String }
```

```
// D: no marker at all, the manifest names the exported types
[wire]
export = ["User", "Order", "LineItem"]
```

D is interesting because it puts the wire surface in one reviewable place rather than scattered across the codebase, and because the set of types you expose over the network is genuinely a project-level architectural fact, not a property of the struct.
It also fits the "capabilities declared at the boundary" instinct already on the wishlist.
Downside: the constraint check has spooky action at a distance, and the error points at a manifest line, not at the offending field.

## Viral constraints

`@wire struct Order { items: List<Item> }` requires `Item` to be wire-expressible too.
So "is wire-expressible" is a property computed structurally and propagated through the whole reachable type graph, closer to Rust's auto traits (`Send`, `Sync`) than to a normal derive.
Two consequences:

- The error message must report the *chain*, not the leaf. "`Order` is `@wire` but `Order.items[].callback` is a function type" is the only useful form of this diagnostic. Error messages are already flagged in \[[haskell-strength-type-system]\] as the real cost of any propagating type property, and this is the same tax.
- Somebody will want `@wire` on a type containing an interface value, and that is dynamic dispatch over the wire, i.e. `google.protobuf.Any` plus a type registry. That is a rabbit hole with its own decision record hiding in it.

## Dead ends worth keeping

- **Reading B, protobuf as the type system's semantic model.** Dies on the table above. Adopting proto3's expressiveness ceiling means giving up sum types, generics, closed enums, and non-nullability, in exchange for a wire format. That is trading the entire type system for a serializer.
- **Declaration-order field numbers.** Dies against the canonical-formatter goal: a reordering formatter would silently break wire compatibility, and the two features cannot coexist.
- **Emission that can fail the build.** Dies per the toggle principle: an output flag that changes which programs are legal is not an output flag.
- **Hardcoding protobuf as the only schema backend.** Not fully dead, but suspect. The interesting artifact is the walkable type graph, and protobuf is one consumer of it. Committing the compiler to one wire format before choosing a paradigm is backwards.

## Open tensions I did not resolve

- Whether wire numbering should be its own lockfile or a projection of the stable per-definition identity that \[[version-control-integration]\] wants for blame purposes. Two identity schemes for the same declarations would be a smell.
- Whether decoded sum types get an implicit unknown-variant arm (protobuf's forward compatibility preserved, exhaustive matching weakened) or a fallible decode (matching preserved, rolling deploys hurt).
- Whether "the compiler knows two definitions of breaking change" is a feature worth its conceptual cost, or a sign the emission tool belongs outside the compiler.
- Whether the wire surface is a per-type attribute or a project-level manifest list, which is really a question about whether "what we expose over the network" is a type property or an architecture property.
- Whether generic types are monomorphized on emission (name explosion, `OptionString`, `ListOfOrderLineItem`) or simply excluded from the wire surface entirely.

## Prior art to actually read

- protobuf-net (C#): types-as-schema with explicit member numbers, and the accumulated pain of that choice over a decade.
- `buf` and its breaking-change detector: the closest thing to a formal definition of wire compatibility, and a ready-made spec for what the compiler would need to know.
- Cap'n Proto's schema evolution rules versus protobuf's: same problem, different answers on ordinals and defaults, and Kenton Varda wrote both.
- Smithy and TypeSpec: IDLs designed to emit *multiple* target schemas, which is the generalized version of this idea.
- Rust `serde` versus `prost`: format-generic derive versus schema-first codegen, as the two poles this note keeps oscillating between.
- Go's struct tags: the low-tech answer to per-field wire metadata, and why the ecosystem tolerated stringly-typed tags for fifteen years.
