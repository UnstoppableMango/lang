run: hello
	./hello

hello: hello.ll
	clang -Wno-override-module hello.ll -o hello

hello.ll: unmangc hello.lang
	./unmangc hello.lang > hello.ll

unmangc: unmangc.c
	clang -Wall -o unmangc unmangc.c

clean:
	rm -f unmangc hello hello.ll

.PHONY: run clean
