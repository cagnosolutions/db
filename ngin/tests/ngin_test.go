package main

import (
	"testing"

	"github.com/cagnosolutions/db/ngin"
)

func Benchmark_Ngin_Put(b *testing.B) {
	e := ngin.OpenNgin(`_db/test`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Put([]byte{'f', 'o', 'o', 'b', 'a', 'r', 'b', 'a', 'z', '!'}, 0)
	}
	b.StopTimer()
	e.CloseNgin()
}

func Benchmark_Ngin_Get(b *testing.B) {
	e := ngin.OpenNgin(`db/test`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if d := e.Get(0); d == nil {
			b.Errorf("Got nil...\n")
		}
	}
	b.StopTimer()
	e.CloseNgin()
}

func Benchmark_Ngin_Del(b *testing.B) {
	e := ngin.OpenNgin(`db/test`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Del(0)
	}
	b.StopTimer()
	e.CloseNgin()
}

func Benchmark_Ngin_Grow(b *testing.B) {
	e := ngin.OpenNgin(`db/test`)
	b.ResetTimer()
	for i := 0; i < 4096*40; i++ {
		e.Put([]byte{'f', 'o', 'o', 'b', 'a', 'r', 'b', 'a', 'z', '!'}, i)
	}
	b.StopTimer()
	e.CloseNgin()
}
