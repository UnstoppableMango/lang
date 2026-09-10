---
name: sketch
description: Promote a wishlist idea to a stage 1 sketch by allocating the next design doc number, scaffolding the template, and researching prior art. Use when the user invokes /sketch, or asks to sketch, promote, or start a design doc for a feature or decision.
argument-hint: <slug or feature phrase>
---

# Sketch: gate 0 → 1

Create a stage 1 design doc for one feature or decision record, per `docs/workflow.md`.
That document is the source of truth; if it disagrees with this skill, the workflow wins.
Invoking this skill is the author's gate 0 → 1 decision for the named feature, so do the work without re-asking permission, but present the finished sketch for review.

## Steps

1. Read `docs/workflow.md` (stage 1 requirements, gate checklist, template) and `docs/wishlist.md`.
2. Resolve the argument to a wishlist entry.
   If no entry matches, add a stage 0 entry to the appropriate wishlist section first (one sentence, optional example, optional prior art), then proceed.
   If the argument is ambiguous between entries, ask which one.
3. Allocate the doc number: the highest `NNNN` in `docs/design/` plus one, zero-padded to four digits, starting at `0001`.
   Create `docs/design/NNNN-slug.md` from the template embedded in `docs/workflow.md`, keeping the stage 1 sections plus Advancement record and Changelog, and deleting the stage 2 design sections.
   For a foundational decision rather than a feature, prefix the title `Decision:`.
4. Research prior art with WebSearch, WebFetch, or deepwiki (e.g. `rust-lang/rfcs`, `WebAssembly/proposals`, academic papers).
   The gate needs at least two cited sources, or a recorded empty search: where you looked, plus the closest known art and how this idea differs.
5. Draft the stage 1 sections:
   - Motivation: the problem, with at least one concrete scenario the language should handle well.
   - Rough shape: one or more candidate directions, each with a small example.
   - Open questions: non-empty, and complete enough that answering them would determine the design.
6. Check open decision records: scan `docs/design/` frontmatter for decisions not yet at stage 3.
   The sketch must not presume their outcomes; anything blocked on one goes in Open questions as an explicit block.
7. Verify the gate 0 → 1 checklist from `docs/workflow.md` item by item.
   Set frontmatter `stage: 1`, `status: active`, `created` and `updated` to today.
   Add one Advancement record line: date, `0 → 1`, one sentence.

## Rules

- Each sentence on its own line; never use the em dash character (U+2014).
- Multiple competing sketches for the same problem may coexist at stage 1; do not reject or edit a rival sketch.
- Do not advance the new doc further; one stage per session, and gate 1 → 2 belongs to a later session via `/advance`.

## Report

End with: the doc path, the wishlist entry it came from, the open questions list, and a one-line reminder that the next advancement happens in a later session.
