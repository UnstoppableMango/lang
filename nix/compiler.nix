{
  pkgs,
  craneLib,
  llvmEnv,
}:
craneLib.buildPackage (
  llvmEnv
  // {
    src = craneLib.cleanCargoSource ../.;
    strictDeps = true;

    nativeBuildInputs = [ pkgs.libllvm.dev ];
    buildInputs = [
      pkgs.libffi
      pkgs.libiconv
    ];
  }
)
