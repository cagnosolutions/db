package ngin

import (
	"os"
	"path/filepath"
)

const (
	db_size = (1 << 26) // 64MB
)

type ngin struct {
	file *os.File
	indx *tree
	data mmap
	free freeset
	curs int
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
		if err := fd.Truncate(db_size); err != nil {
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
	e := &ngin{
		file: fd,
		data: Mmap(fd, 0, int(info.Size())),
	}
	return e
}
