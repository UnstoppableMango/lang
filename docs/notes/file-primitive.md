# File as a language primitive

Riffing on: what if `file` is a primitive in the language, the way Nix has a path type?
Zero standing, nothing here is decided, and everything presumes outcomes of open decision records freely.

## What Nix actually does, and which part is the interesting part

Nix has a path *literal* (`./foo.txt`, `/etc/passwd`, `<nixpkgs>`) that is lexically distinguished from a string.
The lexer sees a slash and produces a different type.
That alone is mildly interesting but not the good part.

The good part is the coercion rule: when a path is used where a string is needed inside a derivation, Nix copies the file into the store and substitutes the store path.
So `"${./config.toml}"` does not mean "the text `./config.toml`", it means "a content-addressed copy of that file exists, and here is where".
The language quietly moved a byte range from the author's machine into the build graph, and the value you hold is a *proof that it got there*.

Three separable ideas hiding in that:

1. Paths are relative to the file that wrote them, not to the process's cwd. Nix resolves `./foo` against the source file's directory.
2. Referencing a file is an act of dependency declaration, not of reading.
3. Identity is content, not location. Two files with the same bytes are the same value.

Any of the three could be stolen without the other two.

## The obvious question: when does the read happen?

Say the language has `file`.
What is `x` here?

```
let x = ./config.toml
```

Candidate answers, and none of them are free:

**(a) x is a path.** A validated, normalized, source-relative name. No IO happened. Reading is a separate operation. This is the most conservative and the most boring, and it is probably right for a first cut.

**(b) x is the contents.** `./config.toml` evaluates to the bytes, at compile time. This is `include_bytes!` / `@embedFile` with nicer syntax. The whole file becomes part of the program. Attractive because it makes the "no IO at runtime" story trivially true, and repulsive because now every path literal is a build input and a 4GB file is a 4GB binary.

**(c) x is a lazy handle.** Contents are a thunk. Forcing it reads. In a lazily evaluated world this is the natural answer, and it is also how you accidentally invent an effect system by the back door, because now forcing a value can fail with ENOENT.

**(d) x is a capability.** `file` is not data at all, it is the right to touch that region of the filesystem. Nothing can read a file it was not handed. This is the most opinionated and the most fun. See below.

I keep wanting (a) at the type level with (b) available as an explicit operator and (d) as the thing the runtime actually deals in, which suggests these are three different types wearing one word.

## Doodle: three types, not one

```
Path      // a name. pure. comparable. no claim that anything exists.
File      // a name plus evidence that it existed. obtained by opening a Path.
Blob      // bytes. content-addressed. no name at all.
```

Then the Nix coercion is just `Path -> Blob` performed at build time, and the store path is the `Blob`'s address rendered as a `Path` again, which is a nice little cycle.

```
let cfg  : Path = ./config.toml
let blob : Blob = embed cfg        // build time, contents baked in
let f    : File = open cfg         // run time, can fail
```

Tension: is `embed` an operator, a function, or a compilation mode?
If it is a function it has to run at compile time, which means the language needs compile-time evaluation, which is a wishlist entry and not a decision.
If it is an operator it is a wart.
If it is a mode (`--embed-all`) it is invisible at the call site and that is worse.

## Doodle: directories

Nix lets you write `./.` and get the whole tree, and this turns out to be both the killer feature and the source of every "why did my build rerun" bug ever filed.

If the language has `file` it will immediately want `dir`.

```
let assets = ./assets/          // trailing slash makes it a Dir?
for f in assets { ... }         // Dir is iterable
let page = assets / "index.html"   // path join as an operator
```

`/` as a join operator is cute and collides with division, unless paths and numbers are distinct enough that overload resolution never gets confused.
Python's `pathlib` proves people like it.
Every language that tried it also proves that `a / b / c` reads as arithmetic on first glance.

Alternative that avoids the collision entirely: make path *literals* extensible, so `./assets/{name}` is a path literal with interpolation, and joining is just writing a longer literal.
Then there is no join operator because there was never a join.
That is more Nix-like than Nix.

