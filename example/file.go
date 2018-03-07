package main

import (
	"fmt"

	"github.com/sevenNt/hera"
)

func main() {
	if err := hera.LoadFromFile("cfg.toml", false); err != nil {
		panic("load from file fail")
	}
	ret := hera.GetString("a.b.c")
	fmt.Printf("a.b.c is %s", ret)
}
