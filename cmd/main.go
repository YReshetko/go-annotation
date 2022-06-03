package main

import (
	"github.com/YReshetko/go-annotation/pkg"

	_ "github.com/YReshetko/go-annotation/cmd/processors/endpoint"
	_ "github.com/YReshetko/go-annotation/cmd/processors/inject"
	_ "github.com/YReshetko/go-annotation/cmd/processors/multiline"
	_ "github.com/YReshetko/go-annotation/cmd/processors/rest"
	_ "github.com/YReshetko/go-annotation/cmd/processors/schema"
)

func main() {
	pkg.Process()
}
