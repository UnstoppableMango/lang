---
name: play
description: Open the docs/notes/ playground to riff on a language idea with zero process, below stage 0. Use when the user invokes /play, or asks to play with, doodle, explore, brainstorm, or take notes on an idea that is not ready for the wishlist or a sketch.
argument-hint: <topic or existing note>
---

# Play: the notes playground

Explore an idea in `docs/notes/`, the scratch space below stage 0 defined in `docs/workflow.md`.
This is the one part of the process with no gates, so the skill's job is mostly to protect the freedom: think out loud on paper and resist the urge to formalize.

## Steps

1. Resolve the argument: if it names an existing file in `docs/notes/`, continue in that file; otherwise create `docs/notes/<slug>.md` with a short kebab-case slug.
   Create the directory if it does not exist.
2. Explore generatively.
   Doodle candidate syntaxes side by side, argue both sides of a tradeoff, chase "what if" tangents, write throwaway example programs in hypothetical syntax.
   Hypotheticals may freely presume outcomes of open decision records ("in a world where the language is interpreted...") and may contradict the wishlist, other notes, or designs at any stage.
3. Capture, do not conclude.
   Dead ends are worth keeping with a line on why they died; a note that ends in open tension is a good note.

## Rules

- No frontmatter, no numbering, no required sections, no stage; use whatever structure helps thinking.
- Notes carry zero standing and nothing advances from them: never create or edit a design doc from this skill, and never treat note content as a decision.
- Exempt from session cadence; playing never counts against the one-stage-per-session rule.
- Notes are deletable at will; prune or delete on request without ceremony.
- Each sentence on its own line; never use the em dash character (U+2014).

## Exit

If a distinct idea crystallized, offer one optional exit: distill it to a single sentence and add it to `docs/wishlist.md` as a stage 0 entry citing the note.
Only do this if the user says yes; otherwise the note just stays a note.
