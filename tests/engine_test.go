package main

import (
	"log"
	"testing"

	"github.com/cagnosolutions/db"
)

var e = db.OpenEngine(`db/test`)

func Benchmark_Engine_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Put([]byte(`foo bar baz 0 0 0 0 0 0 0 0 0 0`), 0)
		e.Put([]byte(`foo bar baz 1 1 1 1 1 1 1 1 1 1`), 1)
		e.Put([]byte(`foo bar baz 2 2 2 2 2 2 2 2 2 2`), 2)
		e.Put([]byte(`foo bar baz 3 3 3 3 3 3 3 3 3 3`), 3)
		e.Put([]byte(`foo bar baz 4 4 4 4 4 4 4 4 4 4`), 4)
		e.Put([]byte(`foo bar baz 5 5 5 5 5 5 5 5 5 5`), 5)
		e.Put([]byte(`foo bar baz 6 6 6 6 6 6 6 6 6 6`), 6)
		e.Put([]byte(`foo bar baz 7 7 7 7 7 7 7 7 7 7`), 7)
		e.Put([]byte(`foo bar baz 8 8 8 8 8 8 8 8 8 8`), 8)
		//e.Put([]byte(`foo bar baz 9 9 9 9 9 9 9 9 9 9`), 9)
		log.Println("next:", e.GetNext())
	}
}
