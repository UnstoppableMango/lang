OPAM ?= opam
DUNE ?= dune

build:
	$(DUNE) build

format fmt:
	$(DUNE) fmt

watch:
	$(DUNE) build -w

_opam:
	$(OPAM) switch create .
