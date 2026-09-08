{ pkgs, lib }:
{
  # The compiler's Source -> LLVM IR stage links against libLLVM via
  # inkwell/llvm-sys; unrelated to the separate clang linking stage
  # (LLVM IR -> binary) that `make hello` drives.
  LLVM_SYS_211_PREFIX = "${pkgs.libllvm.dev}";

  # libllvm.lib holds libLLVM.so, which the dynamically linked compiler needs at
  # run time; the linker finds it through the dev output and records no rpath,
  # so the path has to be baked in explicitly.
  libPath = "${pkgs.libllvm.lib}/lib";
  LIBRARY_PATH = lib.makeLibraryPath [
    pkgs.libffi
    pkgs.libiconv
  ];

  nativeBuildInputs = [ pkgs.libllvm.dev ];

  buildInputs = [
    pkgs.libffi
    pkgs.libiconv
    pkgs.libllvm.lib
  ];
}
