//go:build ignore

package errors

func TestSkip() *int {
	val := 3
	return &val
}

func TestSkip2() int {
	return 1024
}

func p1() {}
func p2() {}

func p3(arg1 int, arg2 bool, arg3 float64) (int, bool, float64) {
	return arg1, arg2, arg3
}
