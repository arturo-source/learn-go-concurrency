package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	c1, c2 := make(chan int), make(chan int)

	go func() {
		Walk(t1, c1)
		close(c1)
	}()
	go func() {
		Walk(t2, c2)
		close(c2)
	}()

	for {
		val1, exist1 := <-c1
		val2, exist2 := <-c2

		if val1 != val2 {
			return false
		}

		if exist1 != exist2 {
			return false
		}

		if exist1 == false {
			break
		}
	}

	return true
}

func main() {
	t := tree.New(1)
	c := make(chan int)
	go func() {
		Walk(t, c)
		close(c)
	}()

	for val := range c {
		fmt.Println(val)
	}

	t1, t2 := tree.New(1), tree.New(1)
	if Same(t1, t2) {
		fmt.Println(t1, t2, "are the same")
	}

	t3, t4 := tree.New(3), tree.New(4)
	if !Same(t3, t4) {
		fmt.Println(t3, t4, "are different")
	}
}
