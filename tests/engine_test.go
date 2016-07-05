package main

import (
	"fmt"
	"testing"

	"github.com/cagnosolutions/db"
)

func Benchmark_Engine_Set(b *testing.B) {
	e := db.OpenEngine(`db/e1`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Set([]byte(fmt.Sprintf("e1-block-number-%x", i)), i)
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Engine_Get(b *testing.B) {
	e := db.OpenEngine(`db/e1`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if d := e.Get(0); d == nil {
			b.Errorf("Got nil...\n")
		}
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Engine_Del(b *testing.B) {
	e := db.OpenEngine(`db/e1`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Del(0)
	}
	b.StopTimer()
	e.CloseEngine()
}

func Benchmark_Engine_Grow(b *testing.B) {
	e := db.OpenEngine(`db/e1_grow`)
	b.ResetTimer()
	for i := 0; i < 4096*40; i++ {
		e.Set([]byte(fmt.Sprintf("e1-block-number-%x", i)), i)
	}
	b.StopTimer()
	e.CloseEngine()
}
