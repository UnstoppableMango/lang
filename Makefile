build:
	nix build .#

update:
	nix flake update

check lint:
	nix flake check

test:
	nix develop -c cargo test

# Rewrite the golden files. Read the diff afterwards; blessing output you have
# not looked at is how a test suite starts agreeing with a bug.
bless:
	nix develop -c env UPDATE_GOLDENS=1 cargo test

format fmt:
	nix fmt

.PHONY: build update check lint test bless format fmt
