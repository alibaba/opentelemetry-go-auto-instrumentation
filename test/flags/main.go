package main

import (
	"fmt"
)

var Placeholder = ""

func main() {
	n, _ := fmt.Printf("helloworld%s", "ingodwetrust")
	println(n)
	println(fmt.Sprintf("placeholder:%s", Placeholder))
}
