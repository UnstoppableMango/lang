# Feature design workflow

This document defines how a feature of the language moves from a one-sentence idea to formally specified behavior.
It is the process document for the whole project.
The [wishlist](wishlist.md) is the intake for this process.

The guiding principle is deliberate slowness.
The feature set grows granularly, point by point, over many sessions.
Nothing advances quickly, nothing skips stages, and every advancement leaves a written trail.

## Overview

Every feature has a stage, modeled after the TC39 process.

| Stage | Name         | Lives in                        | Means                                        |
| ----- | ------------ | ------------------------------- | -------------------------------------------- |
| 0     | Wishlist     | `docs/wishlist.md`              | An idea worth remembering.                   |
| 1     | Sketch       | `docs/design/NNNN-slug.md`      | A problem worth solving, shape still rough.  |
| 2     | Draft design | same file                       | A concrete design with resolved semantics.   |
| 3     | Accepted     | same file                       | Design frozen, awaiting specification.       |
| 4     | Specified    | `docs/spec/` (future)           | Formal spec text exists, design doc is historical record. |

Movement is one stage at a time, in either direction.
A feature can also be rejected or parked at any stage.

## Stage 0: Wishlist

**Artifact:** a single bullet in `docs/wishlist.md`, under the appropriate area section.

**Entry criteria:** none.
Ideas are cheap and the wishlist exists to capture them without ceremony.
Add the sentence, optionally a minimal example and prior art, and move on.

**Rules:**

- One sentence per feature.
- A feature may include a minimal example.
- A feature may reference prior art from other languages, books, videos, or papers.
- Wishlist entries carry no commitment.
  Listing a feature does not mean the language will have it, and contradictory entries may coexist at stage 0.

**Advancement criteria (to stage 1):**

- Someone can articulate the problem the feature solves, not just the feature itself.
- The feature is interesting enough to spend a session sketching.

## Stage 1: Sketch

**Artifact:** a new design doc at `docs/design/NNNN-slug.md` (see [Design docs](#design-docs)).

**Purpose:** establish that the problem is real and survey the solution space before committing to a shape.

**Required content:**

- Motivation: the problem being solved, with at least one concrete scenario the language should handle well.
- Rough shape: one or more candidate directions, sketched with small examples.
- Prior art: how other languages, papers, or systems approach the same problem, and what to learn from each.
- Open questions: an explicit list of what must be resolved before a real design exists.

**Advancement criteria (to stage 2):**

- The motivation holds up when written down.
- One candidate direction has been chosen and the rejected directions are recorded with reasons.
- The open questions list is complete enough that answering it would fully determine the design.

## Stage 2: Draft design

**Artifact:** the same design doc, expanded.

**Purpose:** produce a design precise enough that two people implementing it independently would build the same thing.

**Required content, in addition to the stage 1 content:**

- Syntax: the concrete surface syntax, if the feature has one, with a grammar sketch.
- Semantics: what the feature means, described precisely in prose (formal notation optional at this stage).
- Worked examples: at least three, covering the common case, an edge case, and a misuse that should be rejected with a good error.
- Interactions: an explicit subsection for every other feature currently at stage 2 or higher, stating how the two features compose or conflict.
- Alternatives considered: the designs not chosen, and why.

**Advancement criteria (to stage 3):**

- Every open question from stage 1 has a written answer.
- Every interaction subsection concludes with either "composes cleanly" and a demonstration, or a resolved conflict.
  An unresolved conflict blocks advancement for both features.
- The misuse example produces a describable, helpful error.

## Stage 3: Accepted

**Artifact:** the same design doc, frontmatter updated to stage 3.

**Meaning:** the design is frozen.
The language has committed to this feature in this shape.

**Rules:**

- No changes except amendments discovered during specification or implementation.
- Every amendment is recorded in a Changelog section at the bottom of the design doc: date, what changed, why.
- An amendment that changes semantics rather than clarifying them demotes the feature back to stage 2.

**Advancement criteria (to stage 4):**

- Formal specification text for the feature has been written and merged into `docs/spec/`.

## Stage 4: Specified

**Meaning:** the spec is now the source of truth for this feature.
The design doc remains in `docs/design/` permanently as the historical record of how the decision was reached, but disagreements between it and the spec resolve in favor of the spec.

`docs/spec/` does not exist yet.
Its structure will be designed when the first feature approaches stage 3, and that structure will itself go through this workflow as a decision record.

## Design docs

### Naming

Design docs live at `docs/design/NNNN-slug.md`.
`NNNN` is a zero-padded number assigned in creation order, starting at `0001`, and never reused.
`slug` is a short kebab-case name for the feature.

### Frontmatter

Every design doc starts with YAML frontmatter:

```yaml
---
title: Pattern matching
stage: 1
status: active
created: 2026-08-24
updated: 2026-08-24
---
```

- `stage`: 1 through 4.
  Stage 0 features have no design doc, so 0 never appears here.
- `status`: `active`, `parked`, or `rejected`.
- `created` and `updated`: dates, updated on any meaningful edit.

### Template

Copy this to start a new stage 1 doc, and delete the stage 2 sections until the feature reaches stage 2:

```markdown
---
title: <feature name>
stage: 1
status: active
created: <date>
updated: <date>
---

# <feature name>

## Motivation

## Rough shape

## Prior art

## Open questions

<!-- Stage 2 sections below. -->

## Syntax

## Semantics

## Worked examples

## Interactions

## Alternatives considered

## Changelog
```

## Decision records

Foundational choices that are not features still flow through this pipeline: compiled versus interpreted, paradigm, memory management strategy, implementation language, and similar anchors.
These are decision records.

- They use the same `docs/design/NNNN-slug.md` numbering and frontmatter, with `title` prefixed `Decision:`.
- They follow the same stages, but stage 4 for a decision means the decision is reflected in the spec's front matter or introduction rather than in a feature clause.
- Until a relevant decision record reaches stage 3, feature designs must not presume its outcome.
  A feature that cannot be designed without presuming an open decision is blocked on that decision, and the block is listed in its open questions.

## Rejection, demotion, and parking

- **Rejection** can happen at any stage.
  Set `status: rejected`, add a final Changelog entry explaining why, and never delete the file.
  Rejected wishlist entries move to a Rejected section at the bottom of the wishlist with a one-sentence reason.
- **Demotion** moves a feature back one stage when its assumptions break, most commonly when a new feature's interaction analysis surfaces a conflict.
  Record the demotion and reason in the Changelog.
- **Parking** (`status: parked`) means the feature is neither advancing nor rejected.
  Parked features are skipped in interaction analyses until reactivated, and reactivation at stage 2 or higher requires redoing the interactions section.

## Session cadence

This encodes "take it slow" as process.

- A session advances a feature by at most one stage.
  Sitting at a stage across multiple sessions is normal and expected.
- Prefer depth over breadth: a session that moves one feature meaningfully beats a session that touches five.
- Adding stage 0 wishlist entries is exempt from the pacing rules and welcome in any session.
- Advancement decisions are made by the language author.
  Agents and collaborators can draft any artifact, but a stage change is deliberate and human-approved.

## Current state

There is no hand-maintained index, because hand-maintained indexes rot.
The source of truth is:

- Stage 0: the contents of `docs/wishlist.md`.
- Stages 1 through 4: `ls docs/design/` plus each file's frontmatter.

To see the state of the project, list the design directory and read the frontmatter.

## Future documents

Planned but deliberately not created yet:

- `docs/spec/`: the formal language specification, begun when the first feature nears stage 3.
- Deeper area documents (tutorials, rationale essays) as the design matures.
