# insignificant whitespace

Default assumption in most C-family languages: whitespace between tokens is insignificant.
Newlines, indentation, extra spaces, all just token separators.
Worth doodling what's actually being given up by taking that default, and what the alternatives buy.

## the boring case first

```
fn add(a, b) {
    return a + b;
}
```

vs the same thing squashed:

```
fn add(a,b){return a+b;}
```

Both parse identically.
The parser never sees whitespace as data, it's stripped in the lexer or treated as a token separator and discarded.
This is the "expression-oriented, delimiter-driven" family: braces, semicolons, parens all carry the structural weight so whitespace doesn't have to.

## the significant-whitespace case

Python, Haskell (layout rule), F#, Nim, YAML.
Indentation carries block structure instead of `{}`.
Two very different flavors worth separating:

1. **Python style**: indentation IS the block delimiter.
   No brace ever legal for grouping statements.
   Consistent indent within a block enforced by the lexer (INDENT/DEDENT tokens).
1. **Haskell layout rule**: indentation is sugar for explicit braces/semicolons the parser could equally accept.
   You can always fall back to `let { x = 1; y = 2 } in x + y` explicit form.
   The whitespace-sensitive version is desugared to the explicit-delimiter version before or during parsing.

Haskell's approach is interesting because it means whitespace is significant but not *load-bearing* in the grammar, it's a presentation-layer transform on top of a grammar that doesn't actually need it.
That feels like it could get "the best of both": write the pretty layout form, but the grammar itself is defined in terms of delimiters, and there exists a mechanical desugaring.

## what do you actually get from significant whitespace?

- Forces one canonical visual shape for a given logical structure (sort of; Python still allows inconsistent indent widths across files, just not within a block).
- Removes a whole category of "misleading indentation" bugs, where the visual nesting lies about the actual nesting (the classic C `if` without braces gotcha).
- Fewer tokens on the page, arguably more readable at small scale.
- Couples the language spec to how it looks in a text editor, tabs vs spaces becomes a language question instead of a purely aesthetic one.
- Multi-line expressions get awkward: continuation lines need their own rule (Python's backslash-continuation and implicit-continuation-inside-brackets, Haskell's "further indented than the layout column" rule).
- Copy-pasting code between contexts with different ambient indentation breaks it. Braces are copy-paste-safe. This one keeps nagging at me: it's a "the code is not just text, it's text plus position" property, and it interacts badly with any tooling that treats source as line-oriented (diff, patch context, string templates that inject code).
- Auto-generated code (codegen, macros) is easier to emit correctly when delimiters do the grouping. Emitting well-formed indentation-sensitive output requires a pretty-printer, not just string concatenation.

## what does insignificant whitespace lose?

- Nothing forces you to *use* consistent formatting, hence `gofmt`/`prettier`/etc as an entire tooling category that exists to compensate.
- If the project already lands on \[[single-canonical-formatter]\], is significant whitespace actually buying anything beyond what a mandatory formatter already gives you?
  A formatter-enforced language already has one canonical shape per AST; whitespace being "insignificant" to the grammar doesn't mean it's insignificant to the humans, it's just insignificant to the *parser*.
- This reframes the question: significant whitespace looks like it's really about *where the enforcement lives*.
  Python puts block structure enforcement in the grammar itself (can't even parse malformed indentation).
  A gofmt-style language puts equivalent enforcement in a side tool that's a social/tooling convention, not a grammar rule, so you *can* write (and even compile) ugly or misleading code, you just won't get it past code review or CI fmt-check.

## tangent: whitespace significant only in specific positions?

What if it's not all-or-nothing?
Could imagine whitespace being insignificant almost everywhere (expression-level, delimiter-driven) but mattering in one narrow spot, e.g. for something like a shell-command literal or a heredoc/multiline-string, where the whole point is that whitespace is data.
That's less "is the language whitespace-sensitive" and more "does the language have whitespace-sensitive *sublanguages* embedded via literals."
Most languages already do this for string literals implicitly, nobody calls a triple-quoted Python string "significant whitespace" even though the whitespace inside it obviously matters.
So maybe the real axis isn't "significant vs not" but "how large is the whitespace-insensitive core, and what escapes from it are allowed."

## tangent: significant whitespace and macros/metaprogramming

If \[[ast-in-public-api]\] ever goes anywhere and macros operate on syntax trees rather than raw tokens, does whitespace-sensitivity make macro-generated code harder to splice?
A macro that generates a block of statements to insert at some call site, in a Python-like world, needs to know the ambient indentation of the splice point to emit legally-indented code.
In a brace world it just emits `{ ... }` and does not care where it lands.
This feels like a real tax on indentation-sensitive grammars for any kind of code-generation-as-text approach, though it goes away if macros operate purely on ASTs and only render to text at the very end via a pretty-printer (circles back to the Haskell layout-rule idea: define the grammar in terms of delimiters, treat the indentation-sensitive form as one particular rendering).

## where this is leaning, unresolved

Kind of talked myself toward "the grammar should probably be delimiter-driven underneath (for macro/codegen/copy-paste sanity), with an optional Haskell-style layout rule as sugar if a light-indentation-based surface syntax is wanted later."
But that's a real design commitment with implications for the parser architecture, and it's presuming outcomes on macros/AST-in-public-API that are nowhere near decided.
Also haven't touched: does significant whitespace interact with \[[functional-implicit-return]\] or expression-orientation generally?
Expression-oriented languages with implicit returns (last expression in a block is the value) seem to lean toward needing very clear block delimiters precisely because "what's the last expression of this block" is a question braces answer trivially and indentation answers more ambiguously (what's the last statement at this indent level, going by feel, not rule).
Open question, not chased down yet.
