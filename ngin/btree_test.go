package ngin

import (
	"fmt"
	"testing"
)

func Benchmark_BTree_Set(b *testing.B) {
	t := &btree{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Set([]byte(fmt.Sprintf("e1-block-number-%x", i)), i)
	}
	b.StopTimer()
	e.CloseBTree()
}

func Benchmark_BTree_Get(b *testing.B) {
	e := db.OpenBTree(`db/e1`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if d := e.Get(0); d == nil {
			b.Errorf("Got nil...\n")
		}
	}
	b.StopTimer()
	e.CloseBTree()
}

func Benchmark_BTree_Del(b *testing.B) {
	e := db.OpenBTree(`db/e1`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Del(0)
	}
	b.StopTimer()
	e.CloseBTree()
}

func Benchmark_BTree_Grow(b *testing.B) {
	e := db.OpenBTree(`db/e1_grow`)
	b.ResetTimer()
	for i := 0; i < 4096*40; i++ {
		e.Set([]byte(fmt.Sprintf("e1-block-number-%x", i)), i)
	}
	b.StopTimer()
	e.CloseBTree()
}
