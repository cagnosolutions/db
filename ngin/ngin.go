package ngin

import (
	"bytes"
	"os"
	"path/filepath"
	"syscall"
)

const (
	slab = (1 << 26) //   64 MB Slab Size
	page = (1 << 10) // 4096 KB Page Size
)

var empty = make([]byte, page)

type ngin struct {
	file *os.File
	indx *btree
	data mmap
}

func OpenNgin(path string) *ngin {
	_, err := os.Stat(path + `.dat`)
	// new instance
	if err != nil && !os.IsExist(err) {
		dirs, _ := filepath.Split(path)
		err := os.MkdirAll(dirs, 0755)
		if err != nil {
			panic(err)
		}
		fd, err := os.Create(path + `.dat`)
		if err != nil {
			panic(err)
		}
		if err := fd.Truncate(page + slab); err != nil {
			panic(err)
		}
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}
	// existing
	fd, err := os.OpenFile(path+`.dat`, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	info, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	return &ngin{
		file: fd,
		data: Mmap(fd, 0, int(info.Size())),
	}
}

func (n *ngin) set(d []byte, k int) {
	o := k * page
	if o+page >= len(e.data) {
		e.grow()
	}
	copy(n.data[o:], append(d, make([]byte, (page-len(d)))...))
}

func (n *ngin) get(k int) []byte {
	o := k * page
	if e.data[o] != 0x00 {
		if n := bytes.IndexByte(n.data[o:o+page], byte(0x00)); n > -1 {
			return n.data[o : o+n]
		}
	}
	return nil
}

func (n *ngin) del(k int) {
	o := k * page
	copy(e.data[o:], empty)
}

func (n *ngin) grow() {
	size := ((len(n.data) + (page + slab)) + page - 1) &^ (page - 1)
	n.data.Munmap()
	if err := syscall.Ftruncate(int(e.file.Fd()), int64(size)); err != nil {
		panic(err)
	}
	n.data = Mmap(n.file, 0, size)
}
