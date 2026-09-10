{
  lib,
  craneLib,
  features,
  llvm,
}:
let
  inherit (lib) fileset;

  base = fileset.unions [
    ../Cargo.toml
    ../Cargo.lock
    (fileset.difference ../src (fileset.maybeMissing ../src/features))
  ];
in
craneLib.buildPackage {
  src = craneLib.cleanCargoSource (
    fileset.toSource {
      root = ../.;
      fileset = fileset.union base features.fileset;
    }
  );

  # Derivation features drop into place next to the in-repo ones.
  postUnpack = lib.concatLines (
    lib.mapAttrsToList (
      name: drv: "cp -r --no-preserve=mode ${drv} \"$sourceRoot/src/features/${name}\""
    ) features.drvs
  );

  strictDeps = true;

  inherit (llvm) LLVM_SYS_211_PREFIX nativeBuildInputs buildInputs;

  meta = {
    description = "Compiler for MangoLang (tbd)";
    mainProgram = "unmangc";
  };
}
