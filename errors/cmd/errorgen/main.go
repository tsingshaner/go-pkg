package main

import "github.com/tsingshaner/go-pkg/errors/gen"

func main() {
	config := gen.ReadErrors()
	gen.GeneratePkg(config)
}
