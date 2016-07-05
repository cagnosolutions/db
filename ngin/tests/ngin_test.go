package main

import (
	"fmt"
	"testing"

	"github.com/cagnosolutions/db/ngin"
)

func Benchmark_Ngin_Set(b *testing.B) {
	e := ngin.OpenEngine(`db/e2`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Set([]byte(fmt.Sprintf("e2-block-number-%x", i)), i)
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Ngin_Get(b *testing.B) {
	e := ngin.OpenEngine(`db/e2`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if d := e.Get(0); d == nil {
			b.Errorf("Got nil...\n")
		}
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Ngin_Del(b *testing.B) {
	e := ngin.OpenEngine(`db/e2`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Del(0)
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Ngin_Grow(b *testing.B) {
	e := ngin.OpenEngine(`db/e2_grow`)
	b.ResetTimer()
	for i := 0; i < 4096*40; i++ {
		e.Set([]byte(fmt.Sprintf("e2-block-number-%x", i)), i)
	}
	b.StopTimer()
	e.CloseEngine()
}
