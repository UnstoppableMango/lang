run: hello
	./hello

hello: hello.ll
	clang -Wno-override-module hello.ll -o hello

hello.ll: hello.lang ../Cargo.toml ../src/main.rs
	cargo run --quiet --manifest-path ../Cargo.toml -- hello.lang > hello.ll

clean:
	cargo clean --manifest-path ../Cargo.toml
	rm -f hello hello.ll

.PHONY: run clean
