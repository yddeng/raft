package raft

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

type MapStorage struct {
	mu   sync.Mutex
	m    map[string][]byte
	dir  string
	name string
}

func NewMapStorage(dir, name string) *MapStorage {
	os.Mkdir(dir, os.ModePerm)
	m := make(map[string][]byte)
	return &MapStorage{
		m:    m,
		dir:  dir,
		name: name,
	}
}

func (ms *MapStorage) Get(key string) ([]byte, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	v, found := ms.m[key]
	return v, found
}

func (ms *MapStorage) Set(key string, value []byte) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.m[key] = value
}

func (ms *MapStorage) HasData() bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return len(ms.m) > 0
}

func (ms *MapStorage) storeToDisk() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(ms.m)
	data := w.Bytes()
	filename := path.Join(ms.dir, fmt.Sprintf("rf-%s.state", ms.name))
	ms.mustWriteFile(filename, data)
}

func (ms *MapStorage) loadFromDisk() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	filename := path.Join(ms.dir, fmt.Sprintf("rf-%s.state", ms.name))
	data := ms.mustReadFile(filename)

	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&ms.m)
}

func (ms *MapStorage) mustReadFile(filename string) (data []byte) {
	var err error
	var file *os.File

	// 文件不存在则返回
	if _, err = os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}

	if file, err = os.Open(filename); err != nil {
		log.Fatal(err)
	}
	if data, err = ioutil.ReadAll(file); err != nil {
		log.Fatal(err)
	}
	file.Close()
	return
}

func (ms *MapStorage) mustWriteFile(filename string, data []byte) {
	var err error
	var file *os.File

	// 文件操作失败，直接宕掉即可
	tmpFilename := filename + "-tmp"
	if file, err = os.Create(tmpFilename); err != nil {
		log.Fatal(err)
	}
	if _, err = file.Write(data); err != nil {
		log.Fatal(err)
	}
	if err = file.Sync(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	if err = os.Rename(tmpFilename, filename); err != nil {
		log.Fatal(err)
	}
}
