# strong version control integration

What would it mean for a language to take version control seriously as a first-class concern, instead of treating it as something that happens to text files after the fact?

## the standard story and why it's annoying

Git diffs text.
The language emits text.
Every tool downstream (diff, blame, merge, review) operates on lines of text that happen to also be a parse tree, but the tools don't know that.
This causes familiar pain:

- Reformatting a file (reindent, reflow) creates a diff noise storm even though semantics didn't change.
- Renaming a widely used identifier touches every call site as a text diff, burying the one meaningful rename inside a thousand-line diff.
- Two developers add a new field to a struct in the same place; git sees a textual conflict even though the merge is semantically trivial (union of two additions).
- `git blame` tells you which commit last touched a *line*, not which commit last touched a *node* if intervening commits reformatted or reflowed around it.

If the language controlled its own tree structure and had opinions about identity, a lot of this could dissolve.

## idea: content-addressed / structural storage (Unison-style)

Unison's whole pitch is: definitions are hashed by their structure, not named by the programmer at storage time.
Renaming is a view-layer operation, not an edit to the canonical artifact.
This makes "renaming" a non-event for version control, because the underlying hash never changed.

What if source-of-truth isn't a text file at all, but an AST/IR that gets *rendered* into text for human editing and diffing, with the git-tracked artifact being closer to that IR?

Tension: this fights git itself, which is fundamentally a content-addressed store of *blobs*, usually text blobs.
You'd be building a second content-addressed system on top of / instead of git's.
Unison mostly punts on git entirely and rolls its own codebase manager with its own sync protocol.
That is a big commitment: opt out of the entire git tooling ecosystem (GitHub review UI, blame, bisect, hosting) in exchange for perfect semantic diffs.
Is that worth it?
Depends how much weight "strong version control integration" is supposed to carry versus "integrates with the version control everyone already has and uses."

## idea: stable node identity without leaving git

Middle path: keep text as the source of truth (so git still works, GitHub PRs still render, `git blame` still resolves lines), but give every syntactic node (function, struct field, statement) a stable identity that survives reformatting and most refactors.
Think: each top-level definition carries an invisible/derived ID (content hash of its own subtree, independent of surrounding whitespace/order), and tooling (blame, diff, merge) is aware of these IDs even though git itself is not.

This means:

- `lang blame` could be a real command distinct from `git blame`, walking history at the definition level: "this function's *body* hasn't changed since commit X, even though the file around it has been reformatted three times since."
- `lang diff` renders a semantic diff (added field, renamed parameter, changed control flow) as a layer on top of the textual git diff, the way `difftastic`/`diff-so-fancy` do today but with actual compiler-level understanding instead of heuristic tree-sitter parsing.
- Merge conflicts could be pre-resolved at the AST level for the "two people added different things to the same list/struct/enum" case, which is extremely common and extremely mechanical, before falling back to textual 3-way merge for everything else.

None of this requires abandoning git.
It requires the language's own tooling to sit as a translation layer that git is unaware of, similar to how `git diff` can be configured with custom textconv/diff drivers per file type.
That's actually a real hook point: git already supports pluggable diff/merge drivers per gitattributes pattern.
A language could ship one on day one.

## idea: formatting is not a source of diffs, ever

If the canonical stored form is already fully normalized (there is exactly one valid way to print any given AST, no configurability, gofmt taken to its logical extreme), then "reformatting" as a diff-generating event stops being possible altogether, because there's nothing to configure and no drift to correct.
This is weaker than full structural storage but almost as effective for the "noisy diff" problem, and it's much more compatible with staying inside git's normal text model.
gofmt already mostly achieves this in practice: gofmt diffs are rare because there's nowhere for two gofmt'd files to disagree.

Combine with: the formatter is deterministic enough that a pre-commit hook enforcing "file on disk == canonical print of its own AST" is not just a style nicety but an invariant the language can lean on elsewhere (e.g. the blame/diff tooling above can assume no non-semantic whitespace drift ever needs to be filtered out).

## idea: commits as compilation units / changesets with meaning

What if the language's module system understood the *version control history* directly, e.g. a package version isn't a manually bumped semver string but is derived from the actual set of merged changes since the last release, and the language can tell you "this change added a field to a public struct, that's a minor bump" or "this removed an exported function, that's a breaking change" automatically because it can see both the AST diff and the visibility/API surface?

This is sort of what `cargo semver-checks` and similar tools do after the fact, bolted onto an existing language.
Baking it in as native tooling instead of an ecosystem add-on is a "what if this were designed in from the start" question worth keeping.

Dead end to note: this only works if the language has a strong, checkable notion of "public API surface," which presumes decisions (module system, visibility, exports) that haven't been made and can't be presumed here.
Fine as a note, not fine as a design commitment yet.

## idea: language server / VCS as the same conversation

Editors already have LSP giving live diagnostics/hover/refs.
Git already has blame/log giving historical provenance.
These are currently two totally separate subsystems a developer mentally merges themselves ("who wrote this function, and does it still typecheck").
What if hovering a symbol could show, in one tooltip, both its type and a one-line "last changed in commit X, message Y" without shelling out to a separate blame invocation, because the compiler's own symbol table carries a provenance pointer that was populated by asking git once at a coarse grain and cached?

Not a language design question really, more a tooling/IDE integration question, but the note is: version control integration doesn't have to live in the compiler.
It can live in the language server, which already has the tree and already talks to the editor.
The "language" part of "strong VCS integration" might really be "the language ships a language server that treats VCS as a peer data source," not "the language redefines what a commit is."

## tension: how much of this is actually a language decision vs a tooling decision

Running theme through this whole note: almost every idea above is separable into (a) things a language server / CLI tool could do to *any* sufficiently structured language without changing the language itself, and (b) things that require the language itself to guarantee some invariant (canonical formatting, stable node identity, content-addressed definitions) that make (a) tractable or free.

So maybe the real stage-0 question isn't "how does the language integrate with git" but narrower: does the language guarantee a canonical, unique textual representation for any given AST?
That one guarantee (single-format-only, no configurable style) is cheap to commit to early and unlocks most of the diff/blame/merge wins above for free, without committing to anything as radical as ditching git or building a Unison-style codebase manager.
The content-addressed / no-file-source-of-truth stuff stays interesting but genuinely competes with "just use git like every other language," and that's a much bigger bet to weigh later, not now.

## loose thread: what does "merge" mean for a language with structural editing tools

If the language ever grows a structural editor (edit the tree directly, not text, à la some Lisp/Unison tooling), does "version control" even operate on text anymore, or does it operate on tree edits/patches directly (a la CRDTs for structured data)?
Didn't chase this far, but it rhymes with the content-addressed idea above and might be the same underlying question wearing a different hat: is the artifact under version control text, or is it a tree, and is git even the right substrate once you've decided that.
