package main

import (
	"github.com/YReshetko/go-annotation/example/intpkg"
)

func main() {
	runTest()
}

// runTest comment 1
func runTest() {
	runInternal()
}

// runInternal comment 2
func runInternal() {
	intpkg.Print()
}
