package db

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

const PAGE = 4096

var WIPE = make([]byte, PAGE)

type Engine struct {
	file *os.File
	data mmap
}

func OpenEngine(path string) *Engine {
	// if engine is new, enter if statement
	_, err := os.Stat(path + `.dat`)
	if err != nil && !os.IsExist(err) {
		// create directory path
		dirs, _ := filepath.Split(path)
		err := os.MkdirAll(dirs, 0755)
		if err != nil {
			panic(err)
		}
		// create data file, and truncate it
		fd, err := os.Create(path + `.dat`)
		if err != nil {
			panic(err)
		}
		// write an initial file size of 16MB
		if err := fd.Truncate(1 << 24); err != nil {
			panic(err)
		}
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}
	// open file to construct rest of data structure
	fd, err := os.OpenFile(path+`.dat`, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	info, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	B

	// map file into virtual address space, and sort
	e := &Engine{
		file: fd,
		data: Mmap(fd, 0, int(info.Size())),
	}
	return e
}

func (e *Engine) Set(d []byte, k int) {
	// get byte offset from position k
	of := k * PAGE
	// do a bounds check, grow if nessicary...
	if of+PAGE >= len(e.data) {
		e.grow()
	}
	// copy the data `one-off` the offset
	copy(e.data[of:], append(d, make([]byte, (PAGE-len(d)))...))
}

func (e *Engine) Get(k int) []byte {
	// get byte offset from position k
	of := k * PAGE
	if e.data[of] != 0x00 {
		if n := bytes.IndexByte(e.data[of:of+PAGE], byte(0x00)); n > -1 {
			return e.data[of : of+n]
		}
	}
	return nil
}

func (e *Engine) Del(k int) {
	// get byte offset from position k
	of := k * PAGE
	// copy number of pages * page size worth
	// of nil bytes starting at the k's offset
	copy(e.data[of:], WIPE)
}

func (e *Engine) Iter() <-chan Data {
	ch := make(chan Data)
	go func() {
		for i := 0; i < len(e.data); i += PAGE {
			if e.data[i] != 0x00 {
				if n := bytes.IndexByte(e.data[i:i+PAGE], byte(0x00)); n > -1 {
					ch <- Data{Block: int64(i / PAGE), Value: e.data[i : i+n]}
				}
			}
		}
		close(ch)
	}()
	return ch
}

type Data struct {
	Block int64
	Value []byte
}

func (e *Engine) grow() {
	// resize size to current size + 16MB chunk (grow in 16 MB chunks)
	size := ((len(e.data) + (1 << 24)) + PAGE - 1) &^ (PAGE - 1)
	// unmap current mapping before growing underlying file...
	e.data.Munmap()
	// truncate underlying file to updated size, check for errors
	err := syscall.Ftruncate(int(e.file.Fd()), int64(size))
	if err != nil {
		panic(err)
	}
	// remap underlying file now that it has grown
	e.data = Mmap(e.file, 0, size)
}

func (e *Engine) CloseEngine() {
	e.data.Sync()   // flush data to disk
	e.data.Munmap() // unmap memory mappings
	e.file.Close()  // close underlying file
}

func (e *Engine) Sort() {
	sort.Stable(e.data)
}
