# Lang

A pet project for playing with language design and compilers.

## Building

Go and LLVM are required to be pre-installed.

- [Go](https://go.dev/doc/install) - Toolchain for the version listed in [go.mod](./go.mod)
- [LLVM](https://releases.llvm.org/) - 19.1.0 install

### Optional

These dependencies will be installed to [bin](./bin) automatically by `make`.

- [.NET](https://dotnet.microsoft.com/en-us/download) - dotnet SDK for the version listed in [global.json](./global.json)
- [dprint](https://dprint.dev/install/) - binary for the version listed in [.versions/dprint](./.versions/dprint)
- [buf](https://buf.build/docs/installation/) - binary for the version listed in [.versions/buf](./.versions/buf)

### Build

The default goal is `build`, which will produce `bin/ir` and `bin/lang-host`.

```shell
$ make build
// ...
Build succeeded
```

### Test

Running `make test` will test all modified code.

```shell
$ make test
// ...
Test Suite Passed
```

## Stack

Summary of the tech used to build the language.

### Languages

- F# - Parser
- C# - Language host
- Go - IR Codegen

### Libraries/Tools

- FParsec - Parser/Lexer
- LLVM - Compiler
- tiny-go - Go LLVM bindings
- gRPC - IPC .NET <-> Go
