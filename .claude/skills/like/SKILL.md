---
name: like
description: Record a Like, a trait of another language the author wants the language to draw on, in docs/influences.md. Use when the user invokes /like.
argument-hint: <the language, the trait, and why>
---

# Like: record an influence the command already names

Record each trait in the argument as a **Like** in `docs/influences.md`.
The command fixes the kind, so never read it from the author's words and never ask for it.
Read `CONTEXT.md` for **Influence** and **Like**.

## Steps

Follow the **Influence** steps in `.claude/skills/idea/SKILL.md`, which own how an influence line is written and landed, with two differences:

- Skip sorting: every trait in the argument is an influence, and every one is a Like.
- Ask only for a missing why, or a missing language when the trait plainly comes from one, never for Like or Dislike.

If a trait already has an entry, say so and link it instead of adding a second one; changing or removing an entry is not this skill's job.

The step is done when every trait in the argument has one pull request, and the reply links each.
