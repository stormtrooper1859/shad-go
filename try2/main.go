package main

type cc struct {
	a, b int
}

func Foo(a []int) {
	for i := 0; i < len(a); i++ {
		a[i] = 0
	}


	//_ = a[3]
	//a[0] = 10
	//a[1] = 10
	//a[2] = 10
	//a[3] = 10
}

func main() {
	a := new(int)
	_ = a
}
