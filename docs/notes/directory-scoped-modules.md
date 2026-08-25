directory scoped module system

the pitch: a directory is a module.
no `module Foo` header, no `package foo` line, no explicit export list at the top of a file.
the filesystem layout *is* the module tree.
this is Go's model (package = directory, one flat namespace of source files), and it's worth pulling apart why it feels good and where it pinches.

## the appeal

no bikeshedding about whether the module name in the file matches the file path.
no `import { x } from '../../../foo/bar'` because the directory nesting *is* the addressing scheme, so tooling can rewrite paths on move without touching import statements (rename a folder, every importer's path updates itself, because the path IS the folder).
new contributors have one rule to learn: "where the file lives is what it's part of."
compare to Java's `package` declarations, which just restate the directory path in a comment-shaped line that the compiler enforces, redundant with the filesystem, and yet a real source of copy-paste bugs when a file gets moved and the package line doesn't follow.
if the language *never has* a package line, that whole bug class doesn't exist.

## where it pinches

**multiple files, one namespace, no signaling.**
in Go this is a known sharp edge: every `.go` file in a directory silently shares one namespace, so two files in the same package can define colliding names with no import, no `use`, nothing at the top of the file to say "these two files are related."
you find out from the compiler, not from reading either file in isolation.
is that a feature (files are cheap, split however you like, they're just "the module, chunked") or a bug (a file should be legible on its own, and "which other files am I implicitly sharing scope with" should not require directory-listing)?

**can a directory span multiple modules?**
if directory == module, strictly, then no, and that's clarifying: `src/parser/lexer.mn` and `src/parser/token.mn` are unconditionally the same module, period.
but every real codebase eventually wants a directory that's "this feature, plus its tests, plus a couple of private helpers that shouldn't leak."
Go's answer is `_test.go` suffix conventions and unexported (lowercase) names, layered on top of the one-directory-one-package rule rather than inside it.
worth asking whether visibility (what's exported) should be a separate axis from module boundary (what directory you're in), rather than reusing directory nesting to also imply privacy the way some designs do (private = nested in an underscore-prefixed subdirectory, etc).

**does nesting imply anything?**
in Go, `foo/bar` is a totally separate package from `foo`, no relationship at all beyond string prefix, you must explicitly import `foo/bar` even from code inside `foo`.
alternative: nesting could imply visibility, like npm's underscore-folder convention or Rust's `mod` privacy tree, where a parent module can see into child internals but siblings can't see each other without an explicit `pub`.
that's a real design fork: is the directory tree just an addressing scheme (Go), or is it also a privacy/capability tree (Rust)?
Rust's version needs an explicit `mod foo;` declaration even though the directory structure also encodes it, which reintroduces exactly the redundant declaration line directory-scoping was trying to kill.
so "directory scoped" pulls toward Go's flatter, less expressive but zero-redundant model, almost by definition, if the point is killing the declaration line.

## doodled syntax, Go-flavored (option A: fully flat, no declarations)

```
src/
  parser/
    lexer.mn
    token.mn
  ast/
    node.mn
```

```
// src/parser/lexer.mn
fn lex(src: Str) -> [Token] { ... }
```

```
// src/parser/token.mn
type Token = { kind: Str, text: Str }
```

```
// src/ast/node.mn
import parser  // path = "src/parser", directory name doubles as import name
fn parse(tokens: [parser.Token]) -> Node { ... }
```

no header in lexer.mn or token.mn declaring they're part of "parser", the directory says so.
import path is literally the directory path, no separate logical module name to keep in sync.
question hanging off this: what does the *root* module's name come from, if there's no top-level directory to read?
probably the project manifest (whatever the equivalent of `go.mod`'s module line is) has to name the root once, and everything under it is directory-relative from there.
that's one declaration, but it's one, at the root, not one per file.

## doodled syntax, visibility as a separate axis (option B)

what if directory scoping only decides *grouping*, and an explicit marker decides *visibility*, so the two concerns don't get conflated the way Rust's `mod` + `pub` ends up doing:

```
// src/parser/token.mn
pub type Token = { kind: Str, text: Str }   // visible outside src/parser
fn normalize(t: Token) -> Token { ... }      // not pub, visible only within src/parser
```

no `mod` line needed anywhere, still zero redundant declarations, but now there's a real distinction between "what's in this module" (directory, implicit) and "what this module exposes" (per-item `pub`, explicit).
this seems like the actual sweet spot: keep Go's zero-ceremony grouping, borrow Rust's per-item visibility instead of Go's capitalization-as-visibility convention.
(capitalization-as-visibility is its own tangent: `Token` vs `token` carrying semantic weight is compact but means renaming for style reasons silently changes an API's visibility, which is a landmine. worth a separate note if this direction firms up.)

## tangent: does this fight the compiled-first plan?

if the module system is "directory = module, resolved at compile time by walking the filesystem," that's dead simple to reconcile with a batch compiler: walk the tree once, resolve everything, done.
gets more interesting against the [[compiled-first-scripting-alt]] idea of a `mango run scratch.mn` no-ceremony single file mode: a lone file with no directory-mates is trivially "a module of one," so single-file scripts fall out of directory scoping for free, no special case needed.
that's a nice consistency point between these two notes, worth remembering if both survive to sketch stage.

## dead end considered

considered requiring an explicit manifest file per directory (`module.toml` or similar) declaring the module's name and exports, one per directory, instead of either bare directory-scoping or per-item `pub`.
died because it reintroduces the exact redundancy problem (directory path vs declared name can drift) that directory-scoping exists to avoid, just moved from a line inside each source file to a sibling file.
the only thing it buys over option B is being able to see a module's full export list without reading every file in the directory, which is a real but smaller win, maybe worth revisiting as tooling (a generated/derived listing) rather than a hand-maintained source file.

## open tension, not resolved here

option A (fully flat, capitalization or nothing for visibility) vs option B (directory groups, `pub` marks visibility per item).
option B feels more likely to be "right," but it's also just Go's grouping model plus Rust's visibility model duct-taped together, and duct-taped-together designs from two languages with very different philosophies elsewhere are exactly the kind of thing that looks fine in a two-file example and gets weird at scale (what happens to `pub` across deeply nested directories, does it mean "visible to the immediate parent" or "visible to the whole crate" or "visible to importers regardless of nesting depth", Rust's answer to this is famously a lot of nuance for something that was supposed to save a declaration line).
