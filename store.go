package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

const DELIM = byte(0x1D) // group seporator (alternatives 1C: file seportator, 1E: record seporator, 1F: unit seporator)

var logger = log.New(os.Stderr, "::", log.Ldate|log.Ltime|log.Llongfile)

func Encode(k []byte, v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(k)+1+len(b) > PAGE {
		return nil, errors.New("data too large")
	}
	k = append(k, append([]byte{DELIM}, b...)...)
	return k, nil
}

func Decode(b []byte, v interface{}) error {
	if len(b) > PAGE {
		return errors.New("data too large")
	}
	n := bytes.IndexByte(b, DELIM)
	if n == -1 {
		return errors.New("incorrect data format")
	}
	return json.Unmarshal(b[n+1:], v)
}

type Store struct {
	en    *Engine
	ix    *Tree
	count int
	next  int
	sync.RWMutex
}

func NewStore(path string) *Store {
	st := &Store{
		en: OpenEngine(path),
		ix: NewTree(),
	}
	st.SortAndCompact()
	return st
}

func (st *Store) Add(k []byte, v interface{}) {
	b, err := Encode(k, v)
	if err != nil {
		logger.Fatal(err)
	}

	st.ix.Add(k, Itob(int64(st.next)))

	st.en.Set(b, st.next)
	st.next++
	st.count++
}

func (st *Store) Set() {

}

func (st *Store) Get() {

}

func (st *Store) Del() {

}

func (st *Store) SortAndCompact() {
	st.ix.Close()
	st.en.Sort()
	for dat := range st.en.Iter() {
		n := bytes.IndexByte(dat.Value, DELIM)
		if n == -1 {
			panic("range st.en.Iter() possible file corruption")
		}
		st.ix.Add(dat.Value[:n], Itob(dat.Block))
	}
	st.count = st.ix.Count()
	st.next = st.count
	log.Println("")
}
