llvm-project:
	git clone --depth=1 https://github.com/llvm/llvm-project.git

llvm-build/build.ninja:
	mkdir -p $(dir $@) && cd $(dir $@) && \
	cmake -G Ninja ${CURDIR}/llvm-project/llvm \
	"-DLLVM_TARGETS_TO_BUILD=X86;ARM;AArch64;WebAssembly" \
	-DCMAKE_BUILD_TYPE=Release \
	-DLIBCLANG_BUILD_STATIC=ON \
	-DLLVM_ENABLE_TERMINFO=OFF \
	-DLLVM_ENABLE_ZLIB=OFF \
	-DLLVM_ENABLE_ZSTD=OFF \
	-DLLVM_ENABLE_LIBEDIT=OFF \
	-DLLVM_ENABLE_Z3_SOLVER=OFF \
	-DLLVM_ENABLE_OCAMLDOC=OFF \
	-DLLVM_ENABLE_LIBXML2=OFF \
	-DLLVM_ENABLE_PROJECTS="clang;lld" \
	-DLLVM_TOOL_CLANG_TOOLS_EXTRA_BUILD=OFF \
	-DCLANG_ENABLE_STATIC_ANALYZER=OFF \
	-DCLANG_ENABLE_ARCMT=OFF clang llvm-config llvm-ar llvm-nm lld

llvm-build: llvm-build/build.ninja
	cd $@ && ninja clang llvm-config llvm-ar llvm-nm lld
