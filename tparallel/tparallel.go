// +build !solution

package tparallel

import "fmt"

type T struct {
	completeTest chan int
	unlock chan int
}

func (t *T) Parallel() {

}

func (t *T) Run(subtest func(t *T)) {
	panic("implement me")
}

func Run(topTests []func(t *T)) {
	t := &T{}
	for _, f := range topTests {
		go func(f func(t *T)) {
			defer func() {
				a := recover()
				fmt.Println(a)
			}()
			f(t)

			select {

			}

		}(f)

	}
}
