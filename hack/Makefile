UNMANGC ?= ../result/bin/unmangc

run: hello
	./hello

hello: hello.ll
	clang -Wno-override-module hello.ll -o hello

hello.ll: hello.lang $(UNMANGC)
	$(UNMANGC) hello.lang > hello.ll

$(UNMANGC):
	$(error no compiler binary at $(UNMANGC); build one (e.g. 'nix build' at the repo root) and/or pass UNMANGC=<path>)

clean:
	rm -f hello hello.ll

.PHONY: run clean
