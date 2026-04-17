// Ref: https://go.dev/tour/concurrency/8

package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

// 
//type Tree struct {
//    Left  *Tree
//    Value int
//    Right *Tree
//}

func WalkRecursive(t *tree.Tree, ch chan int) {
	if (t == nil) {
		return 
	}
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	WalkRecursive(t, ch)
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	
	for range cap(ch1) {
		n1 := <- ch1
		n2 := <- ch2
		
		if (n1 != n2) {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
