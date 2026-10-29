# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root, or
- **`CONTEXT-MAP.md`** at the repo root if it exists: it points at one `CONTEXT.md` per context.
  Read each one relevant to the topic.
- **`docs/adr/`**: read ADRs that touch the area you're about to work in.
  In multi-context repos, also check `src/<context>/docs/adr/` for context-scoped decisions.

If any of these files don't exist, **proceed silently**.
Don't flag their absence; don't suggest creating them upfront.
The `/domain-modeling` skill (reached via `/grill-with-docs` and `/improve-codebase-architecture`) creates them lazily when terms or decisions actually get resolved.

## New work

- A feature: a spec issue (`/to-spec`), broken into ticket issues with blocked-by links (`/to-tickets`).
- A decision that is hard to reverse: an ADR in `docs/adr/` (`/domain-modeling`), numbered from ADR-0001.
- A term: `CONTEXT.md`.

Cite an ADR as `ADR-0007`, so it never reads as an issue number.
Write a plan as tickets, never as a plan file.

## The legacy design workflow

This repository is migrating away from the staged workflow in `docs/workflow.md`: `docs/wishlist.md`, the stage 1-3 docs in `docs/design/`, and the scratch notes in `docs/notes/`.
Until those are pruned and converted into the files above, they still hold the project's thinking, so read the ones that touch your area.

- Never add to them: no new wishlist entry, design doc, decision record, or note.
  New work goes where [New work](#new-work) says.
- A decision record (a `docs/design/` doc titled `Decision: ...`) at stage 3 or later is settled exactly as an ADR is.
- Anything below stage 3, and everything in `docs/notes/`, is still open: no output may presume it.
- A `docs/notes/<slug>-interview-answers.md` holds the author's answers to a record's open questions; read it before asking them again.

When one of those open questions is settled, record the outcome as a new ADR rather than advancing the old record.

## File structure

Single-context repo (most repos):

```text
/
├── CONTEXT.md
├── docs/adr/
│   ├── 0001-event-sourced-orders.md
│   └── 0002-postgres-for-write-model.md
└── src/
```

Multi-context repo (presence of `CONTEXT-MAP.md` at the root):

```text
/
├── CONTEXT-MAP.md
├── docs/adr/                          ← system-wide decisions
└── src/
    ├── ordering/
    │   ├── CONTEXT.md
    │   └── docs/adr/                  ← context-specific decisions
    └── billing/
        ├── CONTEXT.md
        └── docs/adr/
```

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `CONTEXT.md`.
Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, that's a signal: either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/domain-modeling`).

## Flag ADR conflicts

If your output contradicts an existing ADR, or a legacy decision record at stage 3 or later, surface it explicitly rather than silently overriding:

> _Contradicts ADR-0007 (event-sourced orders), but worth reopening because…_
