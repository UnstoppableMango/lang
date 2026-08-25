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
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;
      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        { pkgs, lib, ... }:
        {
          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              clang
              gnumake
              nixfmt
              rustc
              cargo
              libllvm.dev
              libffi
              libiconv
            ];

            # hack/'s Source -> LLVM IR stage links against libLLVM via inkwell/llvm-sys;
            # both vars are for that build only, unrelated to the clang linking stage.
            LLVM_SYS_211_PREFIX = "${pkgs.libllvm.dev}";
            LIBRARY_PATH = lib.makeLibraryPath [
              pkgs.libffi
              pkgs.libiconv
            ];
          };

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
