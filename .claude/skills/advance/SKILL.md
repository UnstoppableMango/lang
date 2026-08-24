---
name: advance
description: Check a design doc's artifacts-to-advance checklist and move it up one stage only if every artifact exists. Use when the user invokes /advance, or asks to advance, promote, or accept a design doc or check whether one can advance.
argument-hint: <NNNN or slug>
---

# Advance: verify a gate, then move one stage

Evaluate one design doc against the next gate's checklist in `docs/workflow.md`, and update its state only on a full pass.
That document is the source of truth; if it disagrees with this skill, the workflow wins.

The cardinal rule: this skill verifies artifacts, it never creates them.
If an artifact is missing, report the gap and stop; writing the missing content is authoring work for a separate session, not part of advancement.

## Steps

1. Resolve the argument to `docs/design/NNNN-slug.md` by number or slug.
   Read it and `docs/workflow.md`.
2. Refuse early, with the reason, if any of these hold:
   - `status` is `parked` or `rejected` (reactivation is manual, per the workflow).
   - `stage` is already 4.
   - The latest Advancement record line is dated today: one stage per session, so this doc cannot move again until a later day.
3. Evaluate the target gate's "Artifacts to advance" checklist item by item.
   For each item, either quote the evidence from the doc (pass) or name exactly what is missing (fail).
4. Gate-specific checks:
   - **To stage 2 or 3**: list every other `docs/design/` doc with `status: active` and `stage >= 2`, and confirm the Interactions section has a subsection for each.
     An interaction subsection passes only if it concludes with "composes cleanly" plus a demonstration, or a resolved conflict.
     An unresolved conflict blocks both features; say so.
   - **To stage 3 (FCP)**: if every other item passes and `fcp-started` is absent, set `fcp-started` to today, update `updated`, and stop: acceptance is allowed no earlier than the next session on a later day.
     If `fcp-started` is present, it must be earlier than today.
   - **To stage 4**: the feature's spec text and its conformance examples must be merged in `docs/spec/`.
     If an implementation exists in the repo, the feature must also be implemented behind a flag with the conformance examples passing as tests; if no implementation exists yet, note that this item is conditionally waived.
5. If any item fails: print the full checklist verdict and change nothing, not even `updated`.
6. If all items pass: increment `stage` by exactly one, set `updated` to today, and append one Advancement record line (date, gate, one sentence).
   If the new stage is 3, note in the report that the design is now frozen and only spec-driven amendments recorded in the Changelog are allowed.

## Rules

- Never skip a stage, never advance more than one stage, never advance a second doc in the same invocation.
- Never edit the doc's content sections; the only writable fields are frontmatter and the Advancement record.
- Each sentence on its own line; never use the em dash character (U+2014).

## Report

End with: the per-item checklist verdict, the resulting stage (or unchanged), and the single next action for the author (what to write next, or when acceptance becomes allowed).
