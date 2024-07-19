package main

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	n, _ := fmt.Printf("helloworld%s", "ingodwetrust")
	println(n)

	println(rate.Every(time.Duration(1) * time.Second))
}
