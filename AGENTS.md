# AGENTS.md

This file provides guidance to coding agents when working with code in this repository.

## What this repository is

A programming language being designed from scratch, deliberately slowly.
There is no spec yet; the design happens in GitHub issues, `CONTEXT.md`, and ADRs in `docs/adr/`.
The compiler is implemented in Rust, using `inkwell` for LLVM IR generation and `nom` for parsing; this is a starting point, not a locked-in decision, and may change if a switch proves worthwhile.
Every other foundational decision (paradigm, compilation model, memory strategy, concurrency model, pattern matching semantics) is still open, each a `ready-for-human` issue, so no document or code may presume one until its issue is closed by an ADR.
Decisions are made by the language author, never by an agent; agents may draft an ADR, but the author accepts it.

## Layout

- `src/` — the Rust compiler (`unmangc`): a library (`parse`, `codegen`, `diag`, `ast`) with a thin binary around it, plus `tests/` holding the golden suite.
- `hack/` — `hello.lang` plus a `Makefile` that builds and runs it through the compiler, the one working end-to-end example.
- `README.md` — what the language is for and the principles every design choice is held to.
- `CONTEXT.md` — the glossary: one canonical word per concept, and the words to avoid.
- `docs/influences.md` — the author's Likes and Dislikes from other languages, which inform design choices but never bind them.
- `docs/agents/` — how the engineering skills use the issue tracker, triage labels, and domain docs.
- `nix/feature-flags.nix` — toggles which `src/features/<name>` directories are built into the compiler, so an unfinished feature can land behind a flag.

## Commands

Nix drives everything (a direnv devshell provides `gnumake` and `nixfmt`):

- `command make check` (or `nix flake check`): lint/check, same as CI.
- `command make fmt` (or `nix fmt`): format via treefmt (nixfmt for .nix files).
- `command make build` (or `nix build .#`): build, same as CI.
- `hack/Makefile` (`make run` inside `hack/`) compiles and runs `hack/hello.lang` through the built compiler; `make check` there diffs its output against `hack/hello.expected`, the same thing the `e2e` flake check does.
- `command make test`: run the test suite (`cargo test` in the devshell), which is golden tests plus unit tests.
- `command make bless`: rewrite the golden files under `tests/cases/` after an intended change, then read the diff.

`nix build .#` runs `cargo test` as part of the build, so the golden suite gates CI without extra wiring.

## Writing conventions

- In markdown, each sentence goes on its own line; don't hard-wrap.

## Agent skills

### Issue tracker

Issues are tracked in this repository's GitHub Issues. See `docs/agents/issue-tracker.md`.

### Triage labels

The five default triage labels, each named for its role. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one root `CONTEXT.md` and `docs/adr/`, whose records start at ADR-0001. See `docs/agents/domain.md`.
