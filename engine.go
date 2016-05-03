package db

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const PAGE = 4096

type Engine struct {
	file *os.File
	size int
	mmap Mmap
	next int
}

func OpenEngine(path string) *Engine {
	// if engine is new, enter if statement
	_, err := os.Stat(path + `.dat`)
	if err != nil && !os.IsExist(err) {
		// create directory path
		dirs, _ := filepath.Split(path)
		err := os.MkdirAll(dirs, 0755)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		// create data file, and truncate it
		fd, err := os.Create(path + `.dat`)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		// write an initial file size of 16MB
		if err := fd.Truncate(1 << 24); err != nil {
			log.Fatalf("%s\n", err)
		}
		if err := fd.Close(); err != nil {
			log.Fatalf("%s\n", err)
		}
	}
	// open file to construct rest of data structure
	fd, err := os.OpenFile(path+`.dat`, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	info, err := fd.Stat()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	e := &Engine{
		file: fd,
		size: int(info.Size()),
	}
	e.mmap = OpenMmap(fd, 0, e.size)
	e.Next(0)
	return e
}

func (e *Engine) Next(i int) {
	j := i
	for j < len(e.mmap) {
		if e.mmap[j] == 0x00 {
			e.next = j
			return
		}
		j += int(e.mmap[j])
		return
	}
	if i == 0 {
		// grow file and remap because we
		// have no more empty slots left
		e.grow()
		e.Next(j)
	}
}

func (e *Engine) Put(d []byte, k int) {
	// get byte offset from position k
	of := k * PAGE
	// get page aligned size of record
	sz := (len(d) + 1 + PAGE - 1) &^ (PAGE - 1)
	// do a bounds check
	if of+sz >= len(e.mmap) {
		log.Fatalf("cannot put at offset %d, offset + record exceeds mapped reigon\n", of+sz)
	}
	// check status of record header
	added := (e.mmap[of] == 0x00)
	// resize data according to nearest page offset
	//d = append(d, make([]byte, (sz-len(d)))...)
	// copy the data `one-off` the offset
	copy(e.mmap[of+1:], append(d, make([]byte, (sz-len(d)))...))
	// write the header to the offset
	e.mmap[of] = byte(sz / PAGE)
	// check if we just added, or updated so we
	// know if we should augment e.next's offset
	if added {
		e.Next(of + sz)
	}
}

func (e *Engine) Get(k int) []byte {
	// get byte offset from position k
	of := k * PAGE
	// get number of pages from record header
	sz := int(e.mmap[of])
	// return a copy of the record slice at
	// k's offset, up to number of used pages

	//d := make([]byte, sz*PAGE)
	//copy(d, e.mmap[of:of+sz*PAGE])
	//return d
	return e.mmap[of : of+sz*PAGE]
}

func (e *Engine) Del(k int) {
	// get byte offset from position k
	of := k * PAGE
	// get number of pages from record header
	sz := int(e.mmap[of])
	// copy number of pages * page size worth
	// of nil bytes starting at the k's offset
	copy(e.mmap[of:], make([]byte, sz*PAGE))
	// check to see if the offset we just
	// delted from is earlier in the array
	// and if it is, then reset next to k
	if k < e.next {
		e.next = k
	}
}

func (e *Engine) GetNext() int {
	// return the value next currently holds
	return e.next
}

func (e *Engine) grow() {
	// resize size to current size + 16MB chunk (grow in 16 MB chunks)
	e.size = ((e.size + (1 << 24)) + PAGE - 1) &^ (PAGE - 1)
	// unmap current mapping before growing underlying file...
	e.mmap.Close()
	// truncate underlying file to updated size, check for errors
	err := syscall.Ftruncate(int(e.file.Fd()), int64(e.size))
	if err != nil {
		panic(err)
	}
	// remap underlying file now that it has grown
	e.mmap = OpenMmap(e.file, 0, e.size)
	log.Println("FILE WAS GROWN SUCCESSFULLY")
}

func (e *Engine) CloseEngine() {
	e.mmap.Sync()  // flush mmap to disk
	e.mmap.Close() // unmap memory mappings
	e.file.Close() // close underlying file
}

/*
 *	Memory Mapping Functions & Methods
 */

type Mmap []byte

func OpenMmap(f *os.File, off, len int) Mmap {
	mmap, err := syscall.Mmap(int(f.Fd()), int64(off), len, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return mmap
}

func (m Mmap) Close() {
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
