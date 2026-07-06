# mangolang Wishlist

Unstructured dump of features, syntax ideas, and opinions. No commitment — items get promoted to SPEC.md when settled, dropped when rejected.

---

## Syntax Likes

- **F# pipe operator** `|>` — already in spec, keep it
- **OCaml-style type annotations on the left** — `Int x` not `x: Int`; feels cleaner for function signatures
- **Elm-style `let...in` expressions** inside function bodies
- **Go-style package/module by directory** — no explicit export lists; PascalCase = exported
- **No trailing commas required** in multi-line record/tuple literals
- **`match` as an expression** — value, not a statement
- **String interpolation** with `${}` or `#{}` — leaning toward `#{expr}`
- **Labeled/named arguments** at call site for clarity: `move(from: a, to: b)`
- **Multi-line strings** — triple-quoted or leading-pipe style (F# `"""..."""` or Haskell-ish)
- **Underscore `_` in numeric literals** — `1_000_000`
- **`_` as wildcard** in patterns and unused bindings
- **Block comments** `/* */` in addition to line comments `//`
- **Doc comments** `///` that attach to the next declaration

## Syntax Dislikes

- Semicolons — no
- Braces for block delimiters on their own line (C-style) — prefer Go/Rust inline brace or indentation
- `implements` / `extends` keywords — structural types only
- `class` keyword — use `type` and traits/typeclasses
- Visible lifetime annotations (Rust `'a`) — region-based should be invisible
- `new` keyword for allocation — just call a constructor function
- `void` return type — use `Unit`
- Operator overloading as a free-for-all — restrict to numeric-adjacent ops or ban entirely

## Type System Wants

- **Typeclass/trait constraints** on generics: `fn sort<T : Ord>(List<T>) -> List<T>`
- **Newtype wrappers** with zero runtime cost: `type UserId = UserId Int`
- **Phantom types** for compile-time tagging (e.g. validated vs raw)
- **Row polymorphism** on records (structural subtyping for records) — uncertain, evaluate later
- **`Never` / bottom type** for functions that don't return (diverge/panic/exit)
- **Opaque types** — expose type name but hide constructor from outside module
- **Type aliases** — `type Seconds = Float64` (transparent, for readability)
- **No implicit numeric coercion** — `Int` and `Float64` never silently widen

## Error Handling Wants

- **`?` operator** (Rust-style) for early return on `Err` inside a `Result`-returning function
- **`try/with` expression** as alternative to `match` for error chains — maybe
- **Typed errors** — `Err` variant carries structured type, not raw `String`
- **Error context chaining** — `err |> addContext "while parsing config"`

## Concurrency Wants

- **`go expr` spawn syntax** (Go-style) — familiar, minimal noise
- **Typed channels** `Chan<T>` — send/receive explicit in type
- **`select` expression** over multiple channels
- **Structured concurrency** — goroutines scoped to a region/scope, not fire-and-forget globally
- **No shared mutable state** — enforce at type level or region level

## Memory / Performance Wants

- **Region syntax** to be invisible in 99% of code — compiler infers region from scope
- **Explicit `region` block** as escape hatch for manual arena management
- **`stack` hint** for stack allocation of small structs — or just let compiler decide
- **No GC pauses** — hard requirement; region-based gives this for free

## Module / Package Wants

- **Single namespace per directory** (Go model)
- **Explicit re-exports** when you want to surface internal symbols
- **Qualified imports** `import math` → `math.pi()` OR **selective** `import math (pi, sqrt)`
- **No circular imports** — hard compile error
- **Versioned package manager** built-in (like `go mod`), not an afterthought

## DX Wants

- **Compiler error points to root cause**, not downstream symptom
- **Suggestions in error messages** — "did you mean X?"
- **`mango fmt`** — single canonical formatter, no options
- **`mango check`** — type check without codegen
- **`mango lsp`** — start language server
- **`mango run`** — compile + execute
- **Watch mode** — `mango run --watch`
- **REPL** — `mango repl` for quick experimentation
- **Test runner built-in** — no external framework needed
- **Inline tests** — test functions in the same file, stripped from release build

## Undecided / Open

- Operator overloading: ban entirely, or allow for `Eq`/`Ord`/`Show` typeclasses only?
- Integer types: just `Int` (arbitrary precision) + `Int64`/`Int32`? Or Go-style `int8`..`int64`?
- Float types: just `Float64`, or also `Float32`?
- Tuples: anonymous structs `(Int, String)` or named-only?
- Mutability: does `mut` propagate through a record field, or is each field individually annotated?
- String type: UTF-8 bytes slice under the hood, or rune/codepoint abstraction?
- FFI: C interop via `extern` block (Rust-style) or via Go shim (since Go target ships first)?
- Tail-call optimization: guarantee it for self-tail-calls? Or full TCO?
- Exceptions from C FFI: must catch at FFI boundary, or propagate as `Result`?

---

## Promoted to SPEC.md

- `fn`, `let`, `mut`, `type`, `match`, `module`, `import` keywords
- `camelCase` values, `PascalCase` types
- Type annotation before name: `Int x`
- Pipeline `|>` operator
- ADTs: sum `|` and record `{}`
- No semicolons, tabs for indentation
- `Result<T, E>` and `Option<T>` as stdlib types
- UTF-8 source, `.um` extension
- Directory = module, name matches directory

## Rejected

- Null / nil
- Exceptions
- Global mutable state
- Runtime reflection
- Implicit type coercions
- `implements` / `extends`
- Macros (v1)
- Higher-kinded types (deferred, not rejected)
