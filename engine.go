package db

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

const PAGE = 4096

var WIPE = make([]byte, PAGE)

type Engine struct {
	file *os.File
	size int
	mmap Mmap
	indx *Tree
	recs int // number of records
	next int
}

func (e Engine) PrintMMap() {
	fmt.Printf("%s\n", e.mmap)
}

func (e Engine) SortMmap() {
	sort.Stable(e.mmap)
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
	// map file into virtual address space, and sort
	e.mmap = OpenMmap(fd, 0, e.size)
	e.SortMmap()
	// initialize engine's primary index
	e.indx = NewTree()
	buf := make([]byte, 10)
	for i := 0; i < len(e.mmap); i += PAGE {
		e.indx.Add(e.mmap.Key(i), binary.BigEndian.PutUint64(&buf, int64(i)))
	}
	return e
}

func (e *Engine) SetNext(i int) {
	for i < len(e.mmap) && i < e.next {
		if e.mmap[i] == 0x00 {
			e.next = i
			return
		}
		i += int(e.mmap[i+1])
		return
	}
}

func (e *Engine) Set(d []byte, k int) {
	// get byte offset from position k
	of := k * PAGE
	// do a bounds check, grow if nessicary...
	if of+PAGE >= len(e.mmap) {
		// last record check empty
		if e.count < len(e.mmap)/PAGE {
			sort.Stable(e.mmap)
			// rebuild index
			// get next
		}
		e.grow()
		//log.Fatalf("cannot put at offset %d, offset + record exceeds mapped reigon\n", of+sz)
	}
	// check status of record header
	added := (e.mmap[of] == 0x00)
	// resize data according to nearest page offset
	//d = append(d, make([]byte, (sz-len(d)))...)
	// copy the data `one-off` the offset
	copy(e.mmap[of+2:], append(d, make([]byte, (sz-len(d)))...))
	// write the header to the offset
	e.mmap[of] = 0x01
	e.mmap[of+1] = byte(sz / PAGE)
	// check if we just added, or updated so we
	// know if we should augment e.next's offset
	if added {
		e.SetNext(of + sz)
	}
}

func (e *Engine) Add(d []byte) int {
	return -1
}

func (e *Engine) Get(k int) []byte {
	// get byte offset from position k
	of := k * PAGE
	return e.mmap[of : of+PAGE]
}

func (e *Engine) Del(k int) {
	// get byte offset from position k
	of := k * PAGE
	// copy number of pages * page size worth
	// of nil bytes starting at the k's offset
	copy(e.mmap[of:], WIPE)
	e.count--
}

func (e *Engine) grow() {
	// resize size to current size + 16MB chunk (grow in 16 MB chunks)
	e.size = ((e.size + (1 << 24)) + PAGE - 1) &^ (PAGE - 1)
	// unmap current mapping before growing underlying file...
	e.mmap.Munmap()
	// truncate underlying file to updated size, check for errors
	err := syscall.Ftruncate(int(e.file.Fd()), int64(e.size))
	if err != nil {
		panic(err)
	}
	// remap underlying file now that it has grown
	e.mmap = OpenMmap(e.file, 0, e.size)
}

func (e *Engine) CloseEngine() {
	e.mmap.Sync()   // flush mmap to disk
	e.mmap.Munmap() // unmap memory mappings
	e.file.Close()  // close underlying file
}
