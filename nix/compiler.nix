{
  craneLib,
  llvm,
}:
craneLib.buildPackage {
  src = craneLib.cleanCargoSource ../.;
  strictDeps = true;

  inherit (llvm) LLVM_SYS_211_PREFIX nativeBuildInputs buildInputs;

  meta = {
    description = "Compiler for MangoLang (tbd)";
    mainProgram = "unmangc";
  };
}
