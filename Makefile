OPAM ?= opam
DUNE ?= dune

build:
	$(DUNE) build

watch:
	$(DUNE) build -w

_opam:
	$(OPAM) switch create .
