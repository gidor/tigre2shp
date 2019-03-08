package main

import (
	"fmt"

	"github.com/gen2brain/dlgs"
)

type items map[string]interface{}
type pt struct {
	x int
	y int
}

func pppmain() {
	directory, ok, err := dlgs.File("SElezione", "", true)
	if ok {
		fmt.Println("OK")
	}
	if err == nil {
		fmt.Println("NO error")
	}
	fmt.Println(directory)

	var p rune
	p = 'p'
	fmt.Println(p)
	fmt.Println(string(p))
	fmt.Println(string(p) == "p")
}
