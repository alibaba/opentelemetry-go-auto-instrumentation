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

func TestGetSet(arg1 int, arg2, arg3 bool, arg4 float64, arg5 string,
	arg6 interface{}, arg7, arg8 map[int]bool, arg9 chan int, arg10 []int) (int, bool, bool, float64, string, interface{}, map[int]bool, map[int]bool, chan int, []int) {
	return arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10
}

type Recv struct{ X int }

func (t *Recv) TestGetSetRecv(arg1 int, arg2 float64) (int, float64) {
	return arg1, arg2
}

func OnlyRet() (int, string) {
	return 1024, "gansu"
}

func OnlyArgs(arg1 int, arg2 string) {
	println(arg1, arg2)
}
