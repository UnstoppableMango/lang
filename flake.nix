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
          inherit (inputs'.fenix.packages.stable) toolchain;
          craneLib = (crane.mkLib pkgs).overrideToolchain (_: toolchain);

          llvm = pkgs.callPackage ./nix/llvm.nix { };

          compiler = pkgs.callPackage ./nix/compiler.nix {
            inherit craneLib llvm;
          };
        in
        {
          packages = {
            inherit compiler;
            default = compiler;
          };

          devShells.default = pkgs.mkShellNoCC ({
            packages = [
              pkgs.clang
              pkgs.gnumake
              pkgs.nixfmt
              toolchain
            ]
            ++ llvm.nativeBuildInputs
            ++ llvm.buildInputs;

            inherit (llvm) LLVM_SYS_211_PREFIX LIBRARY_PATH;
            UNMANGC = lib.getExe compiler;
          });

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
