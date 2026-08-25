build:
	nix build .#

update:
	nix flake update

check lint:
	nix flake check

format fmt:
	nix fmt

.PHONY: build update check lint format fmt
