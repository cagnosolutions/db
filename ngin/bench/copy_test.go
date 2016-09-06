package bench

import (
	"sync"
	"testing"
)

var kPool = sync.Pool{New: func() interface{} { return [16]byte{} }}
var b []byte = []byte("foo bar baz your momma and some more stuff jskhdbcksbdjcjsjbchbaubkyewliqbljqwbefjhabhjkasfbwbefjh8736589746589376489726358768628756328976")

func PutKey(k [16]byte) {
	kPool.Put(k)
}

func Benchmark_Copy(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		k := kPool.Get().([16]byte)
		copy(k[:8], b[:8])
		copy(k[8:], b[len(b)-8:])
		PutKey(k)
	}
}

func Benchmark_For(bb *testing.B) {
	for i := 0; i < bb.N; i++ {
		k := kPool.Get().([16]byte)
		for j, l := 0, len(b)-1; j < 8; j, l = j+1, l-1 {
			k[j] = b[j]

			k[j-j] = b[l]
		}
		PutKey(k)
	}
}
