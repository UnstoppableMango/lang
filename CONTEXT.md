# lang

A general-purpose programming language being designed from scratch, and its reference compiler, `unmangc`.
This glossary fixes the words the design uses; it records no decision, and a term defined here may name something still undecided.

## Process

**Decision**:
A foundational choice about the language, tracked as an issue until an ADR records its outcome.
_Avoid_: Decision record, stage

**Candidate**:
One of the alternatives a decision weighs.
_Avoid_: Option, proposal

**Tension**:
A conflict between two answers the author has given, which a decision must resolve before it can close.

## Compilation

**Compiled program**:
What the compiler produces from a user's source.
_Avoid_: Artifact, executable

**Runtime**:
Machinery that executes alongside the compiled program and ships as part of it, such as a scheduler or a reference-count check.
_Avoid_: VM

**Separate runtime**:
A binary or environment that must be installed on a host apart from the compiled program; "no runtime" in the author's answers means no separate runtime.
_Avoid_: Bare "runtime" for this meaning

**Comptime**:
Code the compiler executes while building, whose results are fixed into the compiled program.
_Avoid_: Const eval, macros

**Comptime interpreter**:
The engine inside the compiler that executes comptime code; it is never part of the compiled program.
_Avoid_: Runtime, VM

## Memory

**Regime**:
A strategy by which the language reclaims a value's memory, such as ownership checked statically or compiler-inserted reference counting.
_Avoid_: Memory model, mode

**Release action**:
What a type declares to run to free a non-memory resource.
_Avoid_: Destructor, finalizer, drop

**Pinning**:
Giving a value a stable address, which only a boundary that requires one, such as the C FFI, does.

## Concurrency

**Task**:
A unit of concurrent execution started with `spawn`.
_Avoid_: Thread, goroutine, coroutine, fiber

**Scope**:
The region a task cannot outlive; when it ends, every task spawned in it has finished.
_Avoid_: Nursery, task group

## Types and matching

**Closed sum type**:
A sum type whose variants are all known where it is defined, so a match over it can be checked for exhaustiveness.
_Avoid_: Sealed, final

**Open sum type**:
A sum type marked `open` at its definition, which consumers must match with a wildcard arm.
_Avoid_: Non-exhaustive, extensible

**Method set**:
The functions defined in a type's module that take the type as their first parameter.
_Avoid_: Methods, impl block
