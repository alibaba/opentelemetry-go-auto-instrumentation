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

	arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10 := errors.TestGetSet(1, true, false, 3.14, "str", nil, map[int]bool{1: true}, map[int]bool{2: true}, make(chan int), []int{1, 2, 3})
	fmt.Printf("val%v %v %v %v %v %v %v %v %v %v\n", arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
}
