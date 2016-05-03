package main

import (
	"testing"

	"github.com/cagnosolutions/db"
	"github.com/davecheney/profile"
)

func init() {
	defer profile.Start(profile.MemProfile).Stop()
}

var e = db.OpenEngine(`db/test`)

func Benchmark_Engine_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Put([]byte(`foo bar baz 0 0 0 0 0 0 0 0 0 0`), 0)
	}
}

func Benchmark_Engine_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if d := e.Get(0); d == nil {
			b.Errorf("Got nil...\n")
		}
	}
}

func Benchmark_Engine_Del(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Del(0)
	}
}