## Doodle: the capability reading

This is the one I want to chase.

If `File` can only be constructed by opening a `Path`, and `Path` literals are only written in source, then the set of files a program can touch is statically the set of path literals in its source plus whatever it was handed.
That is a real, checkable property, and it falls out of the type system rather than out of a sandbox.

```
fn main(stdin: File, out: File, cwd: Dir) = ...
```

No ambient filesystem.
No `open("/etc/shadow")` from a library, because the library was never given a `Dir` to resolve against.
Deno and capability-secure languages (E, Monte) went here, and Wasm's preopens are the mainstream version.

Where it falls apart: the moment one function needs to hand a subtree to another, you need attenuation (`cwd / "assets"` yields a Dir that cannot escape upward), which means `..` has to be a type error or a runtime rejection, which means paths are not just strings and normalization is part of the semantics.
That is a real cost and it shows up in the spec, not just the stdlib.

Also: what does a config file do?
Real programs read paths from arguments and config, which are strings at runtime, which means there must be a `Dir -> String -> Result<File>` and now the static story is "the set of *roots* is static", which is weaker but still worth something.

## Doodle: content addressing as the identity

If `Blob` is identified by hash, then:

```
./a.txt == ./b.txt    // true if the bytes match
```

That is a defensible equality and a horrifying one.
Defensible because it is what "the same file" means for reproducible builds.
Horrifying because comparing two paths silently reads two files.

Fix: equality on `Path` is name equality, equality on `Blob` is content equality, and there is no implicit `Path -> Blob`.
Which is exactly the thing that makes Nix nice, removed.
Open tension, unresolved, and I think this is the crux of whether "file as primitive" is worth a decision record at all.

## What "primitive" would actually buy

Being honest about it, because "make it a primitive" is usually a smell:

- Source-relative resolution. This genuinely cannot be a library, because a library function does not know which source file called it. This is the strongest argument, and it may be the *only* one that needs language support.
- Literal syntax that the lexer distinguishes, so no quoting, no escaping, no backslash tragedy on Windows.
- Build-graph participation: the compiler knows which files the program depends on, which feeds incremental builds and reproducibility.
- Capability enforcement, if the type system is doing the work anyway.

Everything else (join, iterate, read, stat) is stdlib and should stay stdlib.

That is a useful narrowing: maybe the primitive is not "file", it is "source-relative path literal", and `File` is an ordinary library type built on it.

## Dead ends worth keeping

- **File as a first-class mutable value** (`x = ./log.txt; x <- "hello"`). Dies because it makes every assignment potentially an IO error and smuggles effects into the most basic syntactic form in the language. Nothing gained over `write(f, s)`.
- **Everything is a file** (Plan 9 in the type system). Dies because it is a runtime/OS idea, not a language idea. The language would just be renaming its handle type.
- **Implicit `Path -> Blob` coercion everywhere**, the full Nix rule. Dies for a general-purpose language because Nix can only afford it by having exactly one consumer (the store) and no runtime. Without a store to copy into, the coercion has no meaning.
- **`/` as path join.** Not dead, but on probation, pending whether path literal interpolation makes it unnecessary.

## Open tensions I did not resolve

- Compile-time read vs runtime read is really a question about whether the language has compile-time evaluation, which is not decided. This note cannot go further without presuming that.
- Content identity vs name identity, per above.
- Whether the capability angle is a file feature or an effects feature that files happen to be a good demo of. Suspicion: the latter, and the note is really about effects.

## Prior art to actually read

- Nix: path type, string coercion, `builtins.path`, why `./.` is a trap.
- Zig `@embedFile`, Rust `include_str!` / `include_bytes!`, Go `//go:embed`. All three are the "(b) contents" answer, all three needed compiler support, and all three are source-relative. That is three independent votes for the narrowing above.
- Deno permissions, Wasm/WASI preopens, capability-secure languages (E, Monte) for the capability angle.
- Python `pathlib`, `.NET` `Path` for the boring ergonomics that any answer still has to get right.
