package db

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

const (
	PAGE_SIZE = 4096
	SLAB_SIZE = 1

	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

type Engine struct {
	file *os.File
	size int
	mmap Data
	page int
	slab int
	next int
}

func NewEngine(path string, slab int) error {
	if info, err := os.Stat(path + ".data"); err != nil && info.Size() > 1 {
		return err
	}
	if slab == 0 || slab > PAGE_SIZE {
		return errors.New("Slab size must be greater than 0, but less than" + strconv.Itoa(PAGE_SIZE))
	}
	fd, err := os.Open(path + ".meta")
	if err != nil {
		return err
	}
	if _, err := fd.Write(`{"page":` + strconv.Itoa(PAGE_SIZE) + `,"slab":` + slab + `,"next":0}`); err != nil {
		return err
	}
	if err := fd.Close(); err != nil {
		return err
	}
	fd, err = os.Open(path + ".data")
	if err != nil {
		return err
	}
	if err := fd.Truncate(16 * MB); err != nil {
		return err
	}
	if err := fd.Close(); err != nil {
		return err
	}
	return nil
}

func OpenEngine(path string) (*Engine, error) {
	info, err := os.Stat(path + ".data")
	if err != nil {
		return nil, err
	}
	var b []byte
	fd, err := os.Open(path + ".meta")
	if err != nil {
		return nil, err
	}
	_, err = fd.Read(b)
	if err != nil {
		return nil, err
	}
	if err := fd.Close(); err != nil {
		return nil, err
	}
	var meta map[string]int
	if err := json.Unmarshal(b, &meta); err != nil {
		return nil, err
	}
	fd, err := os.Open(path + ".data")
	if err != nil {
		return nil, err
	}
	return &Engine{
		file: fd,
		size: info.Size(),
		mmap: nil,
		page: meta["page"],
		slab: meta["slab"],
		next: meta["next"],
	}
}
