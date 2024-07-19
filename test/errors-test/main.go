package main

import (
	"errors"
	"fmt"
)

func main() {
	err := errors.New("wow")
	target := errors.Unwrap(err)
	fmt.Printf("%v\n", target)

	ptr := errors.TestSkip()
	fmt.Printf("ptr%v\n", ptr)

	val := errors.TestSkip2()
	fmt.Printf("val%v\n", val)
}
