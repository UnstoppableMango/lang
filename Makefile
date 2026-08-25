build:
	nix build .#

update:
	nix flake update

check lint:
	nix flake check

format fmt:
	nix fmt

run: hello
	./hello

hello: hello.ll
	clang -Wno-override-module hello.ll -o hello

hello.ll: hello.lang Cargo.toml src/main.rs
	cargo run --quiet -- hello.lang > hello.ll

clean:
	cargo clean
	rm -f hello hello.ll

.PHONY: build update check lint format fmt run clean
