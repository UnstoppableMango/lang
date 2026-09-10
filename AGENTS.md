# AGENTS.md

This file provides guidance to coding agents when working with code in this repository.

## What this repository is

A programming language being designed from scratch, deliberately slowly.
There is no spec yet; most of the repository is a design process and its documents.
The compiler is implemented in Rust, using `inkwell` for LLVM IR generation and `nom` for parsing; this is a starting point, not a locked-in decision, and may change if a switch proves worthwhile.
Every other foundational decision (paradigm, compilation model, memory strategy) is still open, so no document may presume any of them until a decision record reaches stage 3.

## Layout

- `src/` — the Rust compiler (`unmangc`), currently a small stub with no AST or passes.
- `hack/` — `hello.lang` plus a `Makefile` that builds and runs it through the compiler, the one working end-to-end example.
- `docs/` — the design process: `wishlist.md`, `design/`, `notes/`, `workflow.md`.
- `nix/feature-flags.nix` — toggles which `src/features/<name>` directories are built into the compiler; this is the mechanism `docs/workflow.md`'s stage 3 → 4 gate means by "implemented behind a flag."

## Commands

Nix drives everything (a direnv devshell provides `gnumake` and `nixfmt`):

- `command make check` (or `nix flake check`): lint/check, same as CI.
- `command make fmt` (or `nix fmt`): format via treefmt (nixfmt for .nix files).
- `command make build` (or `nix build .#`): build, same as CI.
- `hack/Makefile` (`make run` inside `hack/`) compiles and runs `hack/hello.lang` through the built compiler, the only working end-to-end example in the repo.

There are no tests beyond `nix flake check`.

## The feature design workflow

`docs/workflow.md` is the process document for the whole project and the source of truth; read it before touching anything in `docs/`.
Key points:

- Every feature has a stage 0-4 (modeled on TC39): wishlist bullet → sketch → draft design → accepted → specified.
- Stage 0 lives in `docs/wishlist.md` (one sentence per feature; entries carry no commitment).
- Stages 1-3 live in `docs/design/NNNN-slug.md` with YAML frontmatter (`stage`, `status`, `created`, `updated`); numbers are assigned in creation order and never reused.
- `docs/notes/` is an ungated scratch space below stage 0 with zero standing.
- Each stage gate is a checklist of artifacts that must literally exist in the repo; movement is one stage at a time, and a session advances a feature by at most one stage.
- Stage changes are decided by the language author, never by an agent.
  Agents may draft artifacts, but advancement is human-approved.
- There is no hand-maintained index of features; state is read from `docs/wishlist.md` plus `ls docs/design/` and each file's frontmatter.
- Foundational choices are "decision records" in the same pipeline, titled `Decision: ...`.
- A wishlist entry is exactly one sentence (never semicolon-chained clauses) and states the decision/feature itself, never that a decision needs to be made ("chosen as an explicit decision" is not a decision).
- A decision record's design doc sketch often pairs with `docs/notes/<slug>-interview-answers.md`: the author's answers to its open questions, carrying zero standing until written into the doc at stage 2. Check for one before assuming a record's open questions are unanswered.
- Before re-asking a "new tension" from an interview note, check whether another record's interview-answers note already resolves it by derivation, terms like "no runtime" get reused ambiguously across records and a later clarification in one can settle an earlier tension in another.
- When two records independently reject the same underlying idea under different names (e.g. arenas vs. data-oriented layout), write one shared `docs/notes/` file both point back to instead of duplicating the parallel in each.

Project skills automate the mechanical parts: `/play` (notes playground), `/sketch` (gate 0→1, allocates the next doc number and scaffolds the template), `/advance` (verifies a gate checklist, reports gaps, never writes missing artifacts).
The skills defer to `docs/workflow.md` on any disagreement.

## Writing conventions

- In markdown, each sentence goes on its own line; don't hard-wrap.
- The design doc template (frontmatter, required sections) is in `docs/workflow.md`; copy it exactly when creating a new doc.
