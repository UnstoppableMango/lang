{
  description = "My pet lang so I can tell people I work on stuff like this";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    crane.url = "github:ipetkov/crane";
  };

  outputs =
    inputs@{ flake-parts, crane, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        {
          pkgs,
          lib,
          inputs',
          ...
        }:
        let
          toolchain = inputs'.fenix.packages.stable.toolchain;
          craneLib = (crane.mkLib pkgs).overrideToolchain (_: toolchain);

          # The compiler's Source -> LLVM IR stage links against libLLVM via
          # inkwell/llvm-sys; unrelated to the separate clang linking stage
          # (LLVM IR -> binary) that `make hello` drives.
          llvmEnv = {
            LLVM_SYS_211_PREFIX = "${pkgs.libllvm.dev}";
          };
        in
        {
          packages.default = craneLib.buildPackage (
            llvmEnv
            // {
              src = craneLib.cleanCargoSource ./.;
              strictDeps = true;

              nativeBuildInputs = [ pkgs.libllvm.dev ];
              buildInputs = [
                pkgs.libffi
                pkgs.libiconv
              ];
            }
          );

          devShells.default = pkgs.mkShellNoCC (
            llvmEnv
            // {
              packages = [
                pkgs.clang
                pkgs.gnumake
                pkgs.nixfmt
                toolchain
                pkgs.libllvm.dev
                pkgs.libffi
                pkgs.libiconv
              ];

              LIBRARY_PATH = lib.makeLibraryPath [
                pkgs.libffi
                pkgs.libiconv
              ];
            }
          );

          treefmt.programs = {
            nixfmt.enable = true;
            mdformat.enable = true;
          };

          treefmt.settings.global.excludes = [
            ".agents/skills/**"
            ".claude/skills/**"
          ];
        };
    };
}
