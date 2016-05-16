package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

const PAGE = 4096

type Engine struct {
	file  *os.File
	size  int
	mmap  Mmap
	count int // number of records
	next  int
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
	e.mmap = OpenMmap(fd, 0, e.size)
	e.SetNext(0)
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
	// get page aligned size of record
	sz := (len(d) + 2 + PAGE - 1) &^ (PAGE - 1)
	// do a bounds check, grow if nessicary...
	if of+sz >= len(e.mmap) {
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

func (e *Engine) Put(d []byte, k int) {
	// get byte offset from position k
	of := k * PAGE
	// get page aligned size of record
	sz := (len(d) + 2 + PAGE - 1) &^ (PAGE - 1)

	// do a bounds check, grow if nessicary...
	if of+sz >= len(e.mmap) {
		e.grow()
	}

	// check status of record header
	empty := (e.mmap[of] == 0x00)
	// compair new data size to old data size to see if data will fit
	if !empty && sz > (int(e.mmap[of+1])*PAGE) { // new is larger
		// wipe old data
		copy(e.mmap[of:], make([]byte, (int(e.mmap[of+1])*PAGE)))

		// start at next loop until an empty slot is found where data fits in s "blocks"
		i := e.next
		for j, s := 0, sz/PAGE; i < len(e.mmap); {
			if e.mmap[i] == 0x00 {
				i += PAGE
				j++
				if j == s {
					of = i - ((s - 1) * PAGE)
					break
				}
				continue
			}
			i += int(e.mmap[i+1])
			j = 0
		}
		// check for growth in the case the offset is the same but the data does not fit there (no gaps were found in mmap)
		if i >= len(e.mmap) {
			e.grow()
			of = i * PAGE
		}
	}
	// resize data according to nearest page offset, and copy the data
	copy(e.mmap[of+2:], append(d, make([]byte, (sz-len(d)))...))
	// write the header to the offset
	e.mmap[of] = 0x01
	e.mmap[of+1] = byte(sz / PAGE)
	// check if we just added, or updated so we
	// know if we should augment e.next's offset
	if empty {
		e.SetNext(of + sz)
	}
}

func (e *Engine) Add(d []byte) int {
	return -1
}

func (e *Engine) Get(k int) []byte {
	// get byte offset from position k
	of := k * PAGE
	// get number of pages from record header
	sz := int(e.mmap[of+1])
	// return a copy of the record slice at
	// k's offset, up to number of used pages

	//d := make([]byte, sz*PAGE)
	//copy(d, e.mmap[of:of+(sz*PAGE)])
	//return d
	return e.mmap[of : of+(sz*PAGE)]
}

func (e *Engine) Del(k int) {
	// get byte offset from position k
	of := k * PAGE
	// get number of pages from record header
	sz := int(e.mmap[of+1])
	// copy number of pages * page size worth
	// of nil bytes starting at the k's offset
	copy(e.mmap[of:], make([]byte, sz*PAGE))
	// check to see if the offset we just
	// delted from is earlier in the array
	// and if it is, then reset next to k
	if of < e.next {
		e.next = of
	}
}

func (e *Engine) GetNext(sz int) int {
	// return the value next currently holds
	return e.next
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

/*
 *	Memory Mapping Functions & Methods
 */

/*type Mmap []byte

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
}*/
