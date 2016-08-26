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
		n := gen("%.3d", i)
		t.Set(n, n)
	}

	fmt.Printf("Tree contains %d entries...\n", t.Count())

	for i := 0; i < COUNT; i++ {
		n := gen("%.3d", i)
		x := t.Get(n)
		fmt.Printf("got val: %s\n", x)
	}

	t.Print()

	t.Set(gen("%.3d", 25), gen("%.3d", 25))
	t.Print()
	t.Set(gen("%.3d", 18), gen("%.3d", 18))
	t.Print()
	t.Set(gen("%.3d", 7), gen("%.3d", 777))
	t.Print()

	n := gen("%.3d", 7)
	x := t.Get(n)
	fmt.Printf("got val: %s\n", x)

	//t.BFS() // print out tree...

	t.Close()

}
