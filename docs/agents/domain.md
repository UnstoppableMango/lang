# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

This repository is a language being designed through the staged process in `docs/workflow.md`, which is the source of truth and wins over this file on any disagreement.
Read it before touching anything in `docs/`.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root, if it exists: the glossary.
- **`docs/design/`**: the design docs and decision records (`Decision: ...`) that touch the area you're about to work in, with each file's `stage` frontmatter.
- **`docs/wishlist.md`**: stage 0, one sentence per feature, carrying no commitment.
- **`docs/notes/<slug>-interview-answers.md`**: the author's answers to a record's open questions; check for one before treating those questions as unanswered.

If any of these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The `/domain-modeling` skill creates `CONTEXT.md` lazily when terms actually get resolved.

## Decision records stand in for ADRs

This repository has no `docs/adr/`, and skills must not create one.
Wherever a skill says "ADR", read "decision record": a `docs/design/NNNN-decision-*.md` file in the same stage pipeline as every feature.
A new decision is sketched through `/sketch`, which allocates the next number; it is never written as an ADR.
Cite a record by its number and title, such as `0002 (Decision: compilation model)`, so it never reads as an issue number.

A record's stage is its standing:

- Below stage 3, a decision is still open, and no output may presume it.
- At stage 3 or 4, treat it as settled exactly as you would an ADR.
- `docs/notes/` has zero standing, so nothing there settles a question.

Stage changes are decided by the language author, never by an agent.
A skill may draft an artifact a gate asks for, but never advances a stage itself.

## File structure

Single-context:

```text
/
├── CONTEXT.md
├── docs/
│   ├── workflow.md
│   ├── wishlist.md
│   ├── design/
│   │   ├── 0001-decision-memory-strategy.md
│   │   └── 0005-pattern-matching.md
│   └── notes/
└── src/
```

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a design doc, a hypothesis, a test name), use the term as defined in `CONTEXT.md`. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, that's a signal: either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/domain-modeling`).

## Flag decision conflicts

If your output contradicts a decision record at stage 3 or later, surface it explicitly rather than silently overriding:

> _Contradicts 0002 (Decision: compilation model), but worth reopening because…_

If it presumes an answer to a record still below stage 3, say so as well, since that answer is not yet the author's.
