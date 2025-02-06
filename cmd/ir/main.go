package main

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

func main() {
	c := llvm.NewContext()
	defer c.Dispose()

	m := c.NewModule("test")

	v := llvm.ConstString("test", true)
	fmt.Println(v.String())

	m.Dump()
}
