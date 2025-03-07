package main

import (
	"tinygo.org/x/go-llvm"
)

func main() {
	c := llvm.NewContext()
	defer c.Dispose()

	m := c.NewModule("test")
	// b := c.NewBuilder()

	llvm.FunctionType(c.Int32Type(), []llvm.Type{
		llvm.PointerType(c.Int8Type(), 1),
		llvm.PointerType(c.Int8Type(), 1),
	}, true)

	m.Dump()
}
