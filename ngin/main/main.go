package main

import (
	"fmt"

	"github.com/cagnosolutions/db/ngin"
)

const COUNT = 16

func gen(str string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(str, args...))
}

func main() {

	t := ngin.NewBTree()

	for i := 0; i < COUNT; i++ {
		n := gen("key-val-%.3d", i)
		t.Set(n, n)
	}

	for i := 0; i < COUNT; i++ {
		n := gen("key-val-%.3d", i)
		x := t.Get(n)
		fmt.Printf("got val: %s\n", x)
	}

	fmt.Printf("Tree contains %d entries...\n", t.Count())

	t.Print()
	t.BFS() // print out tree...

	t.Close()

}
