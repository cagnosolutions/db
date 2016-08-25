package ngin

import (
	"fmt"
	"unsafe"
)

type T struct {
	offset, length uint16
}

func (t *T) Encode(s string, v *[]byte) {
	t.offset = uint16(len(*v))
	t.length = uint16(len(s))
	*v = append(*v, []byte(s)...)
}

func (t *T) Bytes(v []byte) []byte {
	return (*[0xffff]byte)(unsafe.Pointer(&v[t.offset]))[:t.length]
}

func (t *T) String(v []byte) string {
	return string(t.Bytes(v))
}

func main() {

	var t T

	// encode
	b := make([]byte, unsafe.Sizeof(t))
	t.Encode("This is a test", &b)
	copy(b, (*[unsafe.Sizeof(t)]byte)(unsafe.Pointer(&t))[:])
	fmt.Printf("[b] addr: %p, type: %T, value: %x, size: %d\n", &b, &b, &b, unsafe.Sizeof(&b))

	// decode
	s := ((*T)(unsafe.Pointer(&b[0])))
	fmt.Printf("[s] addr: %p, type: %T, value: %q, size: %d\n", s, s, s.String(b), unsafe.Sizeof(s))
}
