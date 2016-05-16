package db

import (
	"bytes"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type Mmap []byte

var tmp = make([]byte, PAGE)

func OpenMmap(f *os.File, off, len int) Mmap {
	mmap, err := syscall.Mmap(int(f.Fd()), int64(off), len, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return mmap
}

func (m Mmap) Mlock() {
	err := syscall.Mlock(m)
	if err != nil {
		panic(err)
	}
}

func (m Mmap) Munlock() {
	err := syscall.Munlock(m)
	if err != nil {
		panic(err)
	}
}

func (m Mmap) Munmap() {
	err := syscall.Munmap(m)
	m = nil
	if err != nil {
		panic(err)
	}
}

func (m Mmap) Sync() {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)),
		uintptr(syscall.MS_ASYNC))
	if err != 0 {
		panic(err)
	}
}

func (m Mmap) Mremap(size int) Mmap {
	fd := uintptr(unsafe.Pointer(&m[0]))
	err := syscall.Munmap(m)
	m = nil
	if err != nil {
		panic(err)
	}
	err = syscall.Ftruncate(int(fd), int64(align(size)))
	if err != nil {
		panic(err)
	}
	m, err = syscall.Mmap(int(fd), int64(0), size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return m
}

// open file helper
func OpenFile(path string) (*os.File, string, int) {
	fd, err := os.OpenFile(path, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	fi, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	return fd, sanitize(fi.Name()), int(fi.Size())
}

func sanitize(path string) string {
	if path[len(path)-1] == '/' {
		return path[:len(path)-1]
	}
	if x := strings.Index(path, "."); x != -1 {
		return path[:x]
	}
	return path
}

// round up to nearest pagesize -- helper
func align(size int) int {
	if size > 0 {
		return (size + PAGE - 1) &^ (PAGE - 1)
	}
	return PAGE
}

// resize underlying file -- helper
func resize(fd uintptr, size int) int {
	err := syscall.Ftruncate(int(fd), int64(align(size)))
	if err != nil {
		panic(err)
	}
	return size
}

func (mm Mmap) Len() int {
	return len(mm) / PAGE
}

func (mm Mmap) Less(i, j int) bool {
	pi, pj := i*PAGE, j*PAGE

	if mm[pi] == 0x00 {
		if mm[pi] == mm[pj] {
			return true
		}
		return false
	}
	if mm[pj] == 0x00 {
		return true
	}

	return bytes.Compare(mm[pi:pi+PAGE], mm[pj:pj+PAGE]) == -1

}

func (mm Mmap) Swap(i, j int) {
	pi, pj := i*PAGE, j*PAGE

	copy(tmp, mm[pi:pi+PAGE])
	copy(mm[pi:pi+PAGE], mm[pj:pj+PAGE])
	copy(mm[pj:pj+PAGE], tmp)
}
