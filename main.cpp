// Tell go to skip this file until I figure out a better way
// +build cpp

#include <llvm/ADT/APFloat.h>
#include <llvm/ADT/STLExtras.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/Constants.h>
#include <llvm/IR/DerivedTypes.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/Module.h>
#include <llvm/IR/Type.h>
#include <llvm/IR/Verifier.h>
#include <stdio.h>

int main(void) {
	auto Context = std::make_unique<llvm::LLVMContext>();
	auto Builder = std::make_unique<llvm::IRBuilder<>>(*Context);
	printf("TEST\n");
	return 0;
}
