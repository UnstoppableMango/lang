spec vs reference compiler

the prompt: "fully specified language, primary compiler serves as a reference implementation"
worth unpacking because it smuggles in an ordering claim.
spec first, compiler second, compiler answerable to spec.
that's one whole tradition.
the other tradition is: compiler first, spec (if it ever exists) is just a description of what the compiler already does.

## camp A: spec is the law

C, Scheme (R7RS), SQL, POSIX shell.
the spec is a document independent of any implementation.
multiple compilers can exist and can disagree with each other, and when they do, the spec adjudicates, not vibes, not "well gcc does it this way."
this is what lets you have gcc AND clang AND msvc and still call it "C."

what it costs:
writing a real spec is its own skill, arguably harder than writing a compiler.
ISO C is unreadable to normal humans.
ambiguity in the spec becomes a standing bug that outlives any single implementation ("undefined behavior" as a genre of horror story).
spec work front-loads a huge amount of design rigor before you get to see whether the language is actually pleasant to use.
you can spec yourself into a corner nobody can implement efficiently (see: some corners of C++ template resolution).

## camp B: the compiler is the spec

Python (for a long time, "what CPython does" WAS the language), Ruby (MRI), Perl ("only perl can parse Perl" is a joke about this exact thing).
here the reference implementation isn't "a" implementation, it's the only source of truth, full stop.
docs are a lossy description of it, not an authority over it.

what it costs:
implementation accidents become load-bearing language semantics.
dict ordering in old Python, or numeric coercion quirks in Perl, things nobody designed on purpose but now can't be removed because someone's code depends on them.
alternate implementations (PyPy, JRuby) are perpetually playing catch-up and perpetually slightly wrong.
"is this a bug or a feature" becomes genuinely undecidable from inside the language.

what it buys:
you get to build the thing and let usage tell you what the spec should have said.
way faster iteration.
zero risk of specifying something unimplementable, because if you can't implement it you never finish the sentence.

## hybrid: spec exists, compiler is "the" reference impl but doesn't have final say

this seems to be what the prompt is actually gesturing at.
"fully specified" = a document exists and is authoritative.
"primary compiler serves as a reference implementation" = there could in principle be others, and the compiler is *a* proof the spec is implementable, not the definition itself.

this is Scheme's model, and it's WG14/C's stated model even though in practice camp A langs drift toward "well what does gcc do" for anything underspecified.
question: does "fully specified" mean spec-complete before compiler exists (waterfall-ish, TC39-adjacent, which is oddly close to what docs/workflow.md already does with stage gates) or does it mean spec and compiler co-evolve and "fully specified" is just the state you eventually reach?

tangent: docs/workflow.md's stage 0-4 pipeline is ALREADY camp A shaped.
sketch -> draft -> accepted -> specified, with the spec doc being the artifact that gates advancement, not "does the compiler do it."
so this whole question might already be answered by the workflow itself and not by mangolang's design docs.
worth noticing: the meta-process (how we design the language) and the object-level design (what mangolang's spec/impl relationship is) are two different decisions and it'd be easy to accidentally let one dictate the other without ever deciding on purpose.

## a third option nobody asked for: executable spec

spec written as something the compiler literally consumes or is tested against (like WebAssembly's spec-tests, or how TC39 proposals ship a reference implementation AND test262).
"fully specified" then means "a conformance suite exists," not "a prose document exists."
compiler passing the suite = compiler is conformant, full stop, no interpretation gap.
this sidesteps ISO-C-unreadability (spec is executable, not prose you have to parse in your head) but doesn't sidestep "did we spec something unimplementable" (you'd find out when writing the compiler against it, same as camp B, just with the discovery happening against a written contract instead of vibes).

open tension, not resolved: is "reference implementation" doing real work in that sentence, or is it decorative?
if there's only ever going to be one compiler, ever, "reference implementation" is just "the implementation" wearing a fancier hat.
the word only means something once a second implementation exists to be checked against the first.
so maybe the real question buried in the prompt isn't spec-vs-compiler at all, it's "do we ever want a second compiler to exist," and everything else follows from that answer.
