{ pkgs, lib }:
{
  # The compiler's Source -> LLVM IR stage links against libLLVM via
  # inkwell/llvm-sys; unrelated to the separate clang linking stage
  # (LLVM IR -> binary) that `make hello` drives.
  LLVM_SYS_211_PREFIX = "${pkgs.libllvm.dev}";
  LIBRARY_PATH = lib.makeLibraryPath [
    pkgs.libffi
    pkgs.libiconv
  ];

  nativeBuildInputs = [ pkgs.libllvm.dev ];

  buildInputs = [
    pkgs.libffi
    pkgs.libiconv
  ];
}
