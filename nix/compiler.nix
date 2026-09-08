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
    ../tests
    (fileset.difference ../src (fileset.maybeMissing ../src/features))
  ];
in
craneLib.buildPackage {
  # No cleanCargoSource: it keeps only *.rs, *.toml and Cargo.lock, which would
  # silently drop every golden test case and leave a suite that tests nothing.
  # The fileset above already names exactly what belongs in the source.
  src = fileset.toSource {
    root = ../.;
    fileset = fileset.union base features.fileset;
  };

  # Derivation features drop into place next to the in-repo ones.
  postUnpack = lib.concatLines (
    lib.mapAttrsToList (
      name: drv: "cp -r --no-preserve=mode ${drv} \"$sourceRoot/src/features/${name}\""
    ) features.drvs
  );

  strictDeps = true;

  inherit (llvm) LLVM_SYS_211_PREFIX nativeBuildInputs buildInputs;

  RUSTFLAGS = "-C link-arg=-Wl,-rpath,${llvm.libPath}";

  meta = {
    description = "Compiler for MangoLang (tbd)";
    mainProgram = "unmangc";
  };
}
