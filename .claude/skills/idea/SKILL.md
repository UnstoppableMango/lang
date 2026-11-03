---
name: idea
description: Record a language idea in the tracker, or an influence in docs/influences.md, without triaging it. Use when the user invokes /idea or explicitly asks to record, capture, or jot down an idea.
argument-hint: <the idea, in your own words>
---

# Idea: record it and get out of the way

Record each idea the author gives you, fast, in the author's own words.
Recording is the whole job: shaping, blocking, and deciding belong to `/triage` and the grillings.
Read `CONTEXT.md` for **Idea** and **Influence**, and `docs/agents/issue-tracker.md` for how to reach GitHub.

## Steps

1. Split the argument into ideas, one record each.
2. Sort each idea into one of three homes, in this order:
   - **Influence**: a trait of another language the author likes or dislikes, such as "I want Rust's `?`".
     It goes to `docs/influences.md` alone, with no idea issue beside it.
   - **Owned idea**: an open issue already owns the question it bears on, typically a decision (`ready-for-human`, #97-#101) or a milestone grilling (`wayfinder:grilling`).
     List open issues and read the likely owners before choosing one.
   - **New idea**: everything else, including an idea whose owner you are unsure of.
3. Write the record:
   - **Influence**: on a branch from `main`, add one line under its topic heading in the file's own form (Like or Dislike, the language, the thing, and why), adding the heading if the topic is new.
     Read Like or Dislike and the language from the author's words.
     Open a pull request holding that one line, so each influence merges on its own.
   - **Owned idea**: comment on the owner.
   - **New idea**: open an issue labelled `needs-triage` and nothing else, with a short title you write.
     When a possible owner exists, add `May belong to #N.` to the body.
   - Every record's body (the issue, the comment, or the influence pull request's description) is the author's words verbatim as a quote, then one `Related:` line linking each decision, grilling, or influence section it touches, or `Related: none`.
4. Reply with one link per record and the home each went to, so a wrong guess is easy to move.

The step is done when every idea in the argument has exactly one record and its link.

## Questions

Ask at most one question per run, and only for what the record cannot exist without:

- what the idea is, when the words leave it unclear;
- an influence's why, when the author gave none;
- Like or Dislike, or the language, when neither can be read off.

Everything else, record as given and let the reply's links carry any doubt.
